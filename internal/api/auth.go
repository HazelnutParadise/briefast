package api

import (
	"net/http"
	"strings"

	"github.com/HazelnutParadise/briefast/internal/store"
)

// authorizeRequest checks the Bearer key. An unauthenticated call is logged
// with rejectAction and answered 401; the caller stops when ok is false.
func authorizeRequest(s *store.Store, w http.ResponseWriter, r *http.Request, rejectAction, date string) (store.APIKey, bool) {
	key, authenticated, err := authenticate(s, r)
	if err != nil {
		writeServerError(w, err)
		return key, false
	}
	if authenticated {
		return key, true
	}
	entry := store.UpdateLog{KeyName: "unknown", Action: rejectAction, ReportDate: date, Detail: "missing or invalid bearer token"}
	if key.ID != 0 {
		entry.APIKeyID = &key.ID
		entry.KeyName = key.Name
	}
	if err := s.AddUpdateLog(r.Context(), entry); err != nil {
		writeServerError(w, err)
		return key, false
	}
	writeJSON(w, http.StatusUnauthorized, map[string]any{"ok": false, "errors": []string{"無效或已撤銷的 API key"}})
	return key, false
}

// rejectRequest logs a validation failure and answers 400 with the errors.
func rejectRequest(s *store.Store, w http.ResponseWriter, r *http.Request, key store.APIKey, action, date string, validationErrors []string) {
	if err := s.AddUpdateLog(r.Context(), store.UpdateLog{
		APIKeyID: &key.ID, KeyName: key.Name, ReportDate: date,
		Action: action, Detail: strings.Join(validationErrors, "; "),
	}); err != nil {
		writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "errors": validationErrors})
}
