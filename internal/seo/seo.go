// Package seo 讓公開頁對搜尋引擎與社群分享爬蟲可讀。Syralit 只送空殼 HTML，
// 內容靠 WebSocket 推送，爬蟲拿不到；這個套件在回應送出前改寫 head，另外提供
// robots.txt 與 sitemap.xml。
package seo

import (
	"net/http"
	"strings"
)

// Config 是站台層級的設定。SiteURL 為空時，絕對網址改由請求標頭推導。
type Config struct {
	SiteURL string
}

// BaseURL 回傳不含結尾斜線的站台絕對網址。解析順序是設定值、X-Forwarded-Proto
// 加請求主機、最後才依連線是否為 TLS 決定通訊協定。
func (c Config) BaseURL(r *http.Request) string {
	if configured := strings.TrimRight(strings.TrimSpace(c.SiteURL), "/"); configured != "" {
		return configured
	}
	scheme := "http"
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); forwarded != "" {
		// 代理串接時標頭可能是逗號分隔的清單，最前面那個才是使用者端用的協定。
		if index := strings.Index(forwarded, ","); index >= 0 {
			forwarded = strings.TrimSpace(forwarded[:index])
		}
		scheme = forwarded
	} else if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
