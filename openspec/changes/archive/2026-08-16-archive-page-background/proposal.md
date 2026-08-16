## Why

歷史報告的提示帶捲出畫面後，頁面中段與最新報告完全相同，讀者仍可能把舊內容當最新內容讀。已與使用者確認採整頁換底色：歷史頁改用「舊報紙」的檔案紙色，任何捲動位置都能辨識。

## What Changes

- 從歷史列表開啟且日期早於最新報告時，整頁紙底改用檔案紙色：淺色主題用更灰黃的舊紙色，深色主題用對應的沉色調，透過覆寫 paper 色票達成，不逐一改元件樣式。
- 最新報告（含從歷史列表點到最新日期）維持原紙色，不受影響。
- 新底色上所有既有文字組合的對比度實測仍 ≥4.5:1 才採用，紅綠多空色與其 tint 不變。
- DESIGN.md 的 Colors 一節登錄新的檔案紙色 token 與使用規則，避免設計文件與實作漂移。
- 既有的歷史提示帶保留不動，底色與提示帶並存。

## Non-Goals

- 不做 sticky 提示帶；已與使用者確認先只做底色，避免兩個機制疊加的過度設計。
- 不改歷史列表頁與 404 頁的底色；只有「檢視某份舊報告」的頁面換色。
- 不引入紅綠以外的新彩色系；檔案紙色是中性紙色的深淺變化，不是新色相。

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `report-viewing`: History browsing 需求擴充——早於最新日期的歷史報告檢視必須整頁使用可辨識的檔案紙底色（深淺主題各一組、文字對比 ≥4.5:1），最新日期不得套用。

## Impact

- Affected specs: `report-viewing`
- Affected code:
  - Modified: internal/site/site.go, internal/site/styles.go, internal/site/site_test.go, DESIGN.md
  - New: (none)
  - Removed: (none)
