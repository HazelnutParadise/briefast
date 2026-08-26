package seo

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
	sy "github.com/HazelnutParadise/syralit"
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

func metaFor(t *testing.T, reports Reports, kind PageKind, target string) pageMeta {
	t.Helper()
	deps := Deps{Reports: reports, Config: Config{SiteURL: "https://briefast.example"}}
	return deps.metaFor(kind, httptest.NewRequest(http.MethodGet, target, nil))
}

func TestMetaForHomeUsesLatestReport(t *testing.T) {
	meta := metaFor(t, sampleReports(), PageHome, "/")

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
	meta := metaFor(t, sampleReports(), PageHistory, "/?date=2026-08-20")

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
	meta := metaFor(t, sampleReports(), PageHistory, "/")

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

	meta := metaFor(t, reports, PageHome, "/")

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

	if got := metaFor(t, reports, PageHome, "/").Description; got != siteDescription {
		t.Errorf("Description = %q, want site default %q", got, siteDescription)
	}
}

func TestMetaFallsBackToSiteDefaultsWhenLookupFails(t *testing.T) {
	broken := stubReports{err: errors.New("database is locked")}

	for _, tc := range []struct {
		kind   PageKind
		target string
	}{{PageHome, "/"}, {PageHistory, "/?date=2026-08-20"}} {
		kind, target := tc.kind, tc.target
		meta := metaFor(t, broken, kind, target)
		if meta.Title != siteTitle {
			t.Errorf("%s Title = %q, want %q", target, meta.Title, siteTitle)
		}
		if meta.Description != siteDescription {
			t.Errorf("%s Description = %q, want %q", target, meta.Description, siteDescription)
		}
	}
}

func TestMetaFallsBackWhenDateHasNoReport(t *testing.T) {
	meta := metaFor(t, sampleReports(), PageHistory, "/?date=1999-01-01")

	if meta.Title != siteTitle {
		t.Errorf("Title = %q, want %q", meta.Title, siteTitle)
	}
	if want := "https://briefast.example/history/"; meta.CanonicalURL != want {
		t.Errorf("CanonicalURL = %q, want %q", meta.CanonicalURL, want)
	}
}

// documentFor 走的是正式路徑：依頁面種類建立 DocumentFunc，再以請求呼叫它。
func documentFor(t *testing.T, reports Reports, kind PageKind, target string) sy.Document {
	t.Helper()
	deps := Deps{Reports: reports, Config: Config{SiteURL: "https://briefast.example"}}
	return deps.DocumentFunc(kind)(httptest.NewRequest(http.MethodGet, target, nil))
}

func TestDocumentForHomeCarriesLatestReport(t *testing.T) {
	doc := documentFor(t, sampleReports(), PageHome, "/")

	if want := "台積電領軍收紅，電子權值撐盤｜Briefast"; doc.Title != want {
		t.Errorf("Title = %q, want %q", doc.Title, want)
	}
	for _, want := range []string{
		`<meta name="description" content="台股今日開低走高。">`,
		`<link rel="canonical" href="https://briefast.example/">`,
		`<meta property="og:type" content="article">`,
		`<meta property="og:site_name" content="Briefast">`,
		`<meta property="og:title" content="台積電領軍收紅，電子權值撐盤｜Briefast">`,
		`<meta property="og:url" content="https://briefast.example/">`,
		`<meta property="og:locale" content="zh_TW">`,
		`<meta name="twitter:card" content="summary">`,
	} {
		if !strings.Contains(doc.HeadHTML, want) {
			t.Errorf("HeadHTML missing %s", want)
		}
	}
	if strings.Contains(doc.HeadHTML, "<title>") {
		t.Error("HeadHTML must not carry a title element; Document.Title owns it")
	}
}

func TestDocumentForDatedReport(t *testing.T) {
	doc := documentFor(t, sampleReports(), PageHistory, "/?date=2026-08-20")

	if want := "金融股走弱拖累大盤｜Briefast 2026-08-20"; doc.Title != want {
		t.Errorf("Title = %q, want %q", doc.Title, want)
	}
	if want := `<link rel="canonical" href="https://briefast.example/history/?date=2026-08-20">`; !strings.Contains(doc.HeadHTML, want) {
		t.Errorf("canonical missing: %q", doc.HeadHTML)
	}
}

func TestDocumentForHistoryList(t *testing.T) {
	doc := documentFor(t, sampleReports(), PageHistory, "/")

	if want := "歷史報告｜Briefast"; doc.Title != want {
		t.Errorf("Title = %q, want %q", doc.Title, want)
	}
	if !strings.Contains(doc.HeadHTML, `<meta property="og:type" content="website">`) {
		t.Error("history listing should be og:type website")
	}
}

func TestDocumentEscapesMetadataValues(t *testing.T) {
	reports := sampleReports()
	reports.latest = &report.Report{
		Date:       "2026-08-21",
		Headline:   `"權值股" 反攻 & 收紅`,
		OverviewMD: `<script>alert(1)</script> 台股走高。`,
	}

	doc := documentFor(t, reports, PageHome, "/")

	if strings.Contains(doc.HeadHTML, "<script>alert(1)</script>") {
		t.Error("description was injected without escaping")
	}
	if !strings.Contains(doc.HeadHTML, "&#34;權值股&#34; 反攻 &amp; 收紅") {
		t.Errorf("headline was not escaped: %q", doc.HeadHTML)
	}
}

func TestDocumentFallsBackWhenLookupFails(t *testing.T) {
	doc := documentFor(t, stubReports{err: errors.New("database is locked")}, PageHome, "/")

	if doc.Title != siteTitle {
		t.Errorf("Title = %q, want %q", doc.Title, siteTitle)
	}
}

func TestDocumentLeavesLangToTheShellConfig(t *testing.T) {
	doc := documentFor(t, sampleReports(), PageHome, "/")

	if doc.Lang != "" || doc.Dir != "" {
		t.Errorf("Lang/Dir should stay empty so syralit.toml applies, got %q/%q", doc.Lang, doc.Dir)
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
