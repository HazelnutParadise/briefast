## Context

Briefast 的公開頁由 Syralit v0.7.0 渲染。Syralit 先回傳一張固定的空殼 HTML，再由瀏覽器開 WebSocket（失敗時退回 SSE）把畫面推下來。實測 `GET /` 的回應只有 1,065 bytes，head 內僅有 charset、viewport、`<title>Briefast</title>` 與樣式表連結，body 內文是「連線中…」。

Syralit v0.7.0 沒有任何 head 注入或預渲染介面，`sy.Config` 只有 Title、Host、Port、Theme、MaxUploadSizeMB、SSL 與 UIStrings 欄位，`sy.SetPageConfig` 設的標題是連線後由前端改寫 `document.title`，爬蟲拿不到。因此本批只能從 Go 這一側，在回應送出前改寫 head。

版面字串已經全部在伺服器端組好（`internal/site` 把整頁組成 HTML 字串再交給 Syralit），第二批要改成直出時可以沿用；本批刻意不碰那條路徑。

## Goals / Non-Goals

**Goals:**

- 公開頁的每一頁都有正確的語言標示、專屬標題、描述、canonical 與社群分享中繼資料。
- 站台提供 robots.txt 與涵蓋所有報告的 sitemap.xml。
- 站台絕對網址不寫死，可由部署環境設定，未設定時仍能運作。
- 對 `/admin/`、`/api/` 與 Syralit 的即時通道零影響。

**Non-Goals:**

- 不做伺服器直出報告內文，爬蟲讀到的 body 仍是空殼。
- 不改網址結構、不加 301 導向。
- 不加 JSON-LD。
- 不用 User-Agent 判斷送出不同內容。
- 不對 `/admin/` 加 `noindex` meta，改由 robots.txt 的 Disallow 處理，避免同時使用兩種互相牴觸的機制。

## Decisions

### 用回應改寫中介軟體，不改 Syralit 也不做 UA 分流

改寫既有回應是三個選項裡唯一同時滿足「不動架構」與「人與爬蟲拿到同一份回應」的做法。往上游改 Syralit 屬於另一個 repo，時程不可控；UA 分流要維護兩份輸出，且 Google 官方只把它當過渡手段。

中介軟體只做字串替換，不解析 HTML：把 head 裡第一個 `lang="en"` 換成 `lang="zh-Hant-TW"`，再把 `<title>` 元素整段換成新的標題加上一串 meta 標籤。找不到 `<title>` 就原樣放行，讓 Syralit 未來改版時退化成不改寫而不是輸出壞掉的 HTML。

### 依 Content-Type 決定是否緩衝，WebSocket 與 SSE 一律直通

`/_syralit/ws` 與 `/_syralit/sse` 掛在同一個 handler 底下，緩衝它們會讓即時更新整個停擺。攔截器採兩道保險：

1. 請求路徑以 `/_syralit/` 開頭時完全不包裝，直接交給下游。
2. 包裝後的 ResponseWriter 在寫出標頭時才判斷，Content-Type 不是 `text/html` 或狀態碼不是 200 就切換成直通模式，後續寫入不再經過緩衝。

包裝器必須同時實作 `http.Flusher` 與 `http.Hijacker` 並轉發給底層，否則 SSE 無法即時送出、WebSocket 無法完成升級。改寫完成後必須重設 `Content-Length`，長度對不上瀏覽器會截斷頁面。

### 各頁中繼資料由請求路徑與查詢參數決定

中介軟體持有 store，依請求決定要查哪一份報告：

| 請求 | 標題 | 描述來源 |
| --- | --- | --- |
| `/` | 最新報告頭條 | 最新報告的總覽段落 |
| `/history/` | 歷史報告 | 固定說明句 |
| `/history/?date=YYYY-MM-DD` | 該日報告頭條 | 該日報告的總覽段落 |
| 查無報告或查詢失敗 | 站名與副標 | 固定說明句 |

描述取自 `OverviewMD`，先去掉 Markdown 標記與換行，再截到 150 個字元以內，截斷處補刪節號。查詢失敗一律降級成固定文案，不讓中繼資料的問題變成頁面錯誤。

### 站台絕對網址三層解析

canonical、`og:url` 與 sitemap 需要絕對網址。解析順序為：環境變數 `BRIEFAST_SITE_URL`；否則用請求的 `X-Forwarded-Proto` 加 `Host`；再否則依連線是否為 TLS 決定通訊協定，配 `Host`。正式部署一律設環境變數，README 與 `.env.example` 要寫清楚，反向代理沒送 `X-Forwarded-Proto` 時會退化成 `http`。

### sitemap 直接從 store 全量產出

報告一天一份，數量以年計仍遠低於 sitemap 的五萬筆上限，不需要分頁或索引檔。每份報告的 `lastmod` 用該報告的產生時間，首頁的 `lastmod` 用最新報告的產生時間。網址沿用現行的 `/history/?date=YYYY-MM-DD` 格式，第二批改網址時一併更新。

## Implementation Contract

**行為**

- `GET /robots.txt` 回傳 200、`text/plain`，內容允許所有爬蟲，對 `/admin/` 與 `/api/` 下 Disallow，並以絕對網址宣告 sitemap 位置。
- `GET /sitemap.xml` 回傳 200、`application/xml`，是合法的 urlset，包含首頁、歷史列表頁與每一份報告的網址，每筆帶 `lastmod`。沒有任何報告時仍回傳合法的 urlset，只含首頁與歷史列表頁。
- 公開頁的 HTML 回應，head 內含 `lang="zh-Hant-TW"`、該頁專屬的 `title`、`meta description`、`link rel="canonical"`、`og:type`、`og:site_name`、`og:title`、`og:description`、`og:url`、`og:locale`，以及 `twitter:card` 為 `summary`。
- 同一次請求的 body 內容與版面完全不變，只有 head 改變。

**介面**

- 新套件 `internal/seo` 對外提供三樣東西：robots.txt 的 handler、sitemap.xml 的 handler，以及一個把 `http.Handler` 包成會改寫 head 的中介軟體。三者都接受 `*store.Store` 與站台網址解析設定。
- `main.go` 只在掛載首頁與歷史頁時套用中介軟體，`/admin/` 與 `/api/` 的掛載方式不變。

**失敗模式**

- store 查詢失敗或查無報告：改用站台預設標題與描述，回應照常送出，不回傳錯誤。
- 回應裡找不到可替換的 `<title>`：原樣放行。
- sitemap 查詢失敗：回傳 500 與純文字錯誤訊息，不輸出半份 XML。

**驗收條件**

- `internal/seo` 的測試用 `httptest` 驗證：robots.txt 內容、sitemap 在有報告與零報告兩種情況下的輸出、head 改寫後含所有必要標籤且 `Content-Length` 與 body 長度一致、非 `text/html` 的回應原封不動、路徑以 `/_syralit/` 開頭的請求不被包裝、缺少 `<title>` 時原樣放行、站台網址三層解析各自的結果。
- `go test ./...` 全數通過。
- 實際起 server 後 `curl -s http://127.0.0.1:8600/` 可看到 `zh-Hant-TW` 與完整 meta 標籤，`curl` robots.txt 與 sitemap.xml 皆回 200，且首頁在瀏覽器中仍能正常顯示報告內容。

**範圍邊界**

- 在範圍內：`internal/seo` 新套件、`main.go` 的掛載、`.env.example` 與 README 的新設定鍵說明、AGENTS.md 的約定補充。
- 在範圍外：`internal/site` 的版面程式碼、`internal/admin`、`internal/api`、`internal/store` 的既有函式行為，以及 Syralit 套件本身。若需要 store 新增查詢函式，只新增不修改既有函式。

## Risks / Trade-offs

- 緩衝回應可能拖垮即時通道 → 路徑前綴與 Content-Type 兩道判斷都要有，包裝器必須轉發 Flush 與 Hijack，測試要涵蓋這三種情況。
- 改寫後忘了重設 Content-Length 會讓頁面被截斷 → 測試直接比對標頭與 body 長度。
- 每個公開頁請求多一次資料庫查詢 → 頁面流量低且 SQLite 為本機檔案，先不加快取；若日後量大再考慮短 TTL 快取。
- 反向代理未送 `X-Forwarded-Proto` 時 canonical 會變成 `http` → 正式部署要求設定 `BRIEFAST_SITE_URL`，文件中明列。
- 日後若加入壓縮中介軟體，包在改寫之外才正確 → 記在 AGENTS.md 的約定。
- 爬蟲此時仍讀不到報告內文，搜尋收錄不會因為這批而改善，只有社群分享預覽與網址發現性會改善 → 這是刻意的分批，不是缺漏。
