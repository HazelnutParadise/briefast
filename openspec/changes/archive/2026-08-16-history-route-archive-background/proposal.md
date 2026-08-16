## Why

目前檔案紙底色只在「日期早於最新報告」時出現，從歷史列表點到最新一份時頁面又變回一般紙色，歷史區的視覺身分不連貫。使用者指定改成：`/history/` 底下一律用歷史背景，底色代表「所在區域」而非「內容新舊」。

## What Changes

- `/history/` 下所有頁面一律套用檔案紙底色：歷史列表頁、任何日期的報告檢視（含最新日期）、找不到報告的頁面。
- 歷史提示帶維持原本的日期判斷不變：只有日期早於最新報告才顯示「並非最新內容」，點到最新那份不顯示，避免不實敘述。
- 首頁與其他頁面不受影響。
- DESIGN.md 的檔案紙色使用規則同步改為「用於 /history/ 區全部頁面」。

## Non-Goals

- 不改提示帶的判斷與文案；底色與提示帶自此語意分工——底色標區域、提示帶標內容新舊。
- 不改色票值；沿用已驗證對比的既有檔案紙色與加深連結青。

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `report-viewing`: History browsing 的檔案紙底適用範圍，從「早於最新日期的檢視」改為「/history/ 路由下的所有頁面」；提示帶規則不變。

## Impact

- Affected specs: `report-viewing`
- Affected code:
  - Modified: internal/site/site.go, internal/site/site_test.go, DESIGN.md
  - New: (none)
  - Removed: (none)
