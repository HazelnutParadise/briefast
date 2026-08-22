package seo

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
)

func TestBaseURLPrefersConfiguredSiteURL(t *testing.T) {
	cfg := Config{SiteURL: "https://briefast.example"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "internal.local:8600"
	req.Header.Set("X-Forwarded-Proto", "http")

	if got := cfg.BaseURL(req); got != "https://briefast.example" {
		t.Fatalf("BaseURL() = %q, want %q", got, "https://briefast.example")
	}
}

func TestBaseURLTrimsTrailingSlashFromConfiguredSiteURL(t *testing.T) {
	cfg := Config{SiteURL: "https://briefast.example/"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	if got := cfg.BaseURL(req); got != "https://briefast.example" {
		t.Fatalf("BaseURL() = %q, want %q", got, "https://briefast.example")
	}
}

func TestBaseURLUsesForwardedProto(t *testing.T) {
	cfg := Config{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "briefast.example"
	req.Header.Set("X-Forwarded-Proto", "https")

	if got := cfg.BaseURL(req); got != "https://briefast.example" {
		t.Fatalf("BaseURL() = %q, want %q", got, "https://briefast.example")
	}
}

func TestBaseURLFallsBackToPlainConnectionScheme(t *testing.T) {
	cfg := Config{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "127.0.0.1:8600"

	if got := cfg.BaseURL(req); got != "http://127.0.0.1:8600" {
		t.Fatalf("BaseURL() = %q, want %q", got, "http://127.0.0.1:8600")
	}
}

func TestBaseURLFallsBackToTLSConnectionScheme(t *testing.T) {
	cfg := Config{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "briefast.example"
	req.TLS = &tls.ConnectionState{}

	if got := cfg.BaseURL(req); got != "https://briefast.example" {
		t.Fatalf("BaseURL() = %q, want %q", got, "https://briefast.example")
	}
}

// stubReports 讓中繼資料組裝的失敗路徑可測，*store.Store 無法製造查詢錯誤。
type stubReports struct {
	latest  *report.Report
	byDate  map[string]*report.Report
	summary []store.ReportSummary
	count   int
	err     error
}

func (s stubReports) LatestReport(context.Context) (*report.Report, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.latest == nil {
		return nil, store.ErrNotFound
	}
	return s.latest, nil
}

func (s stubReports) ReportByDate(_ context.Context, date string) (*report.Report, error) {
	if s.err != nil {
		return nil, s.err
	}
	if r, ok := s.byDate[date]; ok {
		return r, nil
	}
	return nil, store.ErrNotFound
}

func (s stubReports) ListReports(context.Context, int, int) ([]store.ReportSummary, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.summary, nil
}

func (s stubReports) CountReports(context.Context) (int, error) {
	if s.err != nil {
		return 0, s.err
	}
	return s.count, nil
}

func sampleReports() stubReports {
	latest := &report.Report{
		Date:       "2026-08-21",
		Headline:   "台積電領軍收紅，電子權值撐盤",
		OverviewMD: "## 盤前總覽\n\n台股今日開低走高。",
	}
	older := &report.Report{
		Date:       "2026-08-20",
		Headline:   "金融股走弱拖累大盤",
		OverviewMD: "美股回檔，**金融族群** 承壓。",
	}
	return stubReports{
		latest: latest,
		byDate: map[string]*report.Report{"2026-08-21": latest, "2026-08-20": older},
		summary: []store.ReportSummary{
			{Date: "2026-08-21", Headline: latest.Headline, GeneratedAt: "2026-08-21T00:30:00Z"},
			{Date: "2026-08-20", Headline: older.Headline, GeneratedAt: "2026-08-20T00:30:00Z"},
		},
		count: 2,
	}
}

func metaFor(t *testing.T, reports Reports, target string) pageMeta {
	t.Helper()
	deps := Deps{Reports: reports, Config: Config{SiteURL: "https://briefast.example"}}
	return deps.metaFor(httptest.NewRequest(http.MethodGet, target, nil))
}

func TestMetaForHomeUsesLatestReport(t *testing.T) {
	meta := metaFor(t, sampleReports(), "/")

	if want := "台積電領軍收紅，電子權值撐盤｜Briefast"; meta.Title != want {
		t.Errorf("Title = %q, want %q", meta.Title, want)
	}
	if want := "台股今日開低走高。"; meta.Description != want {
		t.Errorf("Description = %q, want %q", meta.Description, want)
	}
	if want := "https://briefast.example/"; meta.CanonicalURL != want {
		t.Errorf("CanonicalURL = %q, want %q", meta.CanonicalURL, want)
	}
	if meta.OGType != "article" {
		t.Errorf("OGType = %q, want article", meta.OGType)
	}
}

func TestMetaForDatedReportUsesThatReport(t *testing.T) {
	meta := metaFor(t, sampleReports(), "/history/?date=2026-08-20")

	if want := "金融股走弱拖累大盤｜Briefast 2026-08-20"; meta.Title != want {
		t.Errorf("Title = %q, want %q", meta.Title, want)
	}
	if want := "美股回檔，金融族群承壓。"; meta.Description != want {
		t.Errorf("Description = %q, want %q", meta.Description, want)
	}
	if want := "https://briefast.example/history/?date=2026-08-20"; meta.CanonicalURL != want {
		t.Errorf("CanonicalURL = %q, want %q", meta.CanonicalURL, want)
	}
}

func TestMetaForHistoryListUsesListingCopy(t *testing.T) {
	meta := metaFor(t, sampleReports(), "/history/")

	if want := "歷史報告｜Briefast"; meta.Title != want {
		t.Errorf("Title = %q, want %q", meta.Title, want)
	}
	if meta.Description != historyDescription {
		t.Errorf("Description = %q, want %q", meta.Description, historyDescription)
	}
	if want := "https://briefast.example/history/"; meta.CanonicalURL != want {
		t.Errorf("CanonicalURL = %q, want %q", meta.CanonicalURL, want)
	}
	if meta.OGType != "website" {
		t.Errorf("OGType = %q, want website", meta.OGType)
	}
}

func TestMetaDescriptionIsTruncatedWithEllipsis(t *testing.T) {
	long := strings.Repeat("市", 400)
	reports := sampleReports()
	reports.latest = &report.Report{Date: "2026-08-21", Headline: "長篇總覽", OverviewMD: long}

	meta := metaFor(t, reports, "/")

	runes := []rune(meta.Description)
	if len(runes) != descriptionLimit+1 {
		t.Fatalf("description rune count = %d, want %d", len(runes), descriptionLimit+1)
	}
	if runes[len(runes)-1] != '…' {
		t.Errorf("description does not end with an ellipsis: %q", meta.Description)
	}
}

func TestMetaDescriptionFallsBackWhenOverviewIsEmpty(t *testing.T) {
	reports := sampleReports()
	reports.latest = &report.Report{Date: "2026-08-21", Headline: "沒有總覽", OverviewMD: "   \n\n"}

	if got := metaFor(t, reports, "/").Description; got != siteDescription {
		t.Errorf("Description = %q, want site default %q", got, siteDescription)
	}
}

func TestMetaFallsBackToSiteDefaultsWhenLookupFails(t *testing.T) {
	broken := stubReports{err: errors.New("database is locked")}

	for _, target := range []string{"/", "/history/?date=2026-08-20"} {
		meta := metaFor(t, broken, target)
		if meta.Title != siteTitle {
			t.Errorf("%s Title = %q, want %q", target, meta.Title, siteTitle)
		}
		if meta.Description != siteDescription {
			t.Errorf("%s Description = %q, want %q", target, meta.Description, siteDescription)
		}
	}
}

func TestMetaFallsBackWhenDateHasNoReport(t *testing.T) {
	meta := metaFor(t, sampleReports(), "/history/?date=1999-01-01")

	if meta.Title != siteTitle {
		t.Errorf("Title = %q, want %q", meta.Title, siteTitle)
	}
	if want := "https://briefast.example/history/"; meta.CanonicalURL != want {
		t.Errorf("CanonicalURL = %q, want %q", meta.CanonicalURL, want)
	}
}

const shellHTML = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Briefast</title>
<link rel="stylesheet" href="/_syralit/assets/runtime.css">
</head>
<body>
<div id="syralit-root"><main id="syralit-app"><p>連線中…</p></main></div>
</body>
</html>`

func shellHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(shellHTML))
	})
}

func serveThroughMiddleware(t *testing.T, reports Reports, next http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	deps := Deps{Reports: reports, Config: Config{SiteURL: "https://briefast.example"}}
	w := httptest.NewRecorder()
	deps.Middleware(next).ServeHTTP(w, httptest.NewRequest(http.MethodGet, target, nil))
	return w
}

func TestMiddlewareRewritesHeadOfHomePage(t *testing.T) {
	w := serveThroughMiddleware(t, sampleReports(), shellHandler(), "/")
	body := w.Body.String()

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	for _, want := range []string{
		`<title>台積電領軍收紅，電子權值撐盤｜Briefast</title>`,
		`<meta name="description" content="台股今日開低走高。">`,
		`<link rel="canonical" href="https://briefast.example/">`,
		`<meta property="og:type" content="article">`,
		`<meta property="og:site_name" content="Briefast">`,
		`<meta property="og:title" content="台積電領軍收紅，電子權值撐盤｜Briefast">`,
		`<meta property="og:description" content="台股今日開低走高。">`,
		`<meta property="og:url" content="https://briefast.example/">`,
		`<meta property="og:locale" content="zh_TW">`,
		`<meta name="twitter:card" content="summary">`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %s", want)
		}
	}
}

// 語言標示改由 syralit 的 shell 設定提供，中介軟體不得再碰它——否則框架設定與
// 字串替換會互相打架，而且改寫失效時語言標示會跟著消失。
func TestMiddlewareLeavesLangAttributeAlone(t *testing.T) {
	body := serveThroughMiddleware(t, sampleReports(), shellHandler(), "/").Body.String()

	if !strings.Contains(body, `<html lang="en">`) {
		t.Errorf("middleware rewrote the lang attribute: %q", body)
	}
	if !strings.Contains(body, `<title>台積電領軍收紅，電子權值撐盤｜Briefast</title>`) {
		t.Error("title was not rewritten")
	}
}

func TestMiddlewareLeavesBodyUntouched(t *testing.T) {
	w := serveThroughMiddleware(t, sampleReports(), shellHandler(), "/")

	_, originalBody, _ := strings.Cut(shellHTML, "</head>")
	_, rewrittenBody, _ := strings.Cut(w.Body.String(), "</head>")
	if rewrittenBody != originalBody {
		t.Errorf("body after </head> changed:\n got %q\nwant %q", rewrittenBody, originalBody)
	}
}

func TestMiddlewareEscapesMetadataValues(t *testing.T) {
	reports := sampleReports()
	reports.latest = &report.Report{
		Date:       "2026-08-21",
		Headline:   `"權值股" 反攻 & 收紅`,
		OverviewMD: `<script>alert(1)</script> 台股走高。`,
	}

	body := serveThroughMiddleware(t, reports, shellHandler(), "/").Body.String()

	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Error("description was injected without escaping")
	}
	if !strings.Contains(body, "&#34;權值股&#34; 反攻 &amp; 收紅") {
		t.Errorf("headline was not escaped: %q", body)
	}
}

func TestMiddlewareRewritesDatedReportPage(t *testing.T) {
	body := serveThroughMiddleware(t, sampleReports(), shellHandler(), "/history/?date=2026-08-20").Body.String()

	if !strings.Contains(body, `<title>金融股走弱拖累大盤｜Briefast 2026-08-20</title>`) {
		t.Errorf("dated page title missing: %q", body)
	}
	if !strings.Contains(body, `<link rel="canonical" href="https://briefast.example/history/?date=2026-08-20">`) {
		t.Error("dated page canonical missing")
	}
}

func TestMiddlewareSkipsFrameworkPaths(t *testing.T) {
	var wrapped bool
	probe := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, wrapped = w.(*interceptor)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(shellHTML))
	})

	for _, target := range []string{"/_syralit/ws", "/history/_syralit/sse", "/_syralit/assets/runtime.js"} {
		w := serveThroughMiddleware(t, sampleReports(), probe, target)
		if wrapped {
			t.Errorf("%s was wrapped by the interceptor", target)
		}
		if w.Body.String() != shellHTML {
			t.Errorf("%s body was modified", target)
		}
	}
}

func TestMiddlewarePassesThroughNonHTMLResponses(t *testing.T) {
	payload := `{"ok":true}`
	json := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(payload))
	})

	w := serveThroughMiddleware(t, sampleReports(), json, "/")
	if w.Body.String() != payload {
		t.Errorf("body = %q, want %q", w.Body.String(), payload)
	}
}

func TestMiddlewarePassesThroughNonOKResponses(t *testing.T) {
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(shellHTML))
	})

	w := serveThroughMiddleware(t, sampleReports(), notFound, "/")
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
	if w.Body.String() != shellHTML {
		t.Error("404 body was rewritten")
	}
}

func TestMiddlewarePassesThroughHTMLWithoutTitle(t *testing.T) {
	bare := `<!doctype html><html lang="en"><head></head><body>hi</body></html>`
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(bare))
	})

	w := serveThroughMiddleware(t, sampleReports(), handler, "/")
	if w.Body.String() != bare {
		t.Errorf("body = %q, want unchanged", w.Body.String())
	}
}

func TestMiddlewareSetsAccurateContentLength(t *testing.T) {
	w := serveThroughMiddleware(t, sampleReports(), shellHandler(), "/")

	got := w.Header().Get("Content-Length")
	want := strconv.Itoa(w.Body.Len())
	if got != want {
		t.Errorf("Content-Length = %q, want %q", got, want)
	}
}

func TestMiddlewareServesDefaultsWhenLookupFails(t *testing.T) {
	broken := stubReports{err: errors.New("database is locked")}

	w := serveThroughMiddleware(t, broken, shellHandler(), "/")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "<title>"+siteTitle+"</title>") {
		t.Errorf("default title missing: %q", w.Body.String())
	}
}

// flushRecorder 記錄底層 Flush 有沒有被呼叫到。
type flushRecorder struct {
	*httptest.ResponseRecorder
	flushed bool
}

func (f *flushRecorder) Flush() { f.flushed = true; f.ResponseRecorder.Flush() }

func TestMiddlewareForwardsFlushAndStopsBuffering(t *testing.T) {
	streaming := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(shellHTML))
		w.(http.Flusher).Flush()
	})

	recorder := &flushRecorder{ResponseRecorder: httptest.NewRecorder()}
	deps := Deps{Reports: sampleReports(), Config: Config{SiteURL: "https://briefast.example"}}
	deps.Middleware(streaming).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	if !recorder.flushed {
		t.Error("Flush was not forwarded to the underlying writer")
	}
	if recorder.Body.String() != shellHTML {
		t.Error("flushed response was rewritten instead of streamed through")
	}
}

// hijackRecorder 讓 WebSocket 升級路徑可測。
type hijackRecorder struct {
	*httptest.ResponseRecorder
	hijacked bool
}

func (h *hijackRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.hijacked = true
	return nil, nil, nil
}

func TestMiddlewareForwardsHijack(t *testing.T) {
	upgrade := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			t.Error("wrapped writer does not support hijacking")
			return
		}
		if _, _, err := hijacker.Hijack(); err != nil {
			t.Errorf("Hijack() error = %v", err)
		}
	})

	recorder := &hijackRecorder{ResponseRecorder: httptest.NewRecorder()}
	deps := Deps{Reports: sampleReports(), Config: Config{SiteURL: "https://briefast.example"}}
	deps.Middleware(upgrade).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	if !recorder.hijacked {
		t.Error("Hijack was not forwarded to the underlying writer")
	}
}

func serveEndpoint(t *testing.T, reports Reports, handler func(Deps) http.Handler, target string) *httptest.ResponseRecorder {
	t.Helper()
	deps := Deps{Reports: reports, Config: Config{SiteURL: "https://briefast.example"}}
	w := httptest.NewRecorder()
	handler(deps).ServeHTTP(w, httptest.NewRequest(http.MethodGet, target, nil))
	return w
}

func TestRobotsHandlerDisallowsPrivatePathsAndPointsAtSitemap(t *testing.T) {
	w := serveEndpoint(t, sampleReports(), Deps.RobotsHandler, "/robots.txt")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("Content-Type = %q, want text/plain", ct)
	}
	for _, want := range []string{
		"User-agent: *",
		"Disallow: /admin/",
		"Disallow: /api/",
		"Sitemap: https://briefast.example/sitemap.xml",
	} {
		if !strings.Contains(w.Body.String(), want) {
			t.Errorf("robots.txt missing %q\ngot:\n%s", want, w.Body.String())
		}
	}
}

type sitemapDoc struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod"`
	} `xml:"url"`
}

func parseSitemap(t *testing.T, body string) sitemapDoc {
	t.Helper()
	var doc sitemapDoc
	if err := xml.Unmarshal([]byte(body), &doc); err != nil {
		t.Fatalf("sitemap is not valid XML: %v\n%s", err, body)
	}
	return doc
}

func TestSitemapListsEveryStoredReport(t *testing.T) {
	w := serveEndpoint(t, sampleReports(), Deps.SitemapHandler, "/sitemap.xml")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/xml") {
		t.Errorf("Content-Type = %q, want application/xml", ct)
	}

	doc := parseSitemap(t, w.Body.String())
	locations := map[string]string{}
	for _, u := range doc.URLs {
		locations[u.Loc] = u.LastMod
	}
	for _, want := range []string{
		"https://briefast.example/",
		"https://briefast.example/history/",
		"https://briefast.example/history/?date=2026-08-21",
		"https://briefast.example/history/?date=2026-08-20",
	} {
		if _, ok := locations[want]; !ok {
			t.Errorf("sitemap missing %s", want)
		}
	}
	if got := locations["https://briefast.example/history/?date=2026-08-21"]; got != "2026-08-21T00:30:00Z" {
		t.Errorf("lastmod = %q, want the report generated time", got)
	}
}

func TestSitemapWithNoReportsStillListsSiteEntries(t *testing.T) {
	w := serveEndpoint(t, stubReports{}, Deps.SitemapHandler, "/sitemap.xml")

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	doc := parseSitemap(t, w.Body.String())
	if len(doc.URLs) != 2 {
		t.Fatalf("url count = %d, want 2", len(doc.URLs))
	}
	if doc.URLs[0].Loc != "https://briefast.example/" || doc.URLs[1].Loc != "https://briefast.example/history/" {
		t.Errorf("unexpected urls: %+v", doc.URLs)
	}
}

func TestSitemapReportsQueryFailure(t *testing.T) {
	w := serveEndpoint(t, stubReports{err: errors.New("database is locked")}, Deps.SitemapHandler, "/sitemap.xml")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	if strings.Contains(w.Body.String(), "<urlset") {
		t.Error("a partial sitemap was written on failure")
	}
}
