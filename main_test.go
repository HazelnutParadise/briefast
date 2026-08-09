package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/HazelnutParadise/briefast/internal/store"
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
}
