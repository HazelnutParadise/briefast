## 1. SKILL.md 蒐集流程擴為五批

- [x] 1.1 在 skills/daily-brief/SKILL.md 的分批蒐集章節新增「批次 5：國外財經新聞」，釘選 WSJ Markets RSS（feeds.a.dj.com/rss/RSSMarketsMain.xml）與 CNBC 頭條 RSS（www.cnbc.com/id/100003114/device/rss/rss.html）兩個端點，並要求本批回報兩個來源各取得幾則、其中未讀幾則，落實 Pinned source endpoints 與 Batched collection with per-batch counts 的五批版要求。驗證：以 curl 實際請求兩個 RSS 端點確認回 200 且為有效 RSS，並逐句比對批次 5 段落與 delta spec 的 Batch five 描述一致。
- [x] 1.2 把 SKILL.md 中所有「四批」的流程敘述改為五批：分批蒐集開頭、平行執行說明、彙整表段落（含標題「四批彙整表」）、送出前檢查清單的彙整表條目，並把「全球市場動態……不另設來源」一句改寫為指向批次 5 與台媒報導並用，使 Batched collection with per-batch counts 的判讀前置條件涵蓋五批。驗證：對 SKILL.md 全文 grep「四批」與「不另設來源」確認無殘留，且送出前檢查清單的彙整表條目已寫明五個批次。
- [x] 1.3 在 SKILL.md 的來源成敗把關段落明定批次 5 失敗比照次要來源處理：重試一次後仍失敗就照常發布、只在執行回報列出，落實 Batched collection with per-batch counts 中 batch five 失敗不停發的行為。驗證：內容檢視該段落，確認批次 5 失敗路徑與既有次要來源（如 TWSE）敘述一致，未被歸入主來源停發規則。

## 2. 外電範圍限制與跨語言去重

- [x] 2.1 在 SKILL.md 新增外電使用範圍規則：批次 5 內容只寫入 overview_md 與產業事件背景，不得建立或支持 calls 與 stock_news 條目，落實 Foreign wire scope and cross-language deduplication 的範圍限制。驗證：內容檢視確認規則同時出現在批次 5 段落與判讀原則區，且與既有 wd_stock「只寫入 overview_md」的寫法一致。
- [x] 2.2 在 SKILL.md 去重章節新增跨語言去重規則：seen.py similar 不做跨語言比對，外電文章須逐則對照同窗口台媒報導，事件已被台媒報導者判定不報並以 seen.py record 記為 skipped、以台媒版本為準，僅台媒未報的事件可進報告，落實 Foreign wire scope and cross-language deduplication 的去重行為。驗證：內容檢視去重章節，確認外電流程涵蓋「台媒已報→skipped」與「僅外電有→可用」兩個分支，與 delta spec 兩個 scenario 對應。

## 3. 一致性驗證

- [x] 3.1 全文檢查 SKILL.md 修改後與 delta spec 的三條 requirement（Foreign wire scope and cross-language deduplication、Pinned source endpoints、Batched collection with per-batch counts）逐條對應，無矛盾或漏項。驗證：執行 spectra validate foreign-news-batch 通過，並以接手者視角逐 scenario 對照 SKILL.md 找得到對應規則文字。
