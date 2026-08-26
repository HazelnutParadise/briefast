package seo

import (
	"context"
	"html"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
	sy "github.com/HazelnutParadise/syralit"
)

const (
	siteName  = "Briefast"
	siteTitle = "Briefast｜每日產業與股市報告"
	// siteDescription 在查無報告或查詢失敗時使用，描述站台本身而非某一天的內容。
	siteDescription    = "Briefast 每個交易日開盤前出刊，由 AI 彙整當日產業要聞、個股新聞與多空判斷。"
	historyDescription = "Briefast 歷次每日產業與股市報告，依日期瀏覽完整晨報內容。"
	// descriptionLimit 是描述的字元上限，超過就截斷並補上刪節號。
	descriptionLimit = 150
)

// Reports 是中繼資料需要的報告查詢，*store.Store 直接滿足這個介面。抽成介面是
// 為了讓查詢失敗的降級路徑可以被測試。
type Reports interface {
	LatestReport(ctx context.Context) (*report.Report, error)
	ReportByDate(ctx context.Context, date string) (*report.Report, error)
	ListReports(ctx context.Context, page, pageSize int) ([]store.ReportSummary, error)
	CountReports(ctx context.Context) (int, error)
}

// Deps 是套件對外的相依組合，robots、sitemap 與 head 改寫都從這裡取用。
type Deps struct {
	Reports Reports
	Config  Config
}

// PageKind 指出這一份 DocumentFunc 服務哪一種頁面。歷史頁掛在 StripPrefix 之
// 後，收到的請求路徑一律是 "/"，所以種類必須由掛載端指定，不能從路徑推斷。
type PageKind int

const (
	PageHome PageKind = iota
	PageHistory
)

// DocumentFunc 回傳可直接交給 sy.Config.DocumentFunc 的函式。Lang 與 Dir 留空，
// 讓 syralit.toml 的設定繼續套用到每一頁。
func (d Deps) DocumentFunc(kind PageKind) func(*http.Request) sy.Document {
	return func(r *http.Request) sy.Document {
		meta := d.metaFor(kind, r)
		return sy.Document{Title: meta.Title, HeadHTML: renderMetaTags(meta)}
	}
}

// pageMeta 是單一頁面要寫進 head 的內容。
type pageMeta struct {
	Title        string
	Description  string
	CanonicalURL string
	OGType       string
}

// metaFor 依請求路徑與 date 查詢參數決定這一頁的中繼資料。查詢失敗或查無報告
// 一律退回站台預設值，不讓中繼資料的問題變成頁面錯誤。
func (d Deps) metaFor(kind PageKind, r *http.Request) pageMeta {
	base := d.Config.BaseURL(r)
	fallback := pageMeta{Title: siteTitle, Description: siteDescription, OGType: "website"}

	if kind == PageHistory {
		date := strings.TrimSpace(r.URL.Query().Get("date"))
		if date == "" {
			return pageMeta{
				Title:        "歷史報告｜" + siteName,
				Description:  historyDescription,
				CanonicalURL: base + "/history/",
				OGType:       "website",
			}
		}
		// 找不到就把 canonical 指回列表頁，避免把使用者帶進來的字串回寫進頁面。
		fallback.CanonicalURL = base + "/history/"
		found, err := d.Reports.ReportByDate(r.Context(), date)
		if err != nil || found == nil {
			return fallback
		}
		return pageMeta{
			Title:        found.Headline + "｜" + siteName + " " + found.Date,
			Description:  descriptionOf(found),
			CanonicalURL: base + "/history/?date=" + found.Date,
			OGType:       "article",
		}
	}

	fallback.CanonicalURL = base + "/"
	latest, err := d.Reports.LatestReport(r.Context())
	if err != nil || latest == nil {
		return fallback
	}
	return pageMeta{
		Title:        latest.Headline + "｜" + siteName,
		Description:  descriptionOf(latest),
		CanonicalURL: base + "/",
		OGType:       "article",
	}
}

// descriptionOf 取報告總覽當描述，總覽是空的就退回站台預設描述。
func descriptionOf(r *report.Report) string {
	text := truncateRunes(plainText(r.OverviewMD), descriptionLimit)
	if text == "" {
		return siteDescription
	}
	return text
}

// renderMetaTags 產生 Document.HeadHTML。標題不在這裡，由 Document.Title 帶出。
func renderMetaTags(meta pageMeta) string {
	title := html.EscapeString(meta.Title)
	description := html.EscapeString(meta.Description)
	canonical := html.EscapeString(meta.CanonicalURL)

	var b strings.Builder
	b.WriteString(`<meta name="description" content="` + description + `">`)
	b.WriteString("\n" + `<link rel="canonical" href="` + canonical + `">`)
	b.WriteString("\n" + `<meta property="og:type" content="` + meta.OGType + `">`)
	b.WriteString("\n" + `<meta property="og:site_name" content="` + siteName + `">`)
	b.WriteString("\n" + `<meta property="og:title" content="` + title + `">`)
	b.WriteString("\n" + `<meta property="og:description" content="` + description + `">`)
	b.WriteString("\n" + `<meta property="og:url" content="` + canonical + `">`)
	b.WriteString("\n" + `<meta property="og:locale" content="zh_TW">`)
	b.WriteString("\n" + `<meta name="twitter:card" content="summary">`)
	return b.String()
}

var (
	markdownImage = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	markdownLink  = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	listMarker    = regexp.MustCompile(`^\s*(?:[-*+]|\d+\.)\s+`)
	quoteMarker   = regexp.MustCompile(`^\s*>+\s*`)
)

// plainText 把 Markdown 轉成適合放進 meta description 的一行純文字。標題行整行
// 捨棄，因為那是段落標籤而不是內容。
func plainText(markdown string) string {
	var lines []string
	for _, line := range strings.Split(markdown, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = quoteMarker.ReplaceAllString(line, "")
		line = listMarker.ReplaceAllString(line, "")
		line = markdownImage.ReplaceAllString(line, "")
		line = markdownLink.ReplaceAllString(line, "$1")
		line = strings.NewReplacer("*", "", "`", "").Replace(line)
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, line)
	}
	return tidySpaces(strings.Join(lines, "\n"))
}

// tidySpaces 把連續空白壓成一個空格，並拿掉夾在兩個中日韓字元之間的空格——那個
// 空格是 Markdown 排版留下的，讀起來像斷句錯誤。
func tidySpaces(value string) string {
	var out []rune
	pendingSpace := false
	for _, r := range value {
		if unicode.IsSpace(r) {
			pendingSpace = true
			continue
		}
		if pendingSpace && len(out) > 0 {
			if !(isWide(out[len(out)-1]) && isWide(r)) {
				out = append(out, ' ')
			}
		}
		pendingSpace = false
		out = append(out, r)
	}
	return string(out)
}

func isWide(r rune) bool { return r > unicode.MaxASCII }

func truncateRunes(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit]) + "…"
}
