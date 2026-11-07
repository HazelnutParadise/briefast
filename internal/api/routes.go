package api

import (
	"net/http"

	"github.com/HazelnutParadise/briefast/internal/store"
)

// Route is one public /api endpoint. The same table mounts the handlers and
// generates GET /api and GET /api/openapi.json, so the documentation cannot
// drift from what the server serves.
type Route struct {
	Method      string
	Path        string
	Summary     string
	Auth        bool
	Handler     http.Handler
	Parameters  []map[string]any
	RequestBody map[string]any
	Success     map[string]any
	Errors      []int
}

// Routes returns every /api endpoint, including the overview and OpenAPI
// document that describe the list itself.
func Routes(s *store.Store, notifier Notifier) []Route {
	conferences := NewConferenceHandler(s, notifier)
	seen := NewSeenHandler(s)
	routes := []Route{
		{
			Method: http.MethodGet, Path: "/api/guide", Auth: true, Handler: GuideHandler(s),
			Summary: "讀取每日晨報的完整流程說明（Markdown），呼叫其他端點前先讀",
			Success: map[string]any{"content": map[string]any{"text/markdown": map[string]any{"schema": map[string]any{"type": "string"}}}},
			Errors:  []int{401, 405},
		},
		{
			Method: http.MethodPost, Path: "/api/report", Auth: true, Handler: NewReportHandler(s, notifier),
			Summary:     "發布或全量覆寫當日晨報",
			RequestBody: jsonBody(ref("Report")),
			Success:     jsonContent(object(map[string]any{"ok": constant(true), "date": dateSchema()}, "ok", "date")),
			Errors:      []int{400, 401, 405},
		},
		{
			Method: http.MethodGet, Path: "/api/report/{date}", Auth: true, Handler: NewReadHandler(s),
			Summary:    "讀取指定日期已發布的晨報",
			Parameters: []map[string]any{{"name": "date", "in": "path", "required": true, "schema": dateSchema()}},
			Success:    jsonContent(object(map[string]any{"ok": constant(true), "report": ref("Report")}, "ok", "report")),
			Errors:     []int{400, 401, 404, 405},
		},
		{
			Method: http.MethodPost, Path: "/api/conference", Auth: true, Handler: conferences.Ingest(),
			Summary:     "登記或覆寫法說會會前報告，以 (symbol, held_on) 為鍵",
			RequestBody: jsonBody(ref("Conference")),
			Success:     jsonContent(object(map[string]any{"ok": constant(true), "symbol": str(), "held_on": dateSchema()}, "ok", "symbol", "held_on")),
			Errors:      []int{400, 401, 405},
		},
		{
			Method: http.MethodPost, Path: "/api/conference/settle", Auth: true, Handler: conferences.Settle(),
			Summary: "回填法說會會後結算，由伺服器依會前預測判定命中",
			RequestBody: jsonBody(object(map[string]any{
				"symbol": str(), "held_on": dateSchema(), "pre_close_date": dateSchema(), "pre_close": positive(),
				"post_close_date": dateSchema(), "post_close": positive(),
			}, "symbol", "held_on", "pre_close_date", "pre_close", "post_close_date", "post_close")),
			Success: jsonContent(object(map[string]any{"ok": constant(true), "symbol": str(), "held_on": dateSchema(), "outcome": str()}, "ok", "symbol", "held_on", "outcome")),
			Errors:  []int{400, 401, 404, 405},
		},
		{
			Method: http.MethodGet, Path: "/api/conferences", Auth: true, Handler: conferences.List(),
			Summary: "列出未來場次與待結算場次",
			Success: jsonContent(object(map[string]any{"upcoming": array(ref("ConferenceListItem")), "pending_settlement": array(ref("ConferenceListItem"))}, "upcoming", "pending_settlement")),
			Errors:  []int{401, 405},
		},
		{
			Method: http.MethodPost, Path: "/api/seen/peek", Auth: true, Handler: seen.Peek(),
			Summary:     "從候選新聞中找出還沒處理過的，依原生 ID、網址、標題比對",
			RequestBody: jsonBody(ref("SeenItems")),
			Success: jsonContent(object(map[string]any{
				"ok": constant(true), "new": array(ref("SeenItem")),
				"seen": array(map[string]any{"allOf": []any{ref("SeenItem"), object(map[string]any{"seen_by": enum("source_id", "url", "title")}, "seen_by")}}),
			}, "ok", "new", "seen")),
			Errors: []int{400, 401, 405},
		},
		{
			Method: http.MethodPost, Path: "/api/seen/record", Auth: true, Handler: seen.Record(),
			Summary:     "記下處理過的新聞（含判定不報的），原生 ID 或網址已存在的略過",
			RequestBody: jsonBody(ref("SeenItems")),
			Success:     jsonContent(object(map[string]any{"ok": constant(true), "added": integer(), "duplicates": integer()}, "ok", "added", "duplicates")),
			Errors:      []int{400, 401, 405},
		},
		{
			Method: http.MethodPost, Path: "/api/seen/similar", Auth: true, Handler: seen.Similar(),
			Summary:     "找出與內文相近的已處理文章，用來抓改寫轉載",
			RequestBody: jsonBody(object(map[string]any{"content": str(), "threshold": map[string]any{"type": "number", "minimum": 0, "maximum": 1, "default": defaultSimilarThreshold}}, "content")),
			Success: jsonContent(object(map[string]any{
				"ok": constant(true), "threshold": map[string]any{"type": "number"},
				"matches": array(object(map[string]any{"similarity": map[string]any{"type": "number"}, "title": str(), "url": str(), "event_key": str()}, "similarity", "title", "url", "event_key")),
			}, "ok", "threshold", "matches")),
			Errors: []int{400, 401, 405},
		},
		{
			Method: http.MethodPost, Path: "/api/seen/events", Auth: true, Handler: seen.AddEvent(),
			Summary:     "登記跨天事件，鍵為「代號|事件類型|期間」，已存在的不變",
			RequestBody: jsonBody(object(map[string]any{"key": eventKeySchema(), "note": str()}, "key")),
			Success:     jsonContent(object(map[string]any{"ok": constant(true), "key": str(), "added": boolean()}, "ok", "key", "added")),
			Errors:      []int{400, 401, 405},
		},
		{
			Method: http.MethodPost, Path: "/api/seen/events/check", Auth: true, Handler: seen.CheckEvent(),
			Summary:     "查詢跨天事件是否已登記",
			RequestBody: jsonBody(object(map[string]any{"key": eventKeySchema()}, "key")),
			Success: jsonContent(object(map[string]any{
				"ok": constant(true), "key": str(), "reported": boolean(),
				"record": map[string]any{"oneOf": []any{map[string]any{"type": "null"}, object(map[string]any{"event_key": str(), "note": str(), "first_seen_at": map[string]any{"type": "string", "format": "date-time"}}, "event_key", "note", "first_seen_at")}},
			}, "ok", "key", "reported", "record")),
			Errors: []int{400, 401, 405},
		},
		{
			Method: http.MethodPost, Path: "/api/tx-night", Auth: true, Handler: TXNightHandler(s),
			Summary: "上傳兩份期交所 TX 每日行情 CSV（CP950），取得報告日夜盤對前一交易日同契約日盤的漲跌",
			RequestBody: map[string]any{"required": true, "content": map[string]any{"multipart/form-data": map[string]any{"schema": object(map[string]any{
				"report_date": dateSchema(), "previous_date": dateSchema(),
				"report_csv": map[string]any{"type": "string", "format": "binary"}, "previous_csv": map[string]any{"type": "string", "format": "binary"},
			}, "report_date", "previous_date", "report_csv", "previous_csv")}}},
			Success: jsonContent(object(map[string]any{
				"ok": constant(true), "report_date": dateSchema(), "previous_regular_date": dateSchema(), "contract": str(),
				"night_close": positive(), "previous_regular_close": positive(), "change_points": map[string]any{"type": "number"},
				"change_pct": map[string]any{"type": "number"}, "direction": enum("up", "down", "flat"), "night_volume": integer(),
			}, "ok", "report_date", "previous_regular_date", "contract", "night_close", "previous_regular_close", "change_points", "change_pct", "direction", "night_volume")),
			Errors: []int{400, 401, 405},
		},
	}
	routes = append(routes,
		Route{Method: http.MethodGet, Path: "/api", Summary: "列出所有 API 端點", Success: jsonContent(object(map[string]any{"ok": constant(true), "openapi": str(), "endpoints": array(ref("Endpoint"))}, "ok", "openapi", "endpoints")), Errors: []int{405}},
		Route{Method: http.MethodGet, Path: "/api/openapi.json", Summary: "OpenAPI 3.1 文件", Success: jsonContent(map[string]any{"type": "object"}), Errors: []int{405}},
	)
	n := len(routes)
	routes[n-2].Handler = OverviewHandler(routes)
	routes[n-1].Handler = OpenAPIHandler(routes)
	return routes
}

// Mount registers every route plus the JSON 404 for other /api/ paths.
func Mount(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		mux.Handle(route.Path, route.Handler)
	}
	mux.Handle("/api/", NotFoundHandler())
}
