## 1. SKILL.md 批次 5 明定 CNBC 請求條件

- [x] 1.1 在 skills/daily-brief/SKILL.md 的「批次 5：國外財經新聞（CNBC）」段落寫明：RSS 清單與每篇文章頁的請求都必須帶桌面瀏覽器 User-Agent，並附上一組可直接複製的 UA 字串與 curl 範例（含 --user-agent），同時說明原因是 CNBC 由 Akamai 依 User-Agent 封鎖 curl、python-requests、Claude-User 這類工具型客戶端。落實 Pinned source endpoints requirement 的 CNBC 瀏覽器 UA 條款。驗證：以 SKILL.md 寫的 UA 與 curl 範例實際請求 RSS 端點回 200 且解析出 item；再以同一 UA 請求 RSS 內第一篇文章頁回 200 且能取出正文段落；並確認以 curl 預設 UA 請求同一端點仍回 403，證明差異來自 UA。
- [x] 1.2 在同一段落與「來源成敗把關」段落加入 403 檢查步驟：CNBC 回 403 時先確認請求是否帶了瀏覽器 UA，沒帶就補上重送，帶了仍 403 或連線失敗才算來源失敗，走既有「重試一次後照常發布、執行回報列缺項」的規則；並說明 RSS 的 description 只有一句摘要，不能取代文章頁抓取。落實 Pinned source endpoints requirement 的 403 先查標頭條款與 feed 摘要不可替代條款。驗證：內容檢視 SKILL.md，確認批次 5 段落與把關段落都能找到「403 先查 UA」的分支，且把關段落對 CNBC 的失敗定義已限定為帶瀏覽器 UA 後仍失敗。

## 2. Spec 與 SKILL.md 一致性

- [x] 2.1 逐句對照 openspec/changes/cnbc-browser-user-agent/specs/daily-brief-skill/spec.md 的 Pinned source endpoints requirement 三個 scenario（CNBC requests carry a browser User-Agent、CNBC 403 is checked against the header before being recorded as failure、既有兩個 scenario）與 SKILL.md 批次 5 段落，確認每個 scenario 都能在 SKILL.md 找到對應規則文字，批次 1 到 4 的抓取敘述未被改動。驗證：spectra validate cnbc-browser-user-agent 通過，且對 SKILL.md grep「User-Agent」只命中批次 5 與把關段落。
