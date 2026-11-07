package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func assertPointer(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	if got := w.Header().Get("Link"); got != GuideLink {
		t.Fatalf("Link 標頭 = %q, want %q", got, GuideLink)
	}
	body := decodeBody(t, w)
	if body["guide"] != GuidePointer {
		t.Fatalf("guide 欄位 = %v, want %q", body["guide"], GuidePointer)
	}
	return body
}

func TestEveryResponseKindCarriesPointer(t *testing.T) {
	s, key, _, handler := setupHandler(t)
	seedReport(t, s)
	conferences := NewConferenceHandler(s, nil)

	cases := map[string]func() *httptest.ResponseRecorder{
		"405": func() *httptest.ResponseRecorder {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/report", nil))
			return w
		},
		"401":        func() *httptest.ResponseRecorder { return request(t, handler, "", apiReport()) },
		"400":        func() *httptest.ResponseRecorder { return request(t, handler, key.Token, map[string]any{"date": "x"}) },
		"200 ingest": func() *httptest.ResponseRecorder { return request(t, handler, key.Token, apiReport()) },
		"200 list": func() *httptest.ResponseRecorder {
			req := httptest.NewRequest(http.MethodGet, "/api/conferences", nil)
			req.Header.Set("Authorization", "Bearer "+key.Token)
			w := httptest.NewRecorder()
			conferences.List().ServeHTTP(w, req)
			return w
		},
		"404 read": func() *httptest.ResponseRecorder { return readRequest(t, s, key.Token, "2026-01-01") },
		"500": func() *httptest.ResponseRecorder {
			w := httptest.NewRecorder()
			writeServerError(w, http.ErrAbortHandler)
			return w
		},
	}
	for name, run := range cases {
		t.Run(name, func(t *testing.T) { assertPointer(t, run()) })
	}
}

func TestUnknownAPIPathAnswersJSON(t *testing.T) {
	w := httptest.NewRecorder()
	NotFoundHandler().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/reports", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("狀態碼 = %d, want 404", w.Code)
	}
	if !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("Content-Type = %q", w.Header().Get("Content-Type"))
	}
	body := assertPointer(t, w)
	if body["ok"] != false || !strings.Contains(w.Body.String(), "/api/reports") {
		t.Fatalf("404 內容不完整: %s", w.Body.String())
	}
}

func TestReadWrapsReportAndCanBeResubmitted(t *testing.T) {
	s, key, _, handler := setupHandler(t)
	want := seedReport(t, s)

	w := readRequest(t, s, key.Token, want.Date)
	if w.Code != http.StatusOK {
		t.Fatalf("狀態碼 = %d, body=%s", w.Code, w.Body.String())
	}
	body := assertPointer(t, w)
	if body["ok"] != true {
		t.Fatalf("ok = %v", body["ok"])
	}
	var envelope struct {
		Report json.RawMessage `json:"report"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil || len(envelope.Report) == 0 {
		t.Fatalf("缺少 report 欄位: %v; body=%s", err, w.Body.String())
	}

	req := httptest.NewRequest(http.MethodPost, "/api/report", bytes.NewReader(envelope.Report))
	req.Header.Set("Authorization", "Bearer "+key.Token)
	resubmit := httptest.NewRecorder()
	handler.ServeHTTP(resubmit, req)
	if resubmit.Code != http.StatusOK {
		t.Fatalf("原樣重送 report 被拒: %d %s", resubmit.Code, resubmit.Body.String())
	}
}
