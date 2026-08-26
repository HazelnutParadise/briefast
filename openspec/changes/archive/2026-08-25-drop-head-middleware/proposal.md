## Summary

改用 syralit 的 `Config.DocumentFunc` 產生每頁中繼資料，移除攔截並改寫回應的 SEO 中介軟體。

## Motivation

中介軟體是在框架沒有 per-request 中繼資料能力時的變通做法。它必須緩衝回應、重算 `Content-Length`、轉發 `http.Flusher` 與 `http.Hijacker`、辨識 `/_syralit/` 前綴以避開 WebSocket 與 SSE——這四件事都不是本專案的業務邏輯，而是在應用層重新實作框架的傳輸細節。任何一項寫錯，症狀都是即時更新停擺或頁面被截斷。

syralit v0.11.0 的 `Config.DocumentFunc` 在 shell 渲染前、session 建立前，依請求回傳 `Title` 與 `HeadHTML`，正是中介軟體在做的事，但由框架在正確的時機執行。改用它之後，上述四項顧慮全部消失。

同時處理一個既有落差：`internal/site` 呼叫 `sy.SetPageConfig(sy.PageTitle(...))` 設定固定站名，瀏覽器連線後會用它覆蓋 `document.title`。目前爬蟲拿到正確標題，但有執行 JavaScript 的訪客看到的是通用標題。移除中介軟體會讓標題來源收斂到 `DocumentFunc`，這個覆蓋必須一併移除，否則兩個來源會繼續打架。

## Proposed Solution

`internal/seo` 對外改為提供「依頁面種類產生 `DocumentFunc`」的建構器。首頁與歷史頁各自掛上自己的那一份，後台不掛，因此後台維持沒有中繼資料。

頁面種類必須由呼叫端指定，不能再從請求路徑判斷：歷史頁掛在 StripPrefix 之後，`DocumentFunc` 收到的 `r.URL.Path` 一律是 `/`，只有查詢字串保留下來。日期仍從查詢字串取得，canonical 仍以資料庫中的報告日期組出，不回寫請求帶進來的字串。

`renderMetaTags` 保留並改為產生 `Document.HeadHTML`，標題改由 `Document.Title` 帶出。`Lang` 留空以沿用 `syralit.toml` 的設定。robots.txt 與 sitemap.xml 兩個端點不受影響。

`internal/site` 的 `pageConfig` 移除 `PageTitle`，讓連線後的 `document.title` 維持 `DocumentFunc` 寫入的值。

## Non-Goals

- 不改公開頁的版面、內容或視覺。
- 不動 robots.txt 與 sitemap.xml 的輸出。
- 不做伺服器直出 HTML，爬蟲讀到的 body 仍是應用外殼。
- 不用 `sy.Embed` 埋廣告，那是另一個獨立工作。

## Alternatives Considered

保留中介軟體並只把固定標籤搬到 `Config.HeadHTML`。這會讓 head 同時有兩個來源，而中介軟體仍需為了動態標題而存在，四項傳輸顧慮一項都省不掉。

## Impact

- Affected specs: 修改 `crawler-metadata`
- Affected code:
  - Modified:
    - `main.go`
    - `main_test.go`
    - `internal/seo/meta.go`
    - `internal/seo/seo_test.go`
    - `internal/site/site.go`
    - `AGENTS.md`
  - Removed:
    - `internal/seo/middleware.go`
  - New: (none)
- 對外行為：公開頁的 head 內容與現況相同，來源改為框架。瀏覽器分頁標題改為該頁專屬標題，不再被站名覆蓋。後台與 API 不受影響。
