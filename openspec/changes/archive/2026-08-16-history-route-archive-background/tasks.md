## 1. 實作與驗證

- [x] 1.1 交付 History browsing 需求中 Whole history section uses archive paper background 行為：internal/site/site.go 的 HistoryPage 三個分支（列表、報告檢視、404）一律輸出 archiveStyles，提示帶維持日期判斷；DESIGN.md 檔案紙色規則改為「用於 /history/ 區全部頁面，標示區域而非內容新舊」。驗證方式：internal/site/site_test.go 斷言列表頁、最新日期檢視、404 頁都含 paper 覆寫，最新日期檢視仍不含提示帶，首頁不含覆寫；`go test ./internal/site/` 通過。
- [x] 1.2 全套驗證與視覺確認：`go test ./...` 全綠；本機啟動後開啟首頁、歷史列表、最新日期與舊日期檢視，確認歷史區一律舊紙底、首頁維持原紙色。驗證方式：測試輸出與截圖核對。
