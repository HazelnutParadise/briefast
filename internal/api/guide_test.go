package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HazelnutParadise/briefast/internal/guide"
)

func guideRequest(t *testing.T, handler http.Handler, token string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/guide", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestGuideServesMarkdownWithoutLogging(t *testing.T) {
	s, key, _, _ := setupHandler(t)
	w := guideRequest(t, GuideHandler(s), key.Token)
	if w.Code != http.StatusOK {
		t.Fatalf("狀態碼 = %d, body=%s", w.Code, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/markdown; charset=utf-8" {
		t.Fatalf("Content-Type = %q", ct)
	}
	if w.Header().Get("Link") != GuideLink {
		t.Fatalf("Link 標頭 = %q", w.Header().Get("Link"))
	}
	if w.Body.String() != guide.DailyBrief {
		t.Fatal("回應內容不是嵌入的流程說明")
	}
	for _, step := range []string{"交易日閘門", "## 1. 蒐集產業與股市新聞", "## 2. 判讀產業與個股", "## 3. 法說會", "## 4. 組成並驗證報告 JSON", "## 5. POST 到網站 API"} {
		if !strings.Contains(w.Body.String(), step) {
			t.Fatalf("說明缺少 %q", step)
		}
	}
	if n := countLogsWithAction(t, s, "guide_rejected_auth"); n != 0 {
		t.Fatalf("成功讀取不應寫 log, got %d", n)
	}
}

func TestGuideRejectsRevokedKeyAndLogs(t *testing.T) {
	s, key, _, _ := setupHandler(t)
	if err := s.RevokeAPIKey(context.Background(), key.ID); err != nil {
		t.Fatal(err)
	}
	w := guideRequest(t, GuideHandler(s), key.Token)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("狀態碼 = %d", w.Code)
	}
	if body := assertPointer(t, w); body["ok"] != false {
		t.Fatalf("ok = %v", body["ok"])
	}
	if n := countLogsWithAction(t, s, "guide_rejected_auth"); n != 1 {
		t.Fatalf("guide_rejected_auth 筆數 = %d, want 1", n)
	}
}

func TestGuideRejectsWrongMethod(t *testing.T) {
	s, key, _, _ := setupHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/guide", nil)
	req.Header.Set("Authorization", "Bearer "+key.Token)
	w := httptest.NewRecorder()
	GuideHandler(s).ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("狀態碼 = %d", w.Code)
	}
	assertPointer(t, w)
}
