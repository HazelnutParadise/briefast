package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HazelnutParadise/briefast/internal/store"
)

const reprintBody = `台積電今日公布七月營收，單月合併營收達新台幣三千二百億元，較上月成長百分之十二，較去年同期成長百分之四十五，創下歷史新高。公司表示，先進製程需求持續強勁，三奈米與五奈米產能利用率維持高檔，人工智慧加速器訂單帶動高效能運算平台營收占比進一步提升。`

func seenMux(t *testing.T) (*store.Store, store.APIKey, http.Handler) {
	t.Helper()
	s, key, _, _ := setupHandler(t)
	h := NewSeenHandler(s)
	mux := http.NewServeMux()
	mux.Handle("/api/seen/peek", h.Peek())
	mux.Handle("/api/seen/record", h.Record())
	mux.Handle("/api/seen/similar", h.Similar())
	mux.Handle("/api/seen/events", h.AddEvent())
	mux.Handle("/api/seen/events/check", h.CheckEvent())
	return s, key, mux
}

func seenCall(t *testing.T, mux http.Handler, token, path string, body any) map[string]any {
	t.Helper()
	w := seenRaw(t, mux, token, path, body)
	if w.Code != http.StatusOK {
		t.Fatalf("%s 狀態碼 = %d, body=%s", path, w.Code, w.Body.String())
	}
	return assertPointer(t, w)
}

func seenRaw(t *testing.T, mux http.Handler, token, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

func TestSeenPeekRecordFlow(t *testing.T) {
	s, key, mux := seenMux(t)
	candidate := map[string]any{"url": "https://news.cnyes.com/news/id/6001", "title": "台積電七月營收創新高", "extra": "保留"}

	got := seenCall(t, mux, key.Token, "/api/seen/peek", map[string]any{"items": []any{candidate}})
	if len(got["new"].([]any)) != 1 || len(got["seen"].([]any)) != 0 {
		t.Fatalf("首次 peek 應全部未讀: %v", got)
	}
	if got["new"].([]any)[0].(map[string]any)["extra"] != "保留" {
		t.Fatal("未讀項目應原樣回傳")
	}

	rec := seenCall(t, mux, key.Token, "/api/seen/record", map[string]any{"items": []any{
		map[string]any{"url": candidate["url"], "title": candidate["title"], "content": reprintBody, "decision": "reported", "symbols": []string{"2330"}},
		map[string]any{"url": candidate["url"].(string) + "/?utm_source=line"},
	}})
	if rec["added"] != float64(1) || rec["duplicates"] != float64(1) {
		t.Fatalf("record 計數錯誤: %v", rec)
	}

	got = seenCall(t, mux, key.Token, "/api/seen/peek", map[string]any{"items": []any{
		map[string]any{"url": "https://news.cnyes.com/news/id/6001/?utm_source=line"},
		map[string]any{"url": "https://example.com/other", "title": "台積電七月營收創新高"},
		map[string]any{"url": "https://example.com/fresh", "title": "全新標題"},
	}})
	seen := got["seen"].([]any)
	if len(seen) != 2 || len(got["new"].([]any)) != 1 {
		t.Fatalf("第二次 peek 分類錯誤: %v", got)
	}
	if seen[0].(map[string]any)["seen_by"] != "source_id" || seen[1].(map[string]any)["seen_by"] != "title" {
		t.Fatalf("seen_by 錯誤: %v", seen)
	}

	sim := seenCall(t, mux, key.Token, "/api/seen/similar", map[string]any{"content": "記者王小明／台北報導 " + reprintBody})
	matches := sim["matches"].([]any)
	if sim["threshold"] != 0.25 || len(matches) != 1 || matches[0].(map[string]any)["similarity"].(float64) < 0.9 {
		t.Fatalf("similar 結果錯誤: %v", sim)
	}
	none := seenCall(t, mux, key.Token, "/api/seen/similar", map[string]any{"content": "颱風明天清晨在花蓮登陸，北部山區可能有豪雨，各縣市陸續宣布停班停課。"})
	if len(none["matches"].([]any)) != 0 {
		t.Fatalf("不相關內文不應命中: %v", none)
	}
	if n := countLogsWithAction(t, s, "seen_rejected_validation") + countLogsWithAction(t, s, "seen_rejected_auth"); n != 0 {
		t.Fatalf("成功呼叫不應寫 log, got %d", n)
	}
}

func TestSeenEventsRegisterAndCheck(t *testing.T) {
	_, key, mux := seenMux(t)
	eventKey := "2330|月營收|2026-07"
	before := seenCall(t, mux, key.Token, "/api/seen/events/check", map[string]any{"key": eventKey})
	if before["reported"] != false || before["record"] != nil {
		t.Fatalf("未登記事件: %v", before)
	}
	if add := seenCall(t, mux, key.Token, "/api/seen/events", map[string]any{"key": eventKey, "note": "營收創高"}); add["added"] != true {
		t.Fatalf("第一次登記: %v", add)
	}
	if again := seenCall(t, mux, key.Token, "/api/seen/events", map[string]any{"key": eventKey, "note": "覆寫"}); again["added"] != false {
		t.Fatalf("重複登記: %v", again)
	}
	after := seenCall(t, mux, key.Token, "/api/seen/events/check", map[string]any{"key": eventKey})
	record, _ := after["record"].(map[string]any)
	if after["reported"] != true || record["note"] != "營收創高" || record["event_key"] != eventKey {
		t.Fatalf("已登記事件: %v", after)
	}
}

func TestSeenRejectsBadKeyAndBodyWithLogs(t *testing.T) {
	s, key, mux := seenMux(t)
	if err := s.RevokeAPIKey(context.Background(), key.ID); err != nil {
		t.Fatal(err)
	}
	if w := seenRaw(t, mux, key.Token, "/api/seen/peek", map[string]any{"items": []any{}}); w.Code != http.StatusUnauthorized {
		t.Fatalf("撤銷 key 狀態碼 = %d", w.Code)
	}
	if n := countLogsWithAction(t, s, "seen_rejected_auth"); n != 1 {
		t.Fatalf("seen_rejected_auth = %d", n)
	}

	fresh, err := s.CreateAPIKey(context.Background(), "second", "second-token")
	if err != nil {
		t.Fatal(err)
	}
	bad := []struct {
		path string
		body any
	}{
		{"/api/seen/peek", map[string]any{}},
		{"/api/seen/record", map[string]any{"items": []any{map[string]any{"decision": "skipped"}}}},
		{"/api/seen/record", map[string]any{"items": []any{map[string]any{"url": 3}}}},
		{"/api/seen/similar", map[string]any{"content": " "}},
		{"/api/seen/similar", map[string]any{"content": "x", "threshold": 2}},
		{"/api/seen/events", map[string]any{"key": "2330|月營收"}},
		{"/api/seen/events/check", map[string]any{"key": "2330||2026-07"}},
		{"/api/seen/peek", map[string]any{"items": []any{}, "unknown": true}},
	}
	for _, c := range bad {
		w := seenRaw(t, mux, fresh.Token, c.path, c.body)
		if w.Code != http.StatusBadRequest {
			t.Errorf("%s %v 狀態碼 = %d, body=%s", c.path, c.body, w.Code, w.Body.String())
			continue
		}
		assertPointer(t, w)
	}
	if n := countLogsWithAction(t, s, "seen_rejected_validation"); n != len(bad) {
		t.Fatalf("seen_rejected_validation = %d, want %d", n, len(bad))
	}
}

func TestSeenPeekForgetsExpiredArticles(t *testing.T) {
	s, key, mux := seenMux(t)
	seenCall(t, mux, key.Token, "/api/seen/record", map[string]any{"items": []any{map[string]any{"url": "https://example.com/old/1"}}})
	if _, err := s.DB().Exec("UPDATE seen_articles SET created_at = '2000-01-01T00:00:00.000000000Z'"); err != nil {
		t.Fatal(err)
	}
	got := seenCall(t, mux, key.Token, "/api/seen/peek", map[string]any{"items": []any{map[string]any{"url": "https://example.com/old/1"}}})
	if len(got["new"].([]any)) != 1 {
		t.Fatalf("過期文章應視為未讀: %v", got)
	}
	var n int
	if err := s.DB().QueryRow("SELECT COUNT(*) FROM seen_articles").Scan(&n); err != nil || n != 0 {
		t.Fatalf("過期列應刪除, count = %d err = %v", n, err)
	}
}
