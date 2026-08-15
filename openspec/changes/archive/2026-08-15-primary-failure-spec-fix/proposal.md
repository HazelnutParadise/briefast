## Problem

openspec/specs/daily-brief-skill/spec.md 有兩條 requirement 對「主來源鉅亨（cnyes）失敗時要不要停發報告」規定相反：Source collection completeness gate 要求任何來源（含主來源）重試後仍失敗都照常發布、缺漏只列在執行回報；Batched collection with per-batch counts 卻寫主來源批次失敗要停止整個流程，並附「Primary failure stops the run」scenario。兩條同時有效，實作者無所適從。

## Root Cause

batched-source-collection 變更（2026-08-13 封存）的提案同時包含兩種敘述：一處寫「主來源失敗仍中止整個流程」，另一處明確決定「移除鉅亨網失敗即整份不發的規則，改為所有來源失敗都照常發布」並附完整理由（證交所重大訊息帶公司代號、月營收是硬事實，其他台媒也涵蓋台股，缺 cnyes 損失的是覆蓋密度而非可行性）。最終 SKILL.md 依後者實作，但 delta spec 的 Batched collection requirement 沒同步清掉前者的殘留句，封存時把矛盾帶進了正式 spec。

## Proposed Solution

修改 Batched collection with per-batch counts requirement：刪除「A failure of the primary source batch SHALL stop the workflow regardless of what the other batches returned」一句，改寫為批次失敗一律依 Source collection completeness gate 處理；把「Primary failure stops the run」scenario 改為主來源批次失敗後照常發布、缺漏列入執行回報的版本。其餘內容（五秒節流、彙整表關卡、平行條款）不動。SKILL.md 不需修改，現行內容即為正確行為。

## Non-Goals

- 不改 skills/daily-brief/SKILL.md：其「任何單一來源失敗都不停發」的現行行為就是本次要對齊的基準。
- 不重新檢討「主來源失敗要不要停發」這個決策本身：batched-source-collection 已做過該決策並留有理由，本次只清除殘留。
- 不處理 parked change foreign-news-batch 以外的其他變更；該 parked change 的 delta spec 複製了同一句殘留文字，需一併更正，但其五批擴充內容不在本次範圍。

## Success Criteria

- 正式 spec 中 Batched collection with per-batch counts 不再含任何主來源失敗停止流程的敘述與 scenario。
- 該 requirement 的失敗處理與 Source collection completeness gate 完全一致：任何批次失敗都不停發，缺漏列入執行回報。
- parked change foreign-news-batch 的 delta spec 同步更新，不會在日後封存時把殘留句寫回正式 spec。
- spectra validate primary-failure-spec-fix 通過，且逐句比對 SKILL.md 來源成敗把關一節與修正後 requirement 無矛盾。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `daily-brief-skill`: Batched collection with per-batch counts 的失敗處理從「主來源批次失敗停止流程」改為「所有批次失敗一律依 completeness gate 照常發布並於執行回報揭露」。

## Impact

- Affected specs: `daily-brief-skill`
- Affected code:
  - Modified: openspec/changes/foreign-news-batch/specs/daily-brief-skill/spec.md（該變更目前停放中，需先以 spectra unpark 取回、修正後重新停放；此為取回後的路徑）
  - New: （無）
  - Removed: （無）
