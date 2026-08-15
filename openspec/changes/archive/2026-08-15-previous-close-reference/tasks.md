## 1. SKILL.md 新增昨收價蒐集

- [x] 1.1 在 skills/daily-brief/SKILL.md 第 1 步新增「昨收價參考資料」小節，落實 Previous-close reference data 的抓取契約：兩個端點（openapi.twse.com.tw 的 exchangeReport/STOCK_DAY_ALL、www.tpex.org.tw 的 openapi/v1/tpex_mainboard_daily_close_quotes）各一個一般 HTTP 請求，說明證交所請求併入既有 TWSE 序列節流、櫃買主機不受該節流限制，並要求回報兩來源各幾筆與資料日期。驗證：以 curl 實測兩端點回 200 且回應含民國格式日期欄位，逐句比對小節內容與 delta spec 第一段一致。
- [x] 1.2 在同一小節寫入日期驗證規則：把回應內嵌的民國日期換算西元後與台北時區前一交易日比對，過舊即視為該來源抓取失敗，不得沿用舊價且不標示。驗證：內容檢視確認規則涵蓋「民國格式換算」與「過舊即失敗」兩點，與 delta spec 的 Stale date treated as failure scenario 對應。
- [x] 1.3 在來源成敗把關段落納入昨收抓取的失敗分流：每來源重試一次，仍失敗就改用純新聞脈絡照常判讀與發布，缺漏列入執行回報、不進報告內容，落實 Previous-close reference data 的不停發行為。驗證：內容檢視該段落與既有次要來源失敗處理寫法一致，且送出前檢查清單相應更新。

## 2. 判讀規則接入價格脈絡

- [x] 2.1 在 SKILL.md 判讀章節寫入昨收價用途規則：昨收僅作判讀脈絡（利多是否已反映、量級對照）與報告引用來源，引用時必須標明價格日期；多空 call 仍以窗口內新聞為依據，禁止只憑價格漲跌建立 calls 或 stock_news 條目，落實 Previous-close reference data 的用途限制。驗證：內容檢視規則同時出現在個股判讀原則與送出前檢查清單，並與基本面引用紀律（僅限本次抓取資料、標明資料期間）無矛盾。

## 3. 一致性驗證

- [x] 3.1 全文比對修改後的 SKILL.md 與 Previous-close reference data requirement 的三個 scenario 逐一有對應規則文字，且未動到分批蒐集章節的批次結構敘述。驗證：執行 spectra validate previous-close-reference 通過，並以接手者視角確認與停放中的 foreign-news-batch、primary-failure-spec-fix 無重疊修改點。
