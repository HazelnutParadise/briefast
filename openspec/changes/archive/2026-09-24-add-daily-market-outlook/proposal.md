## Why

現有盤前報告列出新聞、籌碼、個股判斷與法說會，但沒有整體台股大盤的當日方向判斷。讀者須自行拼湊各項訊號，首頁也無法直接回答「今天大盤可能怎麼走」。

## What Changes

- 每日報告新增整體大盤走勢預測，標明方向並用報告內的資料說明理由。
- 報告 API 接受、驗證並保存預測內容，舊報告仍可讀取。
- 首頁與歷史報告在盤前總覽顯示預測；沒有預測的舊報告不顯示空白欄位。
- 每日流程在完成新聞、籌碼與法說資訊整理後產出預測。

## Capabilities

### New Capabilities

無。

### Modified Capabilities

- `report-ingest-api`: 報告可攜帶當日大盤走勢預測，且有欄位驗證與舊資料相容性。
- `report-viewing`: 首頁與歷史報告顯示有內容的預測。
- `daily-brief-skill`: 每日流程依整份報告形成可追溯的大盤判斷。

## Impact

影響 `internal/report`、`internal/site`、`skills/daily-brief/SKILL.md` 與對應測試。SQLite 已保存完整 JSON，不需資料遷移；API 新欄位為選填，以保留舊報告與既有客戶端相容性。
