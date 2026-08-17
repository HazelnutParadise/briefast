## Summary

歷史提示帶從紙面上的細字橫幅改為深色帶（band）樣式：深底、淺色粗體文字、底線連結，大幅提高可見度。

## Motivation

使用者反映現行橫幅不夠明顯。現況只有小型反白標籤加 13.5px 灰字，在寬螢幕與快速捲讀下容易被略過。設計語言中不動用紅綠（保留給多空語意）時，最強的中性強調就是既有的深色帶（頁尾 band 同款）。

## Proposed Solution

- `.archive-note` 背景改 `var(--band)`，整條成為深帶；文字用頁尾同款淺色（`#F2EDE6`／band-ink），字級升到 14px。
- 「歷史報告」標籤改為淺底深字（在深帶上反轉），連結改淺色加底線。
- 內容與連結文字不變；只改視覺。
- 對比：淺字對 band 深底遠高於 4.5:1（頁尾同組合已在用）。

## Non-Goals

- 不改提示帶的出現條件、文案與連結。
- 不使用紅綠或新增色相；band 與淺字都是既有 token 與既有用色。

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

(none)

## Impact

- Affected specs: (none)（spec 只要求 prominent notice，未鎖定樣式）
- Affected code:
  - Modified: internal/site/styles.go, internal/site/site_test.go
  - New: (none)
  - Removed: (none)
