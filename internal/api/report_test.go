package api

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
	"github.com/HazelnutParadise/briefast/internal/store"
)

type countNotifier struct{ count int }

func (n *countNotifier) Notify() { n.count++ }

func setupHandler(t *testing.T) (*store.Store, store.APIKey, *countNotifier, http.Handler) {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "briefast.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	key, err := s.CreateAPIKey(context.Background(), "cowork-daily", "valid-token")
	if err != nil {
		t.Fatal(err)
	}
	n := &countNotifier{}
	return s, key, n, NewReportHandler(s, n)
}

func apiReport() report.Report {
	return report.Report{
		Date: "2026-08-07", Headline: "原始標題", OverviewMD: "總覽", WatchMD: "觀察",
		Calls: report.Calls{ShortBull: []report.CallEntry{{Symbol: "2330", Name: "台積電", Reason: "理由"}}},
		StockNews: []report.StockNews{{Symbol: "2330", Name: "台積電", Call: report.CallShortBull,
			SummaryMD: "摘要", Sources: []report.Source{{Title: "來源", URL: "https://example.com"}}}},
		GeneratedAt: "2026-08-07T07:50:00+08:00",
	}
}

func request(t *testing.T, handler http.Handler, token string, value any) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/report", bytes.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestReportHandlerRejectsMissingAndUnknownKeysWithLogs(t *testing.T) {
	for _, token := range []string{"", "unknown-token"} {
		t.Run(token, func(t *testing.T) {
			s, _, notifier, handler := setupHandler(t)
			w := request(t, handler, token, apiReport())
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
			}
			count, _ := s.CountReports(context.Background())
			if count != 0 || notifier.count != 0 {
				t.Fatalf("reports = %d, notifications = %d", count, notifier.count)
			}
			logs, err := s.ListUpdateLogs(context.Background(), 10)
			if err != nil || len(logs) != 1 || logs[0].Action != "ingest_rejected_auth" {
				t.Fatalf("logs = %+v, err = %v", logs, err)
			}
		})
	}
}

func TestReportHandlerRejectsRevokedKey(t *testing.T) {
	s, key, _, handler := setupHandler(t)
	if err := s.RevokeAPIKey(context.Background(), key.ID); err != nil {
		t.Fatal(err)
	}
	w := request(t, handler, key.Token, apiReport())
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
}

func TestReportHandlerRejectsEverySchemaErrorAtomically(t *testing.T) {
	s, key, _, handler := setupHandler(t)
	value := apiReport()
	value.Date = "bad"
	value.Headline = ""
	w := request(t, handler, key.Token, value)
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "date") || !strings.Contains(w.Body.String(), "headline") {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	count, _ := s.CountReports(context.Background())
	if count != 0 {
		t.Fatalf("reports = %d, want 0", count)
	}
	logs, _ := s.ListUpdateLogs(context.Background(), 10)
	if len(logs) != 1 || logs[0].Action != "ingest_rejected_schema" || !strings.Contains(logs[0].Detail, "headline") {
		t.Fatalf("logs = %+v", logs)
	}
}

func TestReportHandlerAcceptsAndOverwritesSameDate(t *testing.T) {
	s, key, notifier, handler := setupHandler(t)
	first := apiReport()
	w := request(t, handler, key.Token, first)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"ok":true`) {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	second := first
	second.Headline = "覆寫標題"
	w = request(t, handler, key.Token, second)
	if w.Code != http.StatusOK {
		t.Fatalf("second status = %d, body = %s", w.Code, w.Body.String())
	}
	stored, err := s.ReportByDate(context.Background(), first.Date)
	if err != nil || stored.Headline != "覆寫標題" {
		t.Fatalf("stored = %+v, err = %v", stored, err)
	}
	count, _ := s.CountReports(context.Background())
	logs, _ := s.ListUpdateLogs(context.Background(), 10)
	if count != 1 || len(logs) != 2 || notifier.count != 2 || logs[0].KeyName != key.Name {
		t.Fatalf("reports = %d, logs = %+v, notifications = %d", count, logs, notifier.count)
	}
}
