## Why

目前「今日大盤走勢預測」只有收盤方向和一段理由，讀者仍看不出開盤、盤中到尾盤可能如何演變。使用者希望看到有依據的走勢路徑，而非只看到漲跌判斷。

## What Changes

- 新每日報告在大盤預測中加入「預期走法」，描述開盤、盤中到尾盤的基準情境與會改變走法的條件。
- API 保存與讀回這段文字，若有提供卻是空白則拒絕；既有報告繼續可讀。
- 首頁與歷史報告在收盤方向後呈現走勢路徑，再呈現判斷依據。
- 每日流程在資料不足時明確說無法判斷盤中路徑，不編造轉折時間、點位或即時訊號。

## Capabilities

### New Capabilities

無。

### Modified Capabilities

- `report-ingest-api`: 大盤預測可攜帶有效的走勢路徑文字並原樣讀回。
- `report-viewing`: 有走勢路徑時在大盤預測區塊優先顯示。
- `daily-brief-skill`: 新報告必須產生有依據的盤中路徑情境或明示無法判斷。

## Impact

影響 `internal/report/schema.go`、`internal/site/site.go`、`internal/site/styles.go`、`skills/daily-brief/SKILL.md`、文件與對應測試。資料庫仍保存完整報告 JSON，不需 migration。新欄位選填於 API 以維持前版資料與客戶端相容；正式每日流程要求填寫。
