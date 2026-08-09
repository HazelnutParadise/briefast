package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
)

const maxReportBytes = 2 << 20

type Notifier interface {
	Notify()
}

type ReportHandler struct {
	store    *store.Store
	notifier Notifier
}

func NewReportHandler(s *store.Store, notifier Notifier) http.Handler {
	return &ReportHandler{store: s, notifier: notifier}
}

func (h *ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "errors": []string{"只接受 POST"}})
		return
	}

	key, authenticated, err := h.authenticate(r)
	if err != nil {
		writeServerError(w, err)
		return
	}
	if !authenticated {
		entry := store.UpdateLog{KeyName: "unknown", Action: "ingest_rejected_auth", Detail: "missing or invalid bearer token"}
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

	payload, reportValue, validationErrors := decodeAndValidate(r.Body)
	if len(validationErrors) != 0 {
		if err := h.store.AddUpdateLog(r.Context(), store.UpdateLog{
			APIKeyID: &key.ID, KeyName: key.Name, ReportDate: reportValue.Date,
			Action: "ingest_rejected_schema", Detail: strings.Join(validationErrors, "; "),
		}); err != nil {
			writeServerError(w, err)
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "errors": validationErrors})
		return
	}

	if err := h.store.IngestReport(r.Context(), reportValue, payload, key); err != nil {
		writeServerError(w, err)
		return
	}
	if h.notifier != nil {
		h.notifier.Notify()
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "date": reportValue.Date})
}

func (h *ReportHandler) authenticate(r *http.Request) (store.APIKey, bool, error) {
	return authenticate(h.store, r)
}

func authenticate(s *store.Store, r *http.Request) (store.APIKey, bool, error) {
	header := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) || strings.TrimSpace(strings.TrimPrefix(header, prefix)) == "" {
		return store.APIKey{}, false, nil
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if strings.ContainsAny(token, " \t\r\n") {
		return store.APIKey{}, false, nil
	}
	key, err := s.APIKeyByToken(r.Context(), token)
	if errors.Is(err, store.ErrNotFound) {
		return store.APIKey{}, false, nil
	}
	if err != nil {
		return store.APIKey{}, false, err
	}
	return key, key.Active(), nil
}

func decodeAndValidate(body io.Reader) ([]byte, report.Report, []string) {
	limited := io.LimitReader(body, maxReportBytes+1)
	payload, err := io.ReadAll(limited)
	if err != nil {
		return nil, report.Report{}, []string{"讀取 request body 失敗"}
	}
	if len(payload) > maxReportBytes {
		return nil, report.Report{}, []string{"request body 超過 2 MiB"}
	}
	var value report.Report
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		return payload, value, []string{fmt.Sprintf("JSON 格式錯誤：%v", err)}
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return payload, value, []string{"request body 只能包含一個 JSON 物件"}
	}
	return payload, value, value.Validate()
}

func writeServerError(w http.ResponseWriter, err error) {
	log.Printf("report api: %v", err)
	writeJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "errors": []string{"伺服器無法完成請求"}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
