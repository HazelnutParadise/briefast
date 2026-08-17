## 1. 改樣式與驗證

- [x] 1.1 提示帶改深色帶樣式：internal/site/styles.go 的 .archive-note 改 band 深底、文字淺色 14px、標籤淺底深字、連結淺色底線；internal/site/site_test.go 補樣式斷言（.archive-note 使用 var(--band)）。驗證方式：`go test ./internal/site/` 通過。
- [x] 1.2 視覺確認：本機啟動開啟舊日期報告，淺色與深色截圖確認深帶明顯、文字與連結清楚可讀。驗證方式：截圖核對。
