## 1. 報告 schema 支援 chips

- [x] 1.1 依 design.md「chips 資料形狀與驗證規則」在 `internal/report/schema_test.go` 先寫失敗測試，覆蓋 Report schema validation 的三態：chips 合法通過、chips 缺漏或為 null 通過、chips.date 非法時回傳標明 stock_news 索引的錯誤訊息。驗證：`go test ./internal/report/` 出現預期失敗。
- [x] 1.2 在 `internal/report/schema.go` 為 `StockNews` 加上 `Chips *Chips` 欄位與 `Chips` 型別（date、foreign_net、trust_net、dealer_net、total_net，以及指標型別的 margin_change、short_change），並在 `Validate` 檢核 chips.date。行為：帶合法 chips 的報告通過驗證，chips.date 非法者回傳指名索引的錯誤，不驗證 total_net 是否等於三項加總。驗證：1.1 的測試全綠。
- [x] 1.3 在 `internal/api/report_test.go` 補一則 ingest 往返測試，行為：POST 帶 chips 的報告回 200，`GET /api/report/{date}` 讀回完全相同的 chips 物件；另一則 POST chips.date 非法回 400 且不落庫。驗證：`go test ./internal/api/` 全綠。

## 2. 個股詳情籌碼視覺化

- [x] 2.1 依 design.md「個股詳情籌碼視覺化版面」在 `internal/site/site_test.go` 先寫失敗測試，覆蓋 Chip data display：帶 chips 的條目渲染出「籌碼面」區塊、三列長條寬度按最大絕對值等比、買超用多方紅、賣超用空方綠、股數換算張且正值帶加號；無 chips 的條目不出現該區塊也不出現佔位文字。色彩斷言解析色值色相，不只比對選擇器字串。驗證：`go test ./internal/site/` 出現預期失敗。
- [x] 2.2 在 `internal/site/site.go` 的 stock_news 條目渲染中，於觀察重點之後、來源列之前輸出籌碼面區塊，寬度百分比在伺服器端算好以行內樣式輸出，並在 `internal/site/styles.go` 加上零圓角零陰影的長條樣式。行為：帶 chips 的條目顯示三列法人長條與合計、融資、融券、資料日期摘要行，margin_change 或 short_change 缺漏時省略對應片段。驗證：2.1 的測試全綠。
- [x] 2.3 啟動本機服務並以瀏覽器檢視塞入測試報告的首頁與歷史頁，比對版面與 DESIGN.md 的規線、色彩與間距慣例。行為：籌碼區塊在深淺兩種配色下都可讀，紅漲綠跌未反轉，與既有條目排版一致。驗證：兩種配色各一張截圖並逐項比對。

## 3. 晨報 skill 籌碼抓取與判讀

- [x] 3.1 依 design.md「沿用昨收價參考資料模式抓取籌碼資料」在 `skills/daily-brief/SKILL.md` 新增 Chip reference data 條文：四個端點、兩個 www.twse.com.tw 請求併入既有 TWSE 節流序列、回應存檔後以 jq/grep 按代號查值不整份進 context、民國日期換算後須等於前一交易日否則視為抓取失敗、每來源回報筆數與資料日期、失敗照常發布並列入執行回報。驗證：對照 `openspec/changes/add-stock-chip-data/specs/daily-brief-skill/spec.md` 的四個情境逐條核對條文涵蓋。
- [x] 3.2 在 `skills/daily-brief/SKILL.md` 的 JSON 結構與送出前檢查清單補上 Chip block composition 規則：chips 欄位形狀與範例、foreign_net 為外資兩欄相加、total_net 直接取來源合計、margin_change 與 short_change 為今日餘額減前日餘額（張），查無代號時省略該欄位或整塊，市場抓取失敗時該市場條目整塊省略。驗證：對照 spec 的三個情境核對，並確認 JSON 範例與 `internal/report/schema.go` 的欄位名一致。
- [x] 3.3 依 design.md「判讀引用紀律比照昨收價」在判讀原則與檢核清單補上 Chip citation discipline in judgement：籌碼只作既有新聞判斷的佐證、不得單獨成立 calls 或 stock_news 條目、正文引用須標明資料日期、缺料條目不得出現佔位語句。驗證：對照 spec 的兩個情境核對條文。

## 4. 收斂驗證

- [x] 4.1 全專案測試與交付前檢查。行為：`go test ./...` 全綠，`go vet ./...` 無輸出，舊報告（無 chips）渲染與 ingest 行為不變。驗證：兩道指令的實際輸出。
- [x] 4.2 執行 `spectra validate add-stock-chip-data` 與 `spectra analyze add-stock-chip-data`，確認 artifacts 一致無缺口。驗證：兩道指令回報 valid 且無未覆蓋的 requirement。
