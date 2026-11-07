package api

import (
	"io"
	"net/http"

	"github.com/HazelnutParadise/briefast/internal/guide"
	"github.com/HazelnutParadise/briefast/internal/store"
)

// GuideHandler serves GET /api/guide, the daily-brief workflow guide.
func GuideHandler(s *store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		if _, ok := authorizeRequest(s, w, r, "guide_rejected_auth", ""); !ok {
			return
		}
		w.Header().Set("Link", GuideLink)
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, guide.DailyBrief)
	})
}
