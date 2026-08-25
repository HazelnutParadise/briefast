## 1. 讓測試先抓到問題

- [x] 1.1 在 `main_test.go` 新增測試，交付 `Listen address is resolved from configuration` 的驗證：斷言應用程式算出的監聽位址等於 `0.0.0.0:8600`，也就是專案 syralit.toml 宣告的值。驗證方式：`go test ./...` 因為目前算出 `:0` 而失敗，且失敗訊息顯示實際值。

## 2. 修正位址來源

- [x] 2.1 在 `main.go` 收斂出單一設定來源，改用 `sy.ResolveConfig` 取得套用 syralit.toml 與內建預設後的 `Config`，讓建立 handler 與計算監聽位址讀到同一份結果，不再於 page function 外呼叫 `sy.GetOption`。並把位址計算抽成獨立函式，讓 1.1 的測試可以直接呼叫。此修正落實 `Configuration injection` 中位址必須來自解析後設定的規定，且不改變設定的存放位置。驗證方式：1.1 的測試轉綠，`go test ./...` 與 `go vet ./...` 全數通過。

## 3. 端對端確認

- [x] 3.1 交付 `Image runs standalone` 情境中屬於本次修正的部分：實際啟動伺服器，確認 log 顯示 `Briefast listening on http://0.0.0.0:8600`，並確認公開頁、後台、robots.txt 與 sitemap.xml 的回應與升級前一致。驗證方式：保留啟動 log 與各路徑的狀態碼輸出作為憑據。
