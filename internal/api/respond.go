package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// GuidePointer is attached to every JSON response under /api/ so an agent that
// starts from any endpoint learns where the workflow guide is.
const GuidePointer = "先用同一把 API key 呼叫 GET /api/guide，讀完每日晨報流程說明再呼叫其他端點。"

// GuideLink is the matching Link header value.
const GuideLink = `</api/guide>; rel="help"`

func writeServerError(w http.ResponseWriter, err error) {
	log.Printf("briefast api: %v", err)
	writeJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "errors": []string{"伺服器無法完成請求"}})
}

func writeJSON(w http.ResponseWriter, status int, body map[string]any) {
	body["guide"] = GuidePointer
	w.Header().Set("Link", GuideLink)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// NotFoundHandler answers /api/ paths that no endpoint serves, so a mistyped
// call gets JSON with the guide pointer instead of the homepage HTML.
func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"ok":     false,
			"errors": []string{"沒有這個 API 端點：" + r.Method + " " + r.URL.Path + "，端點清單見 GET /api"},
		})
	})
}
