## Why

公開頁目前對搜尋引擎與社群分享爬蟲完全不可讀。實測 `GET /` 只回傳 1,065 bytes 的 Syralit 空殼，內文是「連線中…」，報告內容全靠 WebSocket 推送；Googlebot 的算圖服務不支援 WebSocket，LINE 與 Facebook 的分享爬蟲更不執行 JavaScript。同時 robots.txt 與 sitemap.xml 皆回 404，head 內沒有 description、canonical 與 Open Graph，`lang` 標成 `en`。

先補這一批不動架構的基礎，能立刻換到正確的社群分享預覽卡與可被發現的網址清單，也為第二批的伺服器直出留好接口。

## What Changes

- 新增 `GET /robots.txt`：允許公開頁，`Disallow` `/admin/` 與 `/api/`，並宣告 sitemap 絕對網址。
- 新增 `GET /sitemap.xml`：從 store 撈出所有報告日期動態產生，含首頁、歷史列表頁與每份報告的檢視網址，附 `lastmod`。
- 新增 HTML head 改寫中介軟體，套用於公開頁：把 `lang` 改成 `zh-Hant-TW`，依頁面寫入 `title`、`meta description`、`canonical`、Open Graph 與 Twitter card。
- 新增站台絕對網址解析：優先讀環境變數 `BRIEFAST_SITE_URL`，未設定時由請求的 `Host` 與 `X-Forwarded-Proto` 推導。

## Non-Goals

- 不做公開頁伺服器直出 HTML。爬蟲仍讀不到報告內文，那是第二批的工作，本批只處理 head 與可被發現性。
- 不改網址結構。`/history/?date=` 維持現狀，改成路徑式網址與 301 導向屬於第二批。
- 不加 JSON-LD 結構化資料。它需要內文一起直出才有意義，留給第二批。
- 不用 User-Agent 分流對爬蟲送不同內容。人與爬蟲必須拿到同一份回應。
- 不修改 Syralit 套件本身。上游加 head 注入介面是另一個 repo 的事。

## Capabilities

### New Capabilities

- `crawler-metadata`: 公開頁對搜尋引擎與社群爬蟲提供的 head 中繼資料、robots.txt 與 sitemap.xml。

### Modified Capabilities

(none)

## Impact

- Affected specs: 新增 `crawler-metadata`
- Affected code:
  - New:
    - `internal/seo/seo.go`
    - `internal/seo/seo_test.go`
  - Modified:
    - `main.go`
    - `main_test.go`
    - `.env.example`
    - `docker-compose.yml`
    - `README.md`
    - `AGENTS.md`
  - Removed: (none)
- 對外行為：公開頁回應的 head 內容改變，body 與版面不變；`/admin/` 與 `/api/` 回應完全不受影響。
- 新增部署設定鍵 `BRIEFAST_SITE_URL`，未設定時可正常運作。
