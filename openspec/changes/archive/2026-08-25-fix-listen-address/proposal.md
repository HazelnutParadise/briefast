## Problem

伺服器綁到隨機埠而不是設定檔宣告的位址。啟動 log 從 `Briefast listening on http://0.0.0.0:8600` 變成 `Briefast listening on http://:0`。

`:0` 是合法位址，`http.ListenAndServe` 會成功並綁一個任意空閒埠，過程不報任何錯。容器發布 8600 並以 `http://127.0.0.1:8600/` 做健康檢查，因此每次啟動都會健檢失敗，唯一可見的症狀是不斷重啟。

這同時違反 `deployment` 既有需求「Configuration injection」：該需求明訂映像必須綁定內建設定檔宣告的位址，而不是退回 loopback 預設；其情境「Image runs standalone」也要求服務綁定 syralit.toml 宣告的位址並能從容器外連入。

## Root Cause

`main.go` 透過 `sy.GetOption("server.host")` 與 `sy.GetOption("server.port")` 取得位址。syralit v0.10.1 把設定改為隨 handler 與 session 攜帶之後，`GetOption` 在 page function 以外呼叫時會使用一個零值 `Config`，且不套用內建預設，於是回傳空字串與 0。

`go test ./...` 全數通過也抓不到，因為所有測試都以 `httptest` 直接呼叫 `newHandler`，從不經過 `ListenAndServe`，位址計算完全沒有被任何測試覆蓋。

## Proposed Solution

改用 syralit v0.11.0 新增的 `sy.ResolveConfig`。它會把 syralit.toml 與內建預設套進傳入的 `Config` 並回傳結果，正是 `sy.Handler` 內部所做的解析，也是嵌入模式下取得設定的官方途徑。

在 `main.go` 收斂出單一設定來源，讓建立 handler 與計算監聽位址讀到同一份解析結果，並把位址計算抽成可測函式，讓測試能斷言其結果。

## Non-Goals

- 不用 `sy.Embed` 埋廣告，那是獨立工作。
- 不用 `Config.DocumentFunc` 移除 SEO 中介軟體，那是獨立工作。
- 不改設定的存放方式。host 與 port 維持在版控的 syralit.toml，不移到環境變數。
- 不改 `sy.GetOption` 在頁面內的既有用法，本專案並未在頁面內使用它。

## Success Criteria

- 以專案的 syralit.toml 啟動時，log 顯示 `Briefast listening on http://0.0.0.0:8600`。
- 存在一個測試斷言解析後的監聽位址等於 `0.0.0.0:8600`，且該測試在修正前會失敗。
- `go vet ./...` 與 `go test ./...` 全數通過。
- 公開頁、後台、robots.txt 與 sitemap.xml 的行為不變。

## Impact

- Affected specs: 修改 `deployment`
- Affected code:
  - Modified:
    - `main.go`
    - `main_test.go`
    - `go.mod`
    - `go.sum`
  - New: (none)
  - Removed: (none)
- 依賴：syralit v0.10.1 升為 v0.11.0。
- 對外行為：恢復為設定檔宣告的位址，其餘不變。
