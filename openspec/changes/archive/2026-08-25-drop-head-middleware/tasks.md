## 1. 中繼資料改由 DocumentFunc 產生

- [x] 1.1 先改測試：把 `internal/seo/seo_test.go` 中透過中介軟體驗證改寫結果的測試，改為直接呼叫依頁面種類建立的 `DocumentFunc`，斷言回傳的 `Document` 帶有正確的 `Title` 與 `HeadHTML`，涵蓋首頁取最新報告、歷史列表用固定文案、`?date=` 取該日報告三種情況。驗證方式：`go test ./internal/seo/` 因建構器尚未存在而編譯失敗。
- [x] 1.2 在 `internal/seo` 新增依頁面種類建立 `DocumentFunc` 的建構器，交付 `Each public page carries its own title and description`。頁面種類由呼叫端指定，因為歷史頁掛在 StripPrefix 之後，請求路徑一律是 `/`，只有查詢字串保留；日期仍取自查詢字串，canonical 仍以資料庫中的報告日期組出。`Lang` 留空以沿用 `syralit.toml`。驗證方式：1.1 的測試全數轉綠。
- [x] 1.3 交付 `Public pages carry canonical and social sharing metadata`：確認 `Document.HeadHTML` 仍包含 description、canonical、`og:type`、`og:site_name`、`og:title`、`og:description`、`og:url`、`og:locale` 與 `twitter:card`，且取自報告的文字經過 HTML 逸出。驗證方式：測試以含引號與角括號的頭條斷言輸出已逸出。
- [x] 1.4 交付 `Metadata failures degrade to site defaults`：報告查詢失敗或查無資料時，回傳的 `Document` 使用站台預設標題與描述。驗證方式：以會回傳錯誤的假 store 斷言 `Document.Title` 等於站台預設標題。

## 2. 移除中介軟體

- [x] 2.1 刪除 `internal/seo/middleware.go`，並把 `renderMetaTags` 移到 `internal/seo/meta.go` 供 `DocumentFunc` 使用。中介軟體、回應攔截器與 `rewriteHead` 不再存在，`Metadata rewriting never disturbs other responses` 這條需求也隨之不再適用，因為已經沒有任何程式碼會攔截或改寫回應。驗證方式：`grep` 確認專案中沒有殘留的 `Middleware`、`interceptor` 或 `rewriteHead` 參照，且 `go build ./...` 通過。
- [x] 2.2 在 `main.go` 為首頁與歷史頁各自建立帶 `DocumentFunc` 的 `Config`，後台維持不帶，並移除中介軟體的包裝。robots.txt 與 sitemap.xml 的掛載方式不變。驗證方式：`go test ./...` 全數通過，`main_test.go` 既有的後台無中繼資料斷言仍成立。

## 3. 收斂標題來源

- [x] 3.1 交付 `Browser title is not overwritten after connecting`：移除 `internal/site` 的 `pageConfig` 中設定固定站名的 `PageTitle`，讓連線後的分頁標題維持 `DocumentFunc` 寫入的值。驗證方式：瀏覽器實際開啟首頁，確認分頁標題是當日報告頭條而非站名。

## 4. 驗證與文件

- [x] 4.1 交付 `Public pages declare Traditional Chinese` 中的 `Language applies without per-page metadata` 情境：確認後台頁面仍帶 `lang="zh-Hant-TW"` 且沒有任何中繼資料標籤。驗證方式：以真實 mux 取得後台頁面並斷言兩者。
- [x] 4.2 更新 `AGENTS.md`：中繼資料改由 `DocumentFunc` 產生，刪除中介軟體相關的約定與壓縮中介軟體的警告，並移除已解決的 `SetPageConfig` 覆蓋標題那條 follow-up。驗證方式：閱讀該檔，確認沒有描述已不存在的機制。
- [x] 4.3 端對端確認無回歸：實際啟動伺服器，確認首頁、歷史列表、單日報告的 head 內容與移除前一致，後台無中繼資料，robots.txt 與 sitemap.xml 皆 200，且瀏覽器即時更新正常。驗證方式：保留各路徑輸出與瀏覽器截圖作為憑據。
