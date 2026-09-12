package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
)

// ConferenceHandler serves the three conference endpoints. Now is injectable so
// tests can pin "today in Taipei" without waiting for the calendar.
type ConferenceHandler struct {
	store    *store.Store
	notifier Notifier
	now      func() time.Time
}

func NewConferenceHandler(s *store.Store, notifier Notifier) *ConferenceHandler {
	return &ConferenceHandler{store: s, notifier: notifier, now: time.Now}
}

// WithNow returns a copy whose clock is fixed; tests use it to control today.
func (h *ConferenceHandler) WithNow(now func() time.Time) *ConferenceHandler {
	copy := *h
	copy.now = now
	return &copy
}

// Ingest handles POST /api/conference.
func (h *ConferenceHandler) Ingest() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodPost) {
			return
		}
		key, ok := h.authorize(w, r, "conference_rejected_auth", "")
		if !ok {
			return
		}
		payload, value, validationErrors := decodeConference(r.Body)
		if len(validationErrors) != 0 {
			h.reject(w, r, key, "conference_rejected_schema", value.HeldOn, validationErrors)
			return
		}
		if err := h.store.IngestConference(r.Context(), value, payload, key); err != nil {
			writeServerError(w, err)
			return
		}
		if h.notifier != nil {
			h.notifier.Notify()
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "symbol": value.Symbol, "held_on": value.HeldOn})
	})
}

type settlementRequest struct {
	Symbol        string  `json:"symbol"`
	HeldOn        string  `json:"held_on"`
	PreCloseDate  string  `json:"pre_close_date"`
	PreClose      float64 `json:"pre_close"`
	PostCloseDate string  `json:"post_close_date"`
	PostClose     float64 `json:"post_close"`
}

func (s settlementRequest) validate() []string {
	errs := make([]string, 0)
	if strings.TrimSpace(s.Symbol) == "" {
		errs = append(errs, "symbol 不得為空")
	}
	if !validReportDate(s.HeldOn) {
		errs = append(errs, "held_on 必須是有效的 YYYY-MM-DD 日期")
	}
	if !validReportDate(s.PreCloseDate) {
		errs = append(errs, "pre_close_date 必須是有效的 YYYY-MM-DD 日期")
	} else if validReportDate(s.HeldOn) && s.PreCloseDate >= s.HeldOn {
		errs = append(errs, "pre_close_date 必須早於 held_on")
	}
	if !validReportDate(s.PostCloseDate) {
		errs = append(errs, "post_close_date 必須是有效的 YYYY-MM-DD 日期")
	} else if validReportDate(s.HeldOn) && s.PostCloseDate <= s.HeldOn {
		errs = append(errs, "post_close_date 必須晚於 held_on")
	}
	if s.PreClose <= 0 {
		errs = append(errs, "pre_close 必須大於 0")
	}
	if s.PostClose <= 0 {
		errs = append(errs, "post_close 必須大於 0")
	}
	return errs
}

// Settle handles POST /api/conference/settle.
func (h *ConferenceHandler) Settle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodPost) {
			return
		}
		key, ok := h.authorize(w, r, "settle_rejected_auth", "")
		if !ok {
			return
		}
		var value settlementRequest
		decodeErrors := decodeSingleObject(io.LimitReader(r.Body, maxReportBytes+1), &value)
		validationErrors := decodeErrors
		if len(decodeErrors) == 0 {
			validationErrors = value.validate()
		}
		if len(validationErrors) != 0 {
			h.reject(w, r, key, "settle_rejected_schema", value.HeldOn, validationErrors)
			return
		}
		outcome, err := h.store.SettleConference(r.Context(), value.Symbol, value.HeldOn, store.Settlement{
			PreCloseDate: value.PreCloseDate, PreClose: value.PreClose,
			PostCloseDate: value.PostCloseDate, PostClose: value.PostClose,
		}, key)
		if errors.Is(err, store.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"ok": false, "errors": []string{"找不到該場法說會"}})
			return
		}
		if err != nil {
			writeServerError(w, err)
			return
		}
		if h.notifier != nil {
			h.notifier.Notify()
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "symbol": value.Symbol, "held_on": value.HeldOn, "outcome": outcome})
	})
}

type conferenceListItem struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Market      string `json:"market"`
	HeldOn      string `json:"held_on"`
	HeldAt      string `json:"held_at"`
	Prediction  string `json:"prediction"`
	Venue       string `json:"venue,omitempty"`
	AnnouncedOn string `json:"announced_on,omitempty"`
}

// List handles GET /api/conferences.
func (h *ConferenceHandler) List() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		if _, ok := h.authorize(w, r, "list_rejected_auth", ""); !ok {
			return
		}
		today := report.TaipeiDate(h.now())
		upcoming, err := h.store.ListUpcomingConferences(r.Context(), today)
		if err != nil {
			writeServerError(w, err)
			return
		}
		pending, err := h.store.ListPendingSettlement(r.Context(), today)
		if err != nil {
			writeServerError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"upcoming":           listItems(upcoming, true),
			"pending_settlement": listItems(pending, false),
		})
	})
}

func listItems(rows []store.ConferenceSummary, withAnnounced bool) []conferenceListItem {
	out := make([]conferenceListItem, 0, len(rows))
	for _, row := range rows {
		item := conferenceListItem{Symbol: row.Symbol, Name: row.Name, Market: row.Market, HeldOn: row.HeldOn, HeldAt: row.HeldAt, Prediction: row.Prediction}
		if withAnnounced {
			item.AnnouncedOn = row.AnnouncedOn
			item.Venue = row.Venue
		}
		out = append(out, item)
	}
	return out
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "errors": []string{"只接受 " + method}})
	return false
}

// authorize mirrors the report handlers' behaviour: an unauthenticated call
// is logged with the given action and answered 401.
func (h *ConferenceHandler) authorize(w http.ResponseWriter, r *http.Request, rejectAction, date string) (store.APIKey, bool) {
	key, authenticated, err := authenticate(h.store, r)
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
	if err := h.store.AddUpdateLog(r.Context(), entry); err != nil {
		writeServerError(w, err)
		return key, false
	}
	writeJSON(w, http.StatusUnauthorized, map[string]any{"ok": false, "errors": []string{"無效或已撤銷的 API key"}})
	return key, false
}

func (h *ConferenceHandler) reject(w http.ResponseWriter, r *http.Request, key store.APIKey, action, date string, validationErrors []string) {
	if err := h.store.AddUpdateLog(r.Context(), store.UpdateLog{
		APIKeyID: &key.ID, KeyName: key.Name, ReportDate: date,
		Action: action, Detail: strings.Join(validationErrors, "; "),
	}); err != nil {
		writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "errors": validationErrors})
}

func decodeConference(body io.Reader) ([]byte, report.Conference, []string) {
	limited := io.LimitReader(body, maxReportBytes+1)
	payload, err := io.ReadAll(limited)
	if err != nil {
		return nil, report.Conference{}, []string{"讀取 request body 失敗"}
	}
	if len(payload) > maxReportBytes {
		return nil, report.Conference{}, []string{"request body 超過 2 MiB"}
	}
	var value report.Conference
	if errs := decodeSingleObject(bytes.NewReader(payload), &value); len(errs) != 0 {
		return payload, value, errs
	}
	return payload, value, value.Validate()
}

// decodeSingleObject enforces the same body discipline as report ingestion:
// one JSON object, no unknown fields, nothing after it.
func decodeSingleObject(body io.Reader, target any) []string {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return []string{fmt.Sprintf("JSON 格式錯誤：%v", err)}
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return []string{"request body 只能包含一個 JSON 物件"}
	}
	return nil
}
