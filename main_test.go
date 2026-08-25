package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/site"
	"github.com/HazelnutParadise/briefast/internal/store"
	sy "github.com/HazelnutParadise/syralit"
)

func TestCustomMuxMountsSyralitPagesAndReportAPI(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "briefast.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	handler := newHandler(s)

	for _, path := range []string{"/", "/history/", "/admin/"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("GET %s status = %d", path, w.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTemporaryRedirect || w.Header().Get("Location") != "/history/" {
		t.Fatalf("GET /history status = %d, location = %q", w.Code, w.Header().Get("Location"))
	}

	req = httptest.NewRequest(http.MethodGet, "/api/report", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET /api/report status = %d", w.Code)
	}

	// 唯讀端點需 Bearer key，未帶 key 應為 401 而非落到首頁或 404。
	req = httptest.NewRequest(http.MethodGet, "/api/report/2026-08-07", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/report/{date} status = %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/report/2026-08-07", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST /api/report/{date} status = %d", w.Code)
	}
}

func TestNewsHeadlinesReportEndToEnd(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "briefast.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	key, err := s.CreateAPIKey(context.Background(), "integration", "news-headlines-token")
	if err != nil {
		t.Fatal(err)
	}
	handler := newHandler(s)
	value := report.Report{
		Date:       "2026-08-09",
		Headline:   "標題結構整合驗收",
		OverviewMD: "盤前總覽內容",
		WatchMD:    "- 今日觀察",
		Calls:      report.Calls{},
		Industries: []report.Industry{{
			Name: "科技",
			Events: []report.IndustryEvent{
				{Headline: "記憶體合約價續漲，上游產品組合改善", SummaryMD: "記憶體事件內文"},
				{Headline: "AI 追單推升封測利用率，營收動能轉強", SummaryMD: "封測事件內文"},
			},
			WatchMD: "- 追蹤下一次合約價公告",
		}},
		StockNews: []report.StockNews{{
			Symbol: "2330", Name: "台積電", Call: report.CallNone,
			Headline: "先進製程時程提前，訂單能見度升高", SummaryMD: "個股事件內文", WatchMD: "- 追蹤量產時程",
			Sources: []report.Source{{Title: "測試來源", URL: "https://example.com/news"}},
		}},
		GeneratedAt: "2026-08-09T07:50:00+08:00",
	}

	postReport := func(payload report.Report) *httptest.ResponseRecorder {
		t.Helper()
		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/report", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+key.Token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		return w
	}

	if w := postReport(value); w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"ok":true`) {
		t.Fatalf("valid POST status = %d, body = %s", w.Code, w.Body.String())
	}
	app := site.New(s)
	at := sy.NewAppTest(app.Home)
	at.Run()
	nodes := at.FindAll("html")
	if len(nodes) != 1 {
		t.Fatalf("html nodes = %d, want 1", len(nodes))
	}
	rendered, _ := nodes[0].Props["html"].(string)
	for _, want := range []string{
		"記憶體合約價續漲，上游產品組合改善", "記憶體事件內文",
		"AI 追單推升封測利用率，營收動能轉強", "封測事件內文",
		"先進製程時程提前，訂單能見度升高", "個股事件內文",
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("rendered report missing %q", want)
		}
	}

	invalid := value
	invalid.Industries = append([]report.Industry(nil), value.Industries...)
	invalid.Industries[0].Events = append([]report.IndustryEvent(nil), value.Industries[0].Events...)
	invalid.Industries[0].Events[1].Headline = " "
	w := postReport(invalid)
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "industries[0].events[1].headline") {
		t.Fatalf("invalid POST status = %d, body = %s", w.Code, w.Body.String())
	}
	count, err := s.CountReports(context.Background())
	if err != nil || count != 1 {
		t.Fatalf("report count after rejected POST = %d, err = %v, want 1", count, err)
	}
	stored, err := s.ReportByDate(context.Background(), value.Date)
	if err != nil || stored.Industries[0].Events[1].Headline != value.Industries[0].Events[1].Headline {
		t.Fatalf("rejected POST changed stored report: stored=%+v, err=%v", stored, err)
	}
	logs, err := s.ListUpdateLogs(context.Background(), 10)
	if err != nil || len(logs) != 2 || logs[0].Action != "ingest_rejected_schema" {
		t.Fatalf("update logs = %+v, err = %v", logs, err)
	}
}

func TestCustomMuxServesCrawlerMetadata(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "briefast.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	handler := newHandler(s)

	get := func(target string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodGet, target, nil)
		req.Host = "briefast.example"
		req.Header.Set("X-Forwarded-Proto", "https")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		return w
	}

	robots := get("/robots.txt")
	if robots.Code != http.StatusOK {
		t.Fatalf("GET /robots.txt status = %d", robots.Code)
	}
	if !strings.Contains(robots.Body.String(), "Sitemap: https://briefast.example/sitemap.xml") {
		t.Errorf("robots.txt sitemap line missing: %s", robots.Body.String())
	}

	sitemap := get("/sitemap.xml")
	if sitemap.Code != http.StatusOK {
		t.Fatalf("GET /sitemap.xml status = %d", sitemap.Code)
	}
	if !strings.Contains(sitemap.Body.String(), "<loc>https://briefast.example/</loc>") {
		t.Errorf("sitemap home entry missing: %s", sitemap.Body.String())
	}

	// 公開頁的 head 要被改寫，後台不受影響。
	home := get("/").Body.String()
	for _, want := range []string{`lang="zh-Hant-TW"`, `rel="canonical"`, `property="og:title"`} {
		if !strings.Contains(home, want) {
			t.Errorf("home page missing %s", want)
		}
	}
	// 後台不經 SEO 中介軟體，語言標示由 syralit.toml 的 lang 提供，所以兩件事要
	// 同時成立：沒有中繼資料標籤，但仍有正確的語言標示。
	adminPage := get("/admin/").Body.String()
	if strings.Contains(adminPage, `rel="canonical"`) {
		t.Error("admin page was rewritten by the SEO middleware")
	}
	if !strings.Contains(adminPage, `lang="zh-Hant-TW"`) {
		t.Errorf("admin page missing the shell language attribute: %.200s", adminPage)
	}
}

// 位址計算過去完全沒被測試覆蓋——所有測試都以 httptest 直接呼叫 newHandler，
// 從不經過 ListenAndServe，所以綁定位址退化成 ":0" 時沒有任何測試失敗。
func TestListenAddressComesFromResolvedConfig(t *testing.T) {
	got := listenAddr(appConfig())

	if want := "0.0.0.0:8600"; got != want {
		t.Fatalf("listenAddr() = %q, want %q", got, want)
	}
}
