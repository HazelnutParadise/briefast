## 1. 選色與對比驗證

- [x] 1.1 決定檔案紙色並證明合規：為淺色與深色主題各選一個中性檔案紙色（不引入新色相），以腳本實算 ink、ink-soft、link 與紅綠 up/down 在新底色上的 WCAG 對比度，全部 ≥4.5:1 才採用；同時把選定值與使用規則寫進 DESIGN.md 的 Colors 一節。驗證方式：對比度計算輸出附在執行回報，DESIGN.md 內容審閱含新 token 與「僅用於歷史報告檢視」的規則。

## 2. 實作與驗證

- [x] 2.1 交付 History browsing 需求的 Historical view uses archive paper background 行為：internal/site/site.go 在 date 早於最新日期時，於頁面輸出附加覆寫 paper token 的樣式區塊（淺深主題各一組值，值來自 internal/site/styles.go 或同檔常數），date 等於最新日期時不輸出；歷史提示帶維持不變。驗證方式：internal/site/site_test.go 斷言舊日期輸出含 paper 覆寫與兩組主題值、最新日期不含，`go test ./internal/site/` 通過。
- [x] 2.2 全套驗證與視覺確認：`go test ./...` 全綠；本機啟動後開啟舊日期報告與最新日期報告，淺色與深色都截圖，確認整頁底色可辨、與最新報告有明顯差異、文字可讀、紅綠標籤不受影響。驗證方式：測試輸出與截圖逐狀態核對。
