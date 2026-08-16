## 1. 換色與驗證

- [x] 1.1 檔案紙色換為暖調淡灰並保持對比合規：internal/site/styles.go 的 archiveStyles 淺色 paper 改 `#E4E0D9`（link 維持 `#0B6870`）、深色 paper 改 `#282623`；DESIGN.md 檔案紙色數值同步；internal/site/site_test.go 色值斷言同步。驗證方式：對比度實算數據（淺色最低 4.71、深色最低 4.78）附回報，`go test ./internal/site/` 通過。
- [x] 1.2 視覺確認：本機啟動開啟歷史列表與歷史報告，淺色與深色截圖確認灰調乾淨、與首頁仍有明顯區隔。驗證方式：截圖核對。
