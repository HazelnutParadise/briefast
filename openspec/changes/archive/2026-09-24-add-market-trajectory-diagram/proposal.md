## Why

目前大盤預測以方向標籤和文字敘述呈現，讀者仍要自行把開盤、盤中、尾盤的預期走法轉成視覺形狀。使用者希望直接用圖看出整天的基準情境。

## What Changes

- 新每日報告可提供三個階段相對前收的定性位置，形成盤前的「預期走勢示意」圖。
- 首頁與歷史報告在走法文字前顯示圖、階段標籤與非點位聲明；保留原有走法文字與條件。
- API 驗證並保存圖的定性欄位，舊報告照常顯示；證據不足或無法判斷時不產生圖。
- 每日流程要求圖與文字、收盤方向一致，避免只憑方向自動畫一條上升或下降線。

## Capabilities

### New Capabilities

無。

### Modified Capabilities

- `report-ingest-api`: 大盤預測接受選填且有效的三階段示意圖資料。
- `report-viewing`: 有示意圖資料時在首頁及歷史報告呈現有標示的定性走勢圖。
- `daily-brief-skill`: 有足夠證據時填三階段示意圖，沒有時保留文字並省略圖。

## Impact

影響 `internal/report/schema.go`、`internal/site/site.go`、`internal/site/styles.go`、`skills/daily-brief/SKILL.md`、`README.md`、`DESIGN.md` 與相關測試。完整報告仍以 JSON 存在原 SQLite 欄位，不需資料遷移；不新增圖表依賴或即時行情來源。
