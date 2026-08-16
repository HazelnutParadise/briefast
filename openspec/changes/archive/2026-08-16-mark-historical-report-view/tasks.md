## 1. Store 層

- [x] 1.1 交付 History browsing 需求中「取得最新報告日期不得載入整份 payload」的查詢：在 internal/store/store.go 新增 LatestReportDate 方法，只 SELECT 日期欄位，無報告時回傳 store.ErrNotFound。驗證方式：internal/store/store_test.go 新增測試，寫入多日報告後斷言回傳最大日期，空庫斷言 ErrNotFound，`go test ./internal/store/` 通過。

## 2. Site 層

- [x] 2.1 交付 History browsing 需求的歷史提示帶與導覽：internal/site/site.go 的 HistoryPage 在 date 早於 LatestReportDate 時，於報頭下方渲染載明日期、寫明非最新內容、含「看最新報告」（/）與「返回歷史列表」（/history/）連結的提示帶；date 等於最新日期時不渲染提示帶；歷史報告檢視的報頭導覽同時含回首頁與歷史報告連結；404 空頁補回首頁連結；提示帶樣式加入 internal/site/styles.go，遵循 DESIGN.md（零圓角、零陰影、既有色票與規線語彙）。驗證方式：internal/site/site_test.go 以 sy.NewAppTest 斷言——舊日期含提示帶與兩個連結、最新日期不含提示帶、404 頁含回首頁連結，`go test ./internal/site/` 通過。
- [x] 2.2 全套驗證與視覺確認：`go test ./...` 全綠；本機啟動後實際開啟首頁、歷史列表、舊日期報告與最新日期報告，截圖確認提示帶明顯可辨、版面未跑版、深色模式配色正常。驗證方式：測試輸出與截圖各狀態核對。
