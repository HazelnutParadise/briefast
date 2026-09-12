package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
)

func apiConference(symbol, heldOn, prediction string) report.Conference {
	return report.Conference{
		Symbol: symbol, Name: "公司" + symbol, Market: report.MarketTWSE,
		HeldOn: heldOn, HeldAt: "14:00", Venue: "線上法說會", AnnouncedOn: "2026-09-10",
		Prediction: prediction, Headline: "標題", SummaryMD: "依據", WatchMD: "- 觀察",
		Sources: []report.Source{{Title: "來源", URL: "https://example.com"}}, GeneratedAt: "2026-09-10T07:50:00+08:00",
	}
}

func conferenceMux(t *testing.T) (*store.Store, store.APIKey, *countNotifier, http.Handler) {
	t.Helper()
	s, key, n, _ := setupHandler(t)
	h := NewConferenceHandler(s, n).WithNow(func() time.Time {
		return time.Date(2026, 9, 14, 7, 50, 0, 0, time.FixedZone("Asia/Taipei", 8*3600))
	})
	mux := http.NewServeMux()
	mux.Handle("/api/conference", h.Ingest())
	mux.Handle("/api/conference/settle", h.Settle())
	mux.Handle("/api/conferences", h.List())
	return s, key, n, mux
}

func call(t *testing.T, handler http.Handler, method, path, token string, value any) *httptest.ResponseRecorder {
	t.Helper()
	var body []byte
	if value != nil {
		var err error
		if body, err = json.Marshal(value); err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode body %s: %v", w.Body.String(), err)
	}
	return out
}

func TestConferenceIngestPersistsLogsAndNotifies(t *testing.T) {
	s, key, notifier, mux := conferenceMux(t)
	w := call(t, mux, http.MethodPost, "/api/conference", key.Token, apiConference("2330", "2026-10-16", report.PredictionBull))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	body := decodeBody(t, w)
	if body["ok"] != true || body["symbol"] != "2330" || body["held_on"] != "2026-10-16" {
		t.Fatalf("body = %v", body)
	}
	if notifier.count != 1 {
		t.Fatalf("notifications = %d, want 1", notifier.count)
	}
	stored, err := s.ConferenceByKey(context.Background(), "2330", "2026-10-16")
	if err != nil || stored.Brief.Headline != "標題" {
		t.Fatalf("stored = %+v, err = %v", stored, err)
	}
	logs, _ := s.ListUpdateLogs(context.Background(), 10)
	if len(logs) != 1 || logs[0].Action != "conference_ok" || logs[0].ReportDate != "2026-10-16" || logs[0].KeyName != "cowork-daily" {
		t.Fatalf("logs = %+v", logs)
	}
}

func TestConferenceIngestRejectsInvalidAtomically(t *testing.T) {
	s, key, notifier, mux := conferenceMux(t)
	bad := apiConference("2330", "2026-10-16", "up")
	bad.Headline = ""
	w := call(t, mux, http.MethodPost, "/api/conference", key.Token, bad)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	body := decodeBody(t, w)
	errs, _ := body["errors"].([]any)
	joined := ""
	for _, e := range errs {
		joined += e.(string) + "\n"
	}
	if len(errs) != 2 || !strings.Contains(joined, "prediction") || !strings.Contains(joined, "headline") {
		t.Fatalf("errors = %v", errs)
	}
	if _, err := s.ConferenceByKey(context.Background(), "2330", "2026-10-16"); err != store.ErrNotFound {
		t.Fatalf("row should not exist, err = %v", err)
	}
	if notifier.count != 0 || countLogsWithAction(t, s, "conference_rejected_schema") != 1 {
		t.Fatalf("notifications = %d, logs = %d", notifier.count, countLogsWithAction(t, s, "conference_rejected_schema"))
	}

	w = call(t, mux, http.MethodPost, "/api/conference", key.Token, map[string]any{"symbol": "2330", "unknown": 1})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("unknown field status = %d", w.Code)
	}
}

func TestConferenceEndpointsRejectMissingKeyWithLogs(t *testing.T) {
	s, _, _, mux := conferenceMux(t)
	cases := []struct{ method, path, action string }{
		{http.MethodPost, "/api/conference", "conference_rejected_auth"},
		{http.MethodPost, "/api/conference/settle", "settle_rejected_auth"},
		{http.MethodGet, "/api/conferences", "list_rejected_auth"},
	}
	for _, tc := range cases {
		w := call(t, mux, tc.method, tc.path, "", nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s status = %d", tc.method, tc.path, w.Code)
		}
		if countLogsWithAction(t, s, tc.action) != 1 {
			t.Fatalf("%s log missing", tc.action)
		}
	}
	w := call(t, mux, http.MethodGet, "/api/conference", "", nil)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET /api/conference status = %d", w.Code)
	}
}

func TestConferenceSettlementOutcomesAndRepostKeepsSettlement(t *testing.T) {
	s, key, _, mux := conferenceMux(t)
	cases := []struct {
		symbol, prediction string
		pre, post          float64
		want               string
	}{
		{"1", report.PredictionBull, 1080, 1095, report.OutcomeHit},
		{"2", report.PredictionBull, 1080, 1080, report.OutcomeMiss},
		{"3", report.PredictionBull, 1080, 1060, report.OutcomeMiss},
		{"4", report.PredictionBear, 500, 480, report.OutcomeHit},
		{"5", report.PredictionBear, 500, 510, report.OutcomeMiss},
		{"6", report.PredictionNone, 500, 520, ""},
	}
	for _, tc := range cases {
		if w := call(t, mux, http.MethodPost, "/api/conference", key.Token, apiConference(tc.symbol, "2026-09-11", tc.prediction)); w.Code != http.StatusOK {
			t.Fatalf("ingest %s status = %d", tc.symbol, w.Code)
		}
		w := call(t, mux, http.MethodPost, "/api/conference/settle", key.Token, map[string]any{
			"symbol": tc.symbol, "held_on": "2026-09-11", "pre_close_date": "2026-09-10", "pre_close": tc.pre, "post_close_date": "2026-09-14", "post_close": tc.post,
		})
		if w.Code != http.StatusOK {
			t.Fatalf("settle %s status = %d, body = %s", tc.symbol, w.Code, w.Body.String())
		}
		if body := decodeBody(t, w); body["outcome"] != tc.want {
			t.Fatalf("settle %s outcome = %v, want %q", tc.symbol, body["outcome"], tc.want)
		}
	}
	if countLogsWithAction(t, s, "settle_ok") != 6 {
		t.Fatalf("settle_ok logs = %d", countLogsWithAction(t, s, "settle_ok"))
	}

	// Re-posting the brief must not erase the settlement.
	updated := apiConference("1", "2026-09-11", report.PredictionBull)
	updated.SummaryMD = "新依據"
	if w := call(t, mux, http.MethodPost, "/api/conference", key.Token, updated); w.Code != http.StatusOK {
		t.Fatalf("re-post status = %d", w.Code)
	}
	stored, err := s.ConferenceByKey(context.Background(), "1", "2026-09-11")
	if err != nil || stored.Brief.SummaryMD != "新依據" || stored.Settlement == nil || stored.Settlement.Outcome != report.OutcomeHit || stored.Settlement.PreClose != 1080 {
		t.Fatalf("after re-post = %+v, err = %v", stored, err)
	}
}

func TestConferenceSettlementRejectsOversizedBody(t *testing.T) {
	_, key, _, mux := conferenceMux(t)
	huge := bytes.Repeat([]byte(" "), maxReportBytes+2)
	req := httptest.NewRequest(http.MethodPost, "/api/conference/settle", bytes.NewReader(append([]byte(`{"symbol":"2330"`), huge...)))
	req.Header.Set("Authorization", "Bearer "+key.Token)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("oversized settle status = %d", w.Code)
	}
}

func TestConferenceSettlementRejectsBadDatesAndUnknownPair(t *testing.T) {
	s, key, _, mux := conferenceMux(t)
	call(t, mux, http.MethodPost, "/api/conference", key.Token, apiConference("2330", "2026-09-11", report.PredictionBull))
	w := call(t, mux, http.MethodPost, "/api/conference/settle", key.Token, map[string]any{
		"symbol": "2330", "held_on": "2026-09-11", "pre_close_date": "2026-09-11", "pre_close": 100, "post_close_date": "2026-09-14", "post_close": 101,
	})
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "pre_close_date") {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	stored, _ := s.ConferenceByKey(context.Background(), "2330", "2026-09-11")
	if stored.Settlement != nil || countLogsWithAction(t, s, "settle_rejected_schema") != 1 {
		t.Fatalf("settlement stored on rejected request: %+v", stored.Settlement)
	}
	w = call(t, mux, http.MethodPost, "/api/conference/settle", key.Token, map[string]any{
		"symbol": "9999", "held_on": "2026-09-11", "pre_close_date": "2026-09-10", "pre_close": 100, "post_close_date": "2026-09-14", "post_close": 101,
	})
	if w.Code != http.StatusNotFound || decodeBody(t, w)["ok"] != false {
		t.Fatalf("unknown pair status = %d, body = %s", w.Code, w.Body.String())
	}
}

func TestConferenceListSplitsUpcomingFromPending(t *testing.T) {
	s, key, _, mux := conferenceMux(t)
	w := call(t, mux, http.MethodGet, "/api/conferences", key.Token, nil)
	if w.Code != http.StatusOK || w.Body.String() != "{\"pending_settlement\":[],\"upcoming\":[]}\n" {
		t.Fatalf("empty list status = %d, body = %s", w.Code, w.Body.String())
	}
	for _, c := range []report.Conference{
		apiConference("1314", "2026-09-11", report.PredictionBull),
		apiConference("2330", "2026-09-12", report.PredictionBear),
		apiConference("3443", "2026-09-14", report.PredictionNone),
		apiConference("1101", "2026-09-16", report.PredictionBull),
	} {
		call(t, mux, http.MethodPost, "/api/conference", key.Token, c)
	}
	call(t, mux, http.MethodPost, "/api/conference/settle", key.Token, map[string]any{
		"symbol": "2330", "held_on": "2026-09-12", "pre_close_date": "2026-09-11", "pre_close": 500, "post_close_date": "2026-09-14", "post_close": 480,
	})
	w = call(t, mux, http.MethodGet, "/api/conferences", key.Token, nil)
	var body struct {
		Upcoming []conferenceListItem `json:"upcoming"`
		Pending  []conferenceListItem `json:"pending_settlement"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Upcoming) != 2 || body.Upcoming[0].Symbol != "3443" || body.Upcoming[1].Symbol != "1101" || body.Upcoming[0].AnnouncedOn != "2026-09-10" || body.Upcoming[0].Venue != "線上法說會" {
		t.Fatalf("upcoming = %+v", body.Upcoming)
	}
	if len(body.Pending) != 1 || body.Pending[0].Symbol != "1314" || body.Pending[0].AnnouncedOn != "" || body.Pending[0].Venue != "" {
		t.Fatalf("pending = %+v", body.Pending)
	}
	if strings.Contains(w.Body.String(), "summary_md") {
		t.Fatalf("listing leaks brief content: %s", w.Body.String())
	}
	if countLogsWithAction(t, s, "list_rejected_auth") != 0 {
		t.Fatal("successful listing must not log")
	}
}
