## Why

從歷史列表點進某一天的報告後，頁面與首頁的最新報告完全相同（版面、導覽都一樣），讀者無法辨識自己看的是舊報告，也沒有直接回首頁的連結（頁首導覽只有「歷史報告 →」）。舊報告被當成最新內容讀，對一份有時效性的盤前晨報是實質誤導。

## What Changes

- 從歷史列表開啟的報告，若其日期早於資料庫中最新報告的日期，在報頭下方顯示明顯的歷史報告提示帶：載明「這是 {日期} 的歷史報告，非最新內容」，並提供「看最新報告」（回首頁）與「返回歷史列表」兩個連結。
- 若選到的日期就是最新報告的日期，不顯示歷史提示（它就是最新內容），照首頁方式呈現。
- 歷史報告檢視頁的報頭導覽同時提供「回首頁」與「歷史報告」連結；歷史報告找不到（404）的空頁也補上回首頁連結。
- `internal/store` 新增查詢最新報告日期的輕量方法（只取日期，不載入整份報告 JSON），供比對使用。
- 提示帶樣式遵循 DESIGN.md：零圓角、零陰影、規線與既有色票，不引入新的視覺語彙。

## Non-Goals

- 不改動歷史列表頁本身的排序、分頁與版面。
- 不改動報告五區塊版面；歷史報告內容仍用與首頁相同的 renderReport 呈現。
- 不在 URL 或路由結構上做變更（仍是 `/history/?date=...`）。
- 拒絕的替代方案：一律對從歷史列表進入的報告顯示「歷史報告」標示、不比對最新日期——歷史列表第一列就是最新報告，對它標「非最新內容」是錯的。

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `report-viewing`: History browsing 需求擴充——歷史報告檢視必須可辨識為非最新內容並提供回首頁與返回列表的導覽；與最新報告同日期時不得誤標。

## Impact

- Affected specs: `report-viewing`
- Affected code:
  - Modified: internal/site/site.go, internal/site/styles.go, internal/site/site_test.go, internal/store/store.go, internal/store/store_test.go
  - New: (none)
  - Removed: (none)
