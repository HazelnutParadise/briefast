package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/HazelnutParadise/briefast/internal/store"
)

// ReadHandler 以與 ingest 相同的 Bearer key 提供單日報告的唯讀存取。
type ReadHandler struct {
	store *store.Store
}

func NewReadHandler(s *store.Store) http.Handler {
	return &ReadHandler{store: s}
}

func (h *ReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "errors": []string{"只接受 GET"}})
		return
	}

	date := r.PathValue("date")
	key, authenticated, err := authenticate(h.store, r)
	if err != nil {
		writeServerError(w, err)
		return
	}
	if !authenticated {
		entry := store.UpdateLog{KeyName: "unknown", Action: "read_rejected_auth", ReportDate: date, Detail: "missing or invalid bearer token"}
		if key.ID != 0 {
			entry.APIKeyID = &key.ID
			entry.KeyName = key.Name
		}
		if err := h.store.AddUpdateLog(r.Context(), entry); err != nil {
			writeServerError(w, err)
			return
		}
		writeJSON(w, http.StatusUnauthorized, map[string]any{"ok": false, "errors": []string{"無效或已撤銷的 API key"}})
		return
	}

	if !validReportDate(date) {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "errors": []string{"date 必須是有效的 YYYY-MM-DD 日期"}})
		return
	}

	value, err := h.store.ReportByDate(r.Context(), date)
	if errors.Is(err, store.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]any{"ok": false, "errors": []string{"找不到該日期的報告"}})
		return
	}
	if err != nil {
		writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func validReportDate(value string) bool {
	if len(value) != len("2006-01-02") {
		return false
	}
	parsed, err := time.Parse("2006-01-02", value)
	return err == nil && parsed.Format("2006-01-02") == value
}
