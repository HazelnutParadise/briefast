## Why

讀者掃報告時需要靠標題判斷「這段值不值得讀」。目前產業區塊每段只有短語標籤（例如「記憶體供需」），看不出發生什麼事；個股條目完全沒有標題，只能整段讀完才知道重點。兩處都補上一句話標題，報告才能被快速掃讀。

一句話標題若只寫進 Markdown 慣例，API 無法驗證每段是否真的有標題，版面也無從把標題做成獨立小標。因此改為結構化欄位，由 schema 強制、由版面統一呈現。

## What Changes

- **BREAKING** `industries[]` 的 `summary_md` 改為 `events[]` 陣列，每個事件含 `headline`（一句話標題）與 `summary_md`（內文）。產業層級的 `watch_md` 不變，仍在分類層級。
- **BREAKING** `stock_news[]` 新增 `headline` 欄位，一句話標題。
- API 驗證擴充：每個 `industries` 條目至少一個事件；每個事件的 `headline` 與 `summary_md` 非空白；每個 `stock_news` 條目的 `headline` 非空白。任一違規回 400 且整份不落庫。
- 版面調整：產業分類內逐一渲染事件，標題為獨立小標題、內文接在其下；個股條目的標題渲染在股票名稱與代號所在的識別欄位內。
- `skills/daily-brief/SKILL.md` 規範一句話標題寫法（寫出發生什麼事與影響，不用短語標籤、不重複內文），並更新報告 JSON 範例與送出前檢查清單。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `report-ingest-api`: 報告 schema 改為事件陣列並新增標題驗證。
- `report-viewing`: 產業事件標題與個股標題的版面呈現。
- `daily-brief-skill`: 一句話標題的撰寫規範與新的 industries 結構。

## Impact

- Affected specs: `report-ingest-api`, `report-viewing`, `daily-brief-skill`
- Affected code:
  - New: （無）
  - Modified: internal/report/schema.go, internal/report/schema_test.go, internal/api/report_test.go, internal/api/read_test.go, internal/site/site.go, internal/site/site_test.go, internal/site/styles.go, skills/daily-brief/SKILL.md
  - Removed: （無）
- 資料庫結構不變：報告以完整 JSON 儲存，欄位調整不需要 migration。
- 網站尚未上線、無已發布報告，因此不保留舊 `industries[].summary_md` 格式的相容處理。
