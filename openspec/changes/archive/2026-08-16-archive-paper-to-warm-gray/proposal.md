## Summary

歷史區檔案紙色從黃綠調（淺 `#EAE0CA`／深 `#2A251E`）改為暖調淡灰（淺 `#E4E0D9`／深 `#282623`），其餘行為不變。

## Motivation

使用者反映現行檔案紙色的灰綠感難看。改為帶暖調的淡灰：仍與首頁奶油紙明顯區隔，但不再有黃綠色偏；純冷灰會與全站暖紙面相衝，故取 greige 而非純灰。

## Proposed Solution

- `internal/site/styles.go` 的 archiveStyles 換色：淺色 paper `#E4E0D9`（連結青維持加深的 `#0B6870`，對比 4.94），深色 paper `#282623`。
- 對比度已實算：淺色最低 ink-soft 4.71、深色最低漲紅 4.78，全部 ≥4.5:1。
- DESIGN.md 檔案紙色數值同步更新；site_test.go 斷言的色值同步更新。
- 適用範圍與提示帶邏輯完全不動；不涉及 spec 需求變更（spec 未鎖定色值）。

## Non-Goals

- 不改適用範圍、提示帶、導覽與任何行為。
- 不採純冷灰（與暖紙面相衝）。

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

(none)

## Impact

- Affected specs: (none)
- Affected code:
  - Modified: internal/site/styles.go, internal/site/site_test.go, DESIGN.md
  - New: (none)
  - Removed: (none)
