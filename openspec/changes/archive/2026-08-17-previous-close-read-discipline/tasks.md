## 1. SKILL.md 讀取紀律

- [x] 1.1 在 skills/daily-brief/SKILL.md「昨收價參考資料」一節加入讀取紀律段落：兩個回應抓回後存成工作資料夾檔案；判讀時只用 jq 或 grep 等查詢工具，按新聞提及的代號取值；不得把整份全市場回應讀進對話 context；日期驗證與計數回報同樣以查詢存檔取得（讀 Date 欄位與筆數），不整份讀入。完成後該節同時涵蓋抓取、節流、日期驗證、計數回報與讀取紀律五件事，且與 spec 的 Previous-close reference data requirement 敘述一致。驗證：內容檢視該節，逐句對照 change 的 specs/daily-brief-skill/spec.md 中新增的讀取紀律段落與 Prices queried from saved files by symbol scenario，確認無遺漏也無矛盾。

## 2. 一致性驗證

- [x] 2.1 確認 SKILL.md 其他段落沒有與讀取紀律矛盾的指引（例如要求整份檢視回應的敘述），蒐集完成度回報的要求維持不變。驗證：以「整份」「全市場」「STOCK_DAY_ALL」為關鍵字通讀 skills/daily-brief/SKILL.md 相關段落，確認只有昨收價一節談讀取方式且口徑一致；執行 spectra validate previous-close-read-discipline 通過。
