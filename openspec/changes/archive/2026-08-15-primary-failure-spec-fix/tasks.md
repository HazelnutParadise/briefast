## 1. 驗證修正後 spec 與現行 SKILL.md 一致

- [x] 1.1 逐句比對修正後的 Batched collection with per-batch counts requirement 與 skills/daily-brief/SKILL.md「來源成敗與完整性把關」「分批蒐集」兩節，確認行為完全一致：任何批次（含鉅亨主來源）重試後仍失敗都照常發布、缺漏只列在執行回報、彙整表仍是判讀前置條件。SKILL.md 預期零修改。驗證：內容檢視比對結果逐項列出，且 git status 顯示 skills/daily-brief/SKILL.md 無變更。
- [x] 1.2 確認修正後 requirement 與 Source collection completeness gate 的失敗處理敘述無任何殘留矛盾：本 change 的 delta spec 不存在「primary source batch 失敗停止流程」的規範句或 scenario（正式 spec 於封存合併時才移除該句，封存後需複查同一檢索條件）。驗證：對 openspec/changes/primary-failure-spec-fix/specs/daily-brief-skill/spec.md 檢索 stop the workflow 與 stops the run 字樣，僅允許出現在否定敘述（does not stop）中；並確認正式 spec 中該句仍存在、將由本 change 封存時移除。

## 2. 同步修正停放中的 foreign-news-batch

- [x] 2.1 以 spectra unpark foreign-news-batch 取回停放變更，將其 delta spec（openspec/changes/foreign-news-batch/specs/daily-brief-skill/spec.md）中 Batched collection with per-batch counts 的「A failure of the primary source batch SHALL stop the workflow」規範句與「Primary failure stops the run」scenario，改為與本次修正相同的處理：主來源批次失敗照常發布並列入執行回報；完成後以 spectra park foreign-news-batch 重新停放。驗證：執行 spectra validate foreign-news-batch 通過，且該檔案檢索不到 stop-the-workflow 的肯定敘述。
- [x] 2.2 確認兩個變更對同一 requirement 的封存順序安全：primary-failure-spec-fix 的 delta 是四批版本，foreign-news-batch 的 delta 是五批版本，後者必須在前者之後封存，否則四批版本會覆蓋五批擴充。驗證：在本 change 完成回報與 foreign-news-batch 的 proposal 內容中皆記明「先封存 primary-failure-spec-fix、後套用 foreign-news-batch」的順序約束，內容檢視確認兩處都有寫。
