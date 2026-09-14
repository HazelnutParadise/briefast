## Problem

首頁報告日期列顯示原始字串「2026-09-14T07:50:27+0800 更新」，應顯示「07:50 更新」。

## Root Cause

上傳的每日報告 generated_at 時區寫成 +0800（缺冒號），不是 RFC 3339。首頁的 displayGeneratedAt 用 time.RFC3339 解析失敗後直接回傳原字串。每日報告的 Report.Validate 沒檢查 generated_at（法說會 brief 有檢查），所以錯誤格式被收下並存入資料庫。daily-brief 流程還會拿上一份報告的 generated_at 當新聞窗口起點，格式錯誤會連帶影響下一輪。

## Proposed Solution

- Report.Validate 加入 generated_at 必須是含時區 RFC 3339 的檢查，錯誤訊息與法說會一致，POST /api/report 因此回 400 並不寫入。
- 首頁 displayGeneratedAt 與法說會頁 displayGeneratedFull 改用共用解析函式，先試 RFC 3339，失敗再試 +0800 無冒號時區格式，讓已存入的舊資料正確顯示。

## Non-Goals

- 不做資料遷移改寫已存的 generated_at，舊資料靠顯示端容錯處理，正確值由重新上傳覆寫。
- API 不接受並自動修正 +0800 格式，維持嚴格契約讓 agent 端修正。

## Success Criteria

- 以 generated_at 為 2026-09-14T07:50:27+0800 的報告呼叫 Report.Validate，錯誤包含 generated_at。
- 首頁渲染 generated_at 為 +0800 格式的報告時顯示「07:50 更新」，不含原始字串。
- go test ./... 全數通過。

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `report-ingest-api`: schema 驗證加入 generated_at 格式檢查

## Impact

- Affected specs: report-ingest-api
- Affected code:
  - Modified: internal/report/schema.go, internal/report/schema_test.go, internal/site/site.go, internal/site/conference.go, internal/site/site_test.go
