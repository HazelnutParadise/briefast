## 1. SKILL.md 移除 WSJ

- [x] 1.1 把 skills/daily-brief/SKILL.md 批次 5 小節改為只含 CNBC：標題與說明不再提 WSJ，來源表刪除 WSJ 列，計數回報改為單一來源的取得與未讀筆數，落實 Batched collection with per-batch counts 的 batch five 新涵蓋範圍與 Pinned source endpoints 的來源集。驗證：對 SKILL.md 全文 grep WSJ 與 feeds.a.dj.com 確認零殘留，批次 5 段落與 delta spec 的 Batch five 句一致。
- [x] 1.2 同步清理 SKILL.md 其他提及 WSJ 的位置：來源成敗把關的外電括號註記、去重章節的 WSJ 摘要比對特例句、判讀原則與全球市場動態敘述若有提及一併改為僅指 CNBC，外電的範圍限制與跨語言去重行為文字維持不變，落實 Foreign wire scope and cross-language deduplication 對單一來源的適用。驗證：內容檢視三處（把關、去重、判讀原則）都只指涉 CNBC，且範圍與去重規則的行為敘述未被改動。

## 2. 一致性驗證

- [x] 2.1 全文比對修改後的 SKILL.md 與 delta spec 三條 requirement（Pinned source endpoints、Batched collection with per-batch counts、Foreign wire scope and cross-language deduplication）逐條一致，批次數維持五批、其餘批次與昨收價規則零變動。驗證：執行 spectra validate remove-wsj-source 通過，git diff 確認變更僅限外電相關段落。
