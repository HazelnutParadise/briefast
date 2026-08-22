## 1. 框架設定

- [x] 1.1 交付 `Public pages declare Traditional Chinese` 中屬於框架的部分：在 `syralit.toml` 加 `lang = "zh-Hant-TW"`，讓首頁、歷史頁與後台的第一個 HTML 回應都帶該語言標示。驗證方式：啟動伺服器後 curl 三個路徑，確認 HTML 根元素皆為 `lang="zh-Hant-TW"`，後台不再是 `lang="en"`。

## 2. 中介軟體簡化

- [x] 2.1 先改測試：把中介軟體既有的語言替換斷言，換成驗證 `Rewriting layer leaves the language attribute alone`，也就是餵進一份帶 `lang="en"` 的假 shell 時，改寫後該屬性維持原樣，只有 title 與中繼資料被替換。驗證方式：`go test ./internal/seo/` 因為實作仍在替換而失敗。
- [x] 2.2 移除 `rewriteHead` 內替換語言標示的那一步，改寫流程只剩定位並替換 title 元素。驗證方式：2.1 的測試轉綠，`go test ./internal/seo/` 全數通過。
- [x] 2.3 交付 `Language survives a declined rewrite`：確認 HTML 內沒有 title 元素時整份原樣放行，而語言標示仍由框架提供而正確。驗證方式：在 `main_test.go` 以真實 mux 取得後台頁面，斷言其帶有 `lang="zh-Hant-TW"` 且不含 canonical 標籤。

## 3. 文件與收尾

- [x] 3.1 更新 `AGENTS.md` 的中介軟體約定，改成語言標示由 `syralit.toml` 的 `lang` 提供、中介軟體只負責每頁不同的 title 與中繼資料，並移除已不成立的語言替換描述。驗證方式：閱讀該條約定，確認與實際實作一致且沒有殘留舊敘述。
- [x] 3.2 確認整包升級後行為無回歸：`go vet ./...` 與 `go test ./...` 全綠，並在瀏覽器實際開啟首頁確認 WebSocket 即時更新與報告內容顯示正常。驗證方式：保留指令輸出與瀏覽器截圖作為憑據。
