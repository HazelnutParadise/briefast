package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const openAPIPath = "/api/openapi.json"

// OverviewHandler serves GET /api, the public list of every endpoint.
func OverviewHandler(routes []Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		endpoints := make([]map[string]any, 0, len(routes))
		for _, route := range routes {
			endpoints = append(endpoints, map[string]any{
				"method":    route.Method,
				"path":      route.Path,
				"summary":   route.Summary,
				"auth":      authLabel(route.Auth),
				"operation": operationLink(route),
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "openapi": openAPIPath, "endpoints": endpoints})
	})
}

// OpenAPIHandler serves GET /api/openapi.json. The guide pointer lives in
// info.description and the Link header so the document stays valid OpenAPI.
func OpenAPIHandler(routes []Route) http.Handler {
	document := openAPIDocument(routes)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodGet) {
			return
		}
		w.Header().Set("Link", GuideLink)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(document)
	})
}

func authLabel(required bool) string {
	if required {
		return "bearer"
	}
	return "none"
}

func operationLink(route Route) string {
	pointer := strings.ReplaceAll(strings.ReplaceAll(route.Path, "~", "~0"), "/", "~1")
	return openAPIPath + "#/paths/" + pointer + "/" + strings.ToLower(route.Method)
}

var errorDescriptions = map[int]string{
	400: "請求內容不合格式，errors 列出原因",
	401: "缺少、無效或已撤銷的 API key",
	404: "找不到資料",
	405: "HTTP method 不符",
}

func openAPIDocument(routes []Route) map[string]any {
	paths := map[string]any{}
	for _, route := range routes {
		responses := map[string]any{"200": withDescription(route.Success, "成功")}
		for _, code := range append(append([]int{}, route.Errors...), 500) {
			description := errorDescriptions[code]
			if code == 500 {
				description = "伺服器無法完成請求"
			}
			responses[strconv.Itoa(code)] = withDescription(jsonContent(ref("Error")), description)
		}
		operation := map[string]any{
			"summary":     route.Summary,
			"operationId": strings.ToLower(route.Method) + strings.ReplaceAll(strings.NewReplacer("/", "_", "{", "", "}", "", "-", "_", ".", "_").Replace(route.Path), "__", "_"),
			"responses":   responses,
		}
		if route.Auth {
			operation["security"] = []any{map[string]any{"bearer": []any{}}}
		} else {
			operation["security"] = []any{}
		}
		if len(route.Parameters) != 0 {
			operation["parameters"] = route.Parameters
		}
		if route.RequestBody != nil {
			operation["requestBody"] = route.RequestBody
		}
		item, _ := paths[route.Path].(map[string]any)
		if item == nil {
			item = map[string]any{}
			paths[route.Path] = item
		}
		item[strings.ToLower(route.Method)] = operation
	}
	return map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":       "Briefast API",
			"version":     "1",
			"description": "每日產業與股市晨報網站的 API。" + GuidePointer + "除本文件外，每個 JSON 回應都附 guide 欄位與 Link 標頭提醒這件事。",
		},
		"paths": paths,
		"components": map[string]any{
			"securitySchemes": map[string]any{"bearer": map[string]any{"type": "http", "scheme": "bearer", "description": "後台建立的 API key"}},
			"schemas":         schemas(),
		},
	}
}

func schemas() map[string]any {
	source := object(map[string]any{"title": str(), "url": str()}, "title", "url")
	return map[string]any{
		"Error": object(map[string]any{"ok": constant(false), "errors": array(str()), "guide": str()}, "ok", "errors", "guide"),
		"Endpoint": object(map[string]any{
			"method": str(), "path": str(), "summary": str(), "auth": enum("bearer", "none"), "operation": str(),
		}, "method", "path", "summary", "auth", "operation"),
		"Report": describe(object(map[string]any{
			"date": dateSchema(), "headline": str(), "overview_md": str(), "watch_md": str(),
			"calls": object(map[string]any{
				"short_bull": array(ref("CallEntry")), "short_bear": array(ref("CallEntry")),
				"long_bull": array(ref("CallEntry")), "long_bear": array(ref("CallEntry")),
			}, "short_bull", "short_bear", "long_bull", "long_bear"),
			"industries": array(object(map[string]any{
				"name": str(), "events": array(object(map[string]any{"headline": str(), "summary_md": str()}, "headline", "summary_md")), "watch_md": str(),
			}, "name", "events", "watch_md")),
			"stock_news": array(object(map[string]any{
				"symbol": str(), "name": str(), "call": str(), "headline": str(), "summary_md": str(), "watch_md": str(),
				"sources": array(source), "chips": map[string]any{"type": "object"},
			}, "symbol", "name", "call", "headline", "summary_md", "watch_md", "sources")),
			"generated_at":   map[string]any{"type": "string", "format": "date-time"},
			"market_outlook": map[string]any{"type": "object"},
		}, "date", "headline", "overview_md", "watch_md", "calls", "industries", "stock_news", "generated_at"),
			"完整欄位規則與範例見 GET /api/guide 第 4 步；伺服器拒收未知欄位。"),
		"CallEntry": object(map[string]any{"symbol": str(), "name": str(), "reason": str()}, "symbol", "name", "reason"),
		"Conference": describe(object(map[string]any{
			"symbol": str(), "name": str(), "market": str(), "held_on": dateSchema(), "held_at": str(), "venue": str(),
			"announced_on": dateSchema(), "prediction": str(), "headline": str(), "summary_md": str(), "watch_md": str(),
			"fundamentals": map[string]any{"type": "object"}, "chips": map[string]any{"type": "object"},
			"sources": array(source), "generated_at": map[string]any{"type": "string", "format": "date-time"},
		}, "symbol", "name", "market", "held_on", "announced_on", "prediction", "headline", "summary_md", "watch_md", "sources", "generated_at"),
			"完整欄位規則與範例見 GET /api/guide 第 3 步；伺服器拒收未知欄位。"),
		"ConferenceListItem": object(map[string]any{
			"symbol": str(), "name": str(), "market": str(), "held_on": dateSchema(), "held_at": str(),
			"prediction": str(), "venue": str(), "announced_on": dateSchema(),
		}, "symbol", "name", "market", "held_on", "held_at", "prediction"),
		"SeenItems": object(map[string]any{"items": array(ref("SeenItem"))}, "items"),
		"SeenItem": describe(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"url": str(), "title": str(), "source_id": str(), "content": str(), "summary": str(),
				"symbols": array(str()), "event_key": str(), "published_at": str(),
				"decision": map[string]any{"type": "string", "default": "skipped"}, "call": str(),
			},
		}, "至少要有 url、title 或 source_id 其中一個；其他欄位會原樣保留在 peek 的回應裡。"),
	}
}

func withDescription(response map[string]any, description string) map[string]any {
	out := map[string]any{"description": description}
	for k, v := range response {
		out[k] = v
	}
	return out
}

func describe(schema map[string]any, description string) map[string]any {
	schema["description"] = description
	return schema
}

func jsonBody(schema map[string]any) map[string]any {
	return map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": schema}}}
}

func jsonContent(schema map[string]any) map[string]any {
	return map[string]any{"content": map[string]any{"application/json": map[string]any{"schema": withGuide(schema)}}}
}

// withGuide adds the guide pointer every JSON response carries.
func withGuide(schema map[string]any) map[string]any {
	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		return schema
	}
	properties["guide"] = str()
	return schema
}

func ref(name string) map[string]any { return map[string]any{"$ref": "#/components/schemas/" + name} }

func object(properties map[string]any, required ...string) map[string]any {
	return map[string]any{"type": "object", "properties": properties, "required": required}
}

func array(items map[string]any) map[string]any {
	return map[string]any{"type": "array", "items": items}
}

func str() map[string]any        { return map[string]any{"type": "string"} }
func integer() map[string]any    { return map[string]any{"type": "integer"} }
func boolean() map[string]any    { return map[string]any{"type": "boolean"} }
func positive() map[string]any   { return map[string]any{"type": "number", "exclusiveMinimum": 0} }
func dateSchema() map[string]any { return map[string]any{"type": "string", "format": "date"} }

func constant(value any) map[string]any { return map[string]any{"const": value} }

func enum(values ...string) map[string]any { return map[string]any{"type": "string", "enum": values} }

func eventKeySchema() map[string]any {
	return map[string]any{"type": "string", "pattern": `^[^|]+\|[^|]+\|[^|]+$`, "examples": []any{"2330|月營收|2026-07"}}
}
