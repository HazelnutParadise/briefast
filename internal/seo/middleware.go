package seo

import (
	"bufio"
	"bytes"
	"errors"
	"html"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// frameworkPrefix 是 Syralit 自己的資源與即時通道路徑。首頁掛在 "/" 之下，歷史頁
// 掛在 "/history/" 之下，兩者的 WebSocket 與 SSE 都帶著這一段。
const frameworkPrefix = "/_syralit/"

// Middleware 在回應送出前改寫公開頁的 head。只有狀態 200 且 Content-Type 是
// text/html 的回應會被緩衝，其餘一律原封不動直通。
func (d Deps) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, frameworkPrefix) {
			next.ServeHTTP(w, r)
			return
		}
		catcher := &interceptor{ResponseWriter: w}
		next.ServeHTTP(catcher, r)
		catcher.finish(d, r)
	})
}

// interceptor 先看標頭再決定要不要把 body 收進緩衝區。決定不緩衝之後就退化成
// 單純的轉發，串流與連線升級都不受影響。
type interceptor struct {
	http.ResponseWriter
	status    int
	decided   bool
	buffering bool
	body      bytes.Buffer
}

func (i *interceptor) WriteHeader(status int) {
	i.decide(status)
	if !i.buffering {
		i.ResponseWriter.WriteHeader(status)
	}
}

func (i *interceptor) Write(p []byte) (int, error) {
	if !i.decided {
		i.decide(http.StatusOK)
		if !i.buffering {
			i.ResponseWriter.WriteHeader(i.status)
		}
	}
	if i.buffering {
		return i.body.Write(p)
	}
	return i.ResponseWriter.Write(p)
}

// decide 只跑一次，決定這個回應要不要緩衝。
func (i *interceptor) decide(status int) {
	if i.decided {
		return
	}
	i.decided = true
	i.status = status
	contentType := i.Header().Get("Content-Type")
	i.buffering = status == http.StatusOK && strings.HasPrefix(contentType, "text/html")
}

// Flush 代表下游正在串流。緩衝會把串流卡住，所以此時放棄改寫，把已收到的內容
// 立刻送出並改成直通。
func (i *interceptor) Flush() {
	if i.buffering {
		i.buffering = false
		i.ResponseWriter.WriteHeader(i.status)
		if i.body.Len() > 0 {
			_, _ = i.ResponseWriter.Write(i.body.Bytes())
			i.body.Reset()
		}
	}
	if flusher, ok := i.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack 讓 WebSocket 升級能拿到底層連線。
func (i *interceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := i.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("seo: underlying ResponseWriter does not support hijacking")
	}
	i.buffering = false
	return hijacker.Hijack()
}

// finish 把緩衝起來的 HTML 改寫後送出。找不到可替換的 title 就原樣送出。
func (i *interceptor) finish(d Deps, r *http.Request) {
	if !i.decided {
		i.ResponseWriter.WriteHeader(http.StatusOK)
		return
	}
	if !i.buffering {
		return
	}
	body := rewriteHead(i.body.String(), d.metaFor(r))
	i.Header().Set("Content-Length", strconv.Itoa(len(body)))
	i.ResponseWriter.WriteHeader(i.status)
	_, _ = i.ResponseWriter.Write([]byte(body))
}

// rewriteHead 把整個 title 元素替換成標題加上中繼資料標籤，找不到就整份原樣
// 回傳。語言標示由 syralit.toml 的 lang 提供，這裡不碰。
func rewriteHead(document string, meta pageMeta) string {
	start := strings.Index(document, "<title>")
	if start < 0 {
		return document
	}
	closing := strings.Index(document[start:], "</title>")
	if closing < 0 {
		return document
	}
	end := start + closing + len("</title>")
	return document[:start] + renderMetaTags(meta) + document[end:]
}

func renderMetaTags(meta pageMeta) string {
	title := html.EscapeString(meta.Title)
	description := html.EscapeString(meta.Description)
	canonical := html.EscapeString(meta.CanonicalURL)

	var b strings.Builder
	b.WriteString("<title>" + title + "</title>")
	b.WriteString("\n<meta name=\"description\" content=\"" + description + "\">")
	b.WriteString("\n<link rel=\"canonical\" href=\"" + canonical + "\">")
	b.WriteString("\n<meta property=\"og:type\" content=\"" + meta.OGType + "\">")
	b.WriteString("\n<meta property=\"og:site_name\" content=\"" + siteName + "\">")
	b.WriteString("\n<meta property=\"og:title\" content=\"" + title + "\">")
	b.WriteString("\n<meta property=\"og:description\" content=\"" + description + "\">")
	b.WriteString("\n<meta property=\"og:url\" content=\"" + canonical + "\">")
	b.WriteString("\n<meta property=\"og:locale\" content=\"zh_TW\">")
	b.WriteString("\n<meta name=\"twitter:card\" content=\"summary\">")
	return b.String()
}
