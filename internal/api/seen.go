package api

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/HazelnutParadise/briefast/internal/dedup"
	"github.com/HazelnutParadise/briefast/internal/store"
)

const defaultSimilarThreshold = 0.25

// SeenHandler serves the /api/seen endpoints that keep the daily brief's read
// log on the server.
type SeenHandler struct {
	store *store.Store
}

func NewSeenHandler(s *store.Store) *SeenHandler {
	return &SeenHandler{store: s}
}

type seenItemsRequest struct {
	Items *[]map[string]any `json:"items"`
}

// begin runs the checks every dedup call shares: POST only, a valid key, a
// single JSON object body, and pruning expired records. ok is false when a
// response has already been written.
func (h *SeenHandler) begin(w http.ResponseWriter, r *http.Request, target any) (store.APIKey, bool) {
	if !requireMethod(w, r, http.MethodPost) {
		return store.APIKey{}, false
	}
	key, ok := authorizeRequest(h.store, w, r, "seen_rejected_auth", "")
	if !ok {
		return key, false
	}
	if errs := decodeSingleObject(io.LimitReader(r.Body, maxReportBytes+1), target); len(errs) != 0 {
		rejectRequest(h.store, w, r, key, "seen_rejected_validation", "", errs)
		return key, false
	}
	if err := h.store.PruneSeen(r.Context()); err != nil {
		writeServerError(w, err)
		return key, false
	}
	return key, true
}

func (h *SeenHandler) reject(w http.ResponseWriter, r *http.Request, key store.APIKey, errs []string) {
	rejectRequest(h.store, w, r, key, "seen_rejected_validation", "", errs)
}

// Peek handles POST /api/seen/peek.
func (h *SeenHandler) Peek() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body seenItemsRequest
		key, ok := h.begin(w, r, &body)
		if !ok {
			return
		}
		items, errs := validateItems(body.Items)
		if len(errs) != 0 {
			h.reject(w, r, key, errs)
			return
		}
		fresh, known := make([]map[string]any, 0), make([]map[string]any, 0)
		for _, item := range items {
			f := fingerprint(item)
			reason, err := h.store.SeenMatch(r.Context(), f.SourceID, f.URLHash, f.TitleHash)
			if err != nil {
				writeServerError(w, err)
				return
			}
			if reason == "" {
				fresh = append(fresh, item)
				continue
			}
			marked := make(map[string]any, len(item)+1)
			for k, v := range item {
				marked[k] = v
			}
			marked["seen_by"] = reason
			known = append(known, marked)
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "new": fresh, "seen": known})
	})
}

// Record handles POST /api/seen/record.
func (h *SeenHandler) Record() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body seenItemsRequest
		key, ok := h.begin(w, r, &body)
		if !ok {
			return
		}
		items, errs := validateItems(body.Items)
		if len(errs) != 0 {
			h.reject(w, r, key, errs)
			return
		}
		added := 0
		for _, item := range items {
			inserted, err := h.store.RecordSeenArticle(r.Context(), fingerprint(item))
			if err != nil {
				writeServerError(w, err)
				return
			}
			if inserted {
				added++
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "added": added, "duplicates": len(items) - added})
	})
}

type similarRequest struct {
	Content   string   `json:"content"`
	Threshold *float64 `json:"threshold"`
}

type similarMatch struct {
	Similarity float64 `json:"similarity"`
	Title      string  `json:"title"`
	URL        string  `json:"url"`
	EventKey   string  `json:"event_key"`
}

// Similar handles POST /api/seen/similar.
func (h *SeenHandler) Similar() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body similarRequest
		key, ok := h.begin(w, r, &body)
		if !ok {
			return
		}
		threshold := defaultSimilarThreshold
		if body.Threshold != nil {
			threshold = *body.Threshold
		}
		var errs []string
		if strings.TrimSpace(body.Content) == "" {
			errs = append(errs, "content 不得為空")
		}
		if threshold < 0 || threshold > 1 || math.IsNaN(threshold) {
			errs = append(errs, "threshold 必須介於 0 到 1")
		}
		if len(errs) != 0 {
			h.reject(w, r, key, errs)
			return
		}
		stored, err := h.store.SeenSketches(r.Context())
		if err != nil {
			writeServerError(w, err)
			return
		}
		target := dedup.Sketch(body.Content)
		matches := make([]similarMatch, 0)
		for _, article := range stored {
			score := dedup.Similarity(target, article.Sketch)
			if score >= threshold {
				matches = append(matches, similarMatch{Similarity: math.Round(score*1000) / 1000, Title: article.Title, URL: article.URL, EventKey: article.EventKey})
			}
		}
		sort.SliceStable(matches, func(i, j int) bool { return matches[i].Similarity > matches[j].Similarity })
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "threshold": threshold, "matches": matches})
	})
}

type eventRequest struct {
	Key  string `json:"key"`
	Note string `json:"note"`
}

func (e eventRequest) validate() []string {
	parts := strings.Split(e.Key, "|")
	if len(parts) != 3 {
		return []string{"key 必須是「代號|事件類型|期間」三段"}
	}
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			return []string{"key 的三段都不得為空"}
		}
	}
	return nil
}

// AddEvent handles POST /api/seen/events.
func (h *SeenHandler) AddEvent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body eventRequest
		key, ok := h.begin(w, r, &body)
		if !ok {
			return
		}
		if errs := body.validate(); len(errs) != 0 {
			h.reject(w, r, key, errs)
			return
		}
		added, err := h.store.AddSeenEvent(r.Context(), body.Key, body.Note)
		if err != nil {
			writeServerError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "key": body.Key, "added": added})
	})
}

// CheckEvent handles POST /api/seen/events/check.
func (h *SeenHandler) CheckEvent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body eventRequest
		key, ok := h.begin(w, r, &body)
		if !ok {
			return
		}
		if errs := body.validate(); len(errs) != 0 {
			h.reject(w, r, key, errs)
			return
		}
		event, err := h.store.SeenEventByKey(r.Context(), body.Key)
		if errors.Is(err, store.ErrNotFound) {
			writeJSON(w, http.StatusOK, map[string]any{"ok": true, "key": body.Key, "reported": false, "record": nil})
			return
		}
		if err != nil {
			writeServerError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "key": body.Key, "reported": true, "record": event})
	})
}

var stringFields = []string{"url", "title", "source_id", "content", "summary", "event_key", "published_at", "decision", "call"}

func validateItems(items *[]map[string]any) ([]map[string]any, []string) {
	if items == nil {
		return nil, []string{"items 必須是陣列"}
	}
	var errs []string
	for i, item := range *items {
		if item == nil {
			errs = append(errs, fmt.Sprintf("items[%d] 必須是物件", i))
			continue
		}
		for _, field := range stringFields {
			if value, present := item[field]; present {
				if _, isString := value.(string); !isString {
					errs = append(errs, fmt.Sprintf("items[%d].%s 必須是字串", i, field))
				}
			}
		}
		if value, present := item["symbols"]; present {
			if _, ok := toStrings(value); !ok {
				errs = append(errs, fmt.Sprintf("items[%d].symbols 必須是字串陣列", i))
			}
		}
		if text(item, "url") == "" && text(item, "title") == "" && text(item, "source_id") == "" {
			errs = append(errs, fmt.Sprintf("items[%d] 至少要有 url、title 或 source_id 其中一個", i))
		}
	}
	return *items, errs
}

// fingerprint turns a submitted item into the stored article shape.
func fingerprint(item map[string]any) store.SeenArticle {
	rawURL := text(item, "url")
	a := store.SeenArticle{
		SourceID:    dedup.SourceID(text(item, "source_id"), rawURL),
		URL:         rawURL,
		Title:       text(item, "title"),
		EventKey:    text(item, "event_key"),
		PublishedAt: text(item, "published_at"),
		Decision:    text(item, "decision"),
		Call:        text(item, "call"),
	}
	a.Symbols, _ = toStrings(item["symbols"])
	if a.Decision == "" {
		a.Decision = "skipped"
	}
	if rawURL != "" {
		a.URLHash = dedup.Digest(dedup.CanonicalURL(rawURL))
		if u, err := url.Parse(rawURL); err == nil {
			a.SourceDomain = strings.ToLower(u.Host)
		}
	}
	if title := dedup.NormalizeText(a.Title); title != "" {
		a.TitleHash = dedup.Digest(title)
	}
	content := text(item, "content")
	if content == "" {
		content = text(item, "summary")
	}
	if normalized := dedup.NormalizeText(content); normalized != "" {
		a.ContentHash = dedup.Digest(normalized)
		a.LeadHash = dedup.Digest(string([]rune(normalized)[:min(100, len([]rune(normalized)))]))
		a.Sketch = dedup.Sketch(content)
	}
	return a
}

func text(item map[string]any, field string) string {
	value, _ := item[field].(string)
	return strings.TrimSpace(value)
}

func toStrings(value any) ([]string, bool) {
	if value == nil {
		return nil, true
	}
	list, ok := value.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(list))
	for _, entry := range list {
		s, ok := entry.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}
