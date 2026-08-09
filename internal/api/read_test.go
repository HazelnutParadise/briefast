package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
)

// readRequest 走 mux 以確保 {date} path value 與正式路由一致。
func readRequest(t *testing.T, s *store.Store, token, date string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	mux.Handle("/api/report/{date}", NewReadHandler(s))
	req := httptest.NewRequest(http.MethodGet, "/api/report/"+date, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

func seedReport(t *testing.T, s *store.Store) report.Report {
	t.Helper()
	value := apiReport()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertReport(context.Background(), value, payload); err != nil {
		t.Fatal(err)
	}
	return value
}

func countLogsWithAction(t *testing.T, s *store.Store, action string) int {
	t.Helper()
	logs, err := s.ListUpdateLogs(context.Background(), 100)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, entry := range logs {
		if entry.Action == action {
			n++
		}
	}
	return n
}

func TestReadReturnsStoredReport(t *testing.T) {
	s, _, _, _ := setupHandler(t)
	want := seedReport(t, s)

	w := readRequest(t, s, "valid-token", want.Date)
	if w.Code != http.StatusOK {
		t.Fatalf("狀態碼 = %d, want 200, body=%s", w.Code, w.Body.String())
	}
	var got report.Report
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Date != want.Date || got.Headline != want.Headline {
		t.Fatalf("報告基本欄位不符: %+v", got)
	}
	if len(got.Industries) != 1 || got.Industries[0].WatchMD != want.Industries[0].WatchMD {
		t.Fatalf("industries watch_md 未原樣回傳: %+v", got.Industries)
	}
	if len(got.StockNews) != 1 || got.StockNews[0].WatchMD != want.StockNews[0].WatchMD {
		t.Fatalf("stock_news watch_md 未原樣回傳: %+v", got.StockNews)
	}
	if n := countLogsWithAction(t, s, "read_rejected_auth"); n != 0 {
		t.Fatalf("成功讀取不應留下拒絕紀錄, got %d", n)
	}
}

func TestReadRejectsRevokedKey(t *testing.T) {
	s, key, _, _ := setupHandler(t)
	seedReport(t, s)
	if err := s.RevokeAPIKey(context.Background(), key.ID); err != nil {
		t.Fatal(err)
	}

	w := readRequest(t, s, "valid-token", "2026-08-07")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("狀態碼 = %d, want 401", w.Code)
	}
	if n := countLogsWithAction(t, s, "read_rejected_auth"); n != 1 {
		t.Fatalf("read_rejected_auth 筆數 = %d, want 1", n)
	}
}

func TestReadRejectsMissingHeader(t *testing.T) {
	s, _, _, _ := setupHandler(t)
	seedReport(t, s)

	w := readRequest(t, s, "", "2026-08-07")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("狀態碼 = %d, want 401", w.Code)
	}
	if n := countLogsWithAction(t, s, "read_rejected_auth"); n != 1 {
		t.Fatalf("read_rejected_auth 筆數 = %d, want 1", n)
	}
}

func TestReadRejectsMalformedDate(t *testing.T) {
	s, _, _, _ := setupHandler(t)
	seedReport(t, s)

	for _, date := range []string{"2026-8-7", "2026-02-30", "not-a-date"} {
		w := readRequest(t, s, "valid-token", date)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("date=%s 狀態碼 = %d, want 400", date, w.Code)
		}
	}
}

func TestReadMissingReportReturns404(t *testing.T) {
	s, _, _, _ := setupHandler(t)
	seedReport(t, s)

	w := readRequest(t, s, "valid-token", "2026-08-08")
	if w.Code != http.StatusNotFound {
		t.Fatalf("狀態碼 = %d, want 404", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["ok"] != false {
		t.Fatalf("body ok = %v, want false", body["ok"])
	}
	if n := countLogsWithAction(t, s, "read_rejected_auth"); n != 0 {
		t.Fatalf("查無報告不應寫入拒絕紀錄, got %d", n)
	}
}
