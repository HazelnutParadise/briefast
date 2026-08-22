## 1. 站台網址解析

- [x] 1.1 先寫 `internal/seo/seo_test.go` 的失敗測試，涵蓋「站台絕對網址三層解析」的三種情況：設了 `BRIEFAST_SITE_URL` 時以它為準、未設但有 `X-Forwarded-Proto` 時用該通訊協定加請求主機、兩者都沒有時依連線是否為 TLS 決定 `https` 或 `http`。驗證方式：`go test ./internal/seo/` 因函式未實作而失敗。
- [x] 1.2 實作站台網址解析函式，讓 `Site base URL resolution` 的三個情境全數通過。驗證方式：`go test ./internal/seo/` 三個情境測試轉綠。

## 2. 頁面中繼資料組裝

- [x] 2.1 先寫失敗測試，涵蓋「各頁中繼資料由請求路徑與查詢參數決定」：首頁取最新報告頭條、`/history/` 取歷史列表文案、`/history/?date=` 取該日報告頭條；以及描述取自總覽段落並去除 Markdown 標記、超過 150 字元時截斷補刪節號。驗證方式：`go test ./internal/seo/` 失敗。
- [x] 2.2 實作中繼資料組裝，交付 `Each public page carries its own title and description` 的行為：每一頁產出專屬的標題與描述。驗證方式：2.1 的測試全數轉綠。
- [x] 2.3 交付 `Metadata failures degrade to site defaults` 的行為：報告查詢回錯或查無資料時改用站台預設標題與描述，回應狀態碼不變。驗證方式：以會回傳錯誤的假 store 執行測試，確認頁面仍以 200 送出且帶預設文案。

## 3. head 改寫中介軟體

- [x] 3.1 先寫失敗測試，涵蓋「用回應改寫中介軟體，不改 Syralit 也不做 UA 分流」的替換行為：改寫後的 HTML 帶 `lang="zh-Hant-TW"`、專屬 `title`、`meta description`、`link rel="canonical"`、`og:type`、`og:site_name`、`og:title`、`og:description`、`og:url`、`og:locale` 與 `twitter:card`，且 body 內容一字未改。驗證方式：`go test ./internal/seo/` 失敗。
- [x] 3.2 實作中介軟體的字串替換，交付 `Public pages declare Traditional Chinese` 與 `Public pages carry canonical and social sharing metadata` 兩項行為。驗證方式：3.1 的測試轉綠。
- [x] 3.3 先寫失敗測試，涵蓋「依 Content-Type 決定是否緩衝，WebSocket 與 SSE 一律直通」與 `Metadata rewriting never disturbs other responses`：路徑以 `/_syralit/` 開頭時完全不包裝、Content-Type 非 `text/html` 或狀態非 200 時原樣輸出、HTML 內找不到 `title` 時原樣輸出、改寫後 `Content-Length` 等於新 body 的位元組長度、包裝器有轉發 `http.Flusher` 與 `http.Hijacker`。驗證方式：`go test ./internal/seo/` 失敗。
- [x] 3.4 實作回應攔截器的直通判斷與介面轉發，讓 3.3 的所有情境轉綠。驗證方式：`go test ./internal/seo/` 全綠。

## 4. robots.txt 與 sitemap.xml

- [x] 4.1 交付 `Site serves robots.txt` 的行為：`GET /robots.txt` 回 200 與 `text/plain`，內容對 `/admin/` 與 `/api/` 下 Disallow，並以絕對網址宣告 sitemap 位置。驗證方式：先寫 `httptest` 測試比對狀態碼、Content-Type 與三行內容，再實作至轉綠。
- [x] 4.2 交付 `Site serves a sitemap covering every report` 的行為，依「sitemap 直接從 store 全量產出」用 `CountReports` 取總數後以 `ListReports` 一次撈完，輸出含首頁、歷史列表頁與每份報告的 `urlset`，每筆帶 `lastmod`；零報告時仍輸出合法 `urlset`，查詢失敗時回 500 且不輸出半份 XML。驗證方式：先寫涵蓋有報告、零報告、查詢失敗三種情況的 `httptest` 測試，再實作至轉綠。

## 5. 掛載與設定

- [x] 5.1 在 `main.go` 把首頁與歷史頁的 handler 包上 head 改寫中介軟體，並掛上 `/robots.txt` 與 `/sitemap.xml` 兩條路由，`/admin/` 與 `/api/` 的掛載方式維持不變。驗證方式：`go test ./...` 全數通過，且既有的 `main_test.go` 沒有出現新的失敗。
- [x] 5.2 把新的部署設定鍵 `BRIEFAST_SITE_URL` 寫進 `.env.example` 與 README 的設定說明，講明未設定時會依請求標頭推導、正式部署必須設定。驗證方式：閱讀兩份檔案，確認鍵名、用途與預設行為三項資訊齊全。
- [x] 5.3 在 `AGENTS.md` 的開發約定補一條：公開頁的 head 中繼資料由 SEO 中介軟體產生，日後若加入壓縮之類的中介軟體必須包在改寫之外。驗證方式：閱讀 `AGENTS.md`，確認該條約定存在且與現有條目風格一致。

## 6. 端對端驗證

- [x] 6.1 實際啟動伺服器後確認公開頁行為：`curl` 首頁可看到 `zh-Hant-TW` 與完整 meta 標籤、`/robots.txt` 與 `/sitemap.xml` 皆回 200、瀏覽器開首頁仍能正常顯示報告內容且即時更新未受影響。驗證方式：保留 `curl` 輸出與瀏覽器截圖作為憑據。
