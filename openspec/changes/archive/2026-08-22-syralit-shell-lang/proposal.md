## Summary

升級 syralit 到 v0.9.0，把文件語言標示交還給框架，中介軟體只留每頁不同的中繼資料。

## Motivation

`crawler-metadata` 上線時，syralit v0.7.0 把 `lang="en"` 寫死在 app shell 裡，沒有任何設定能改，所以語言標示只能靠 SEO 中介軟體攔下回應做字串替換。

v0.9.0 補上了 `Config.Lang`（連同 `Config.Dir` 與 `Config.HeadHTML`），語言標示終於能在框架層設定。交還給框架有兩個實質好處：

- 中介軟體找不到可替換的 `<title>` 時會整份原樣放行，目前那條降級路徑會連語言標示一起放掉。交給框架之後，即使中繼資料改寫失效，語言標示仍然正確。
- `syralit.toml` 的設定由三個 Syralit app 共用，後台頁面目前是 `lang="en"`，會一併修好。後台介面是繁體中文，卻對瀏覽器與輔助技術宣告為英文。

## Proposed Solution

在 `syralit.toml` 加 `lang = "zh-Hant-TW"`。`applyToConfig` 只在 Go 端的 `Config.Lang` 為空時才套用檔案設定，而程式建立的 `sy.Config` 沒有指定 Lang，因此檔案設定會生效。

接著把中介軟體裡替換語言標示的那一段移除，改寫流程只剩替換 title 元素。對應的測試改成驗證中介軟體不再碰語言標示，端對端測試仍然驗證回應帶有正確語言標示，只是來源換成框架。

## Non-Goals

- 不改用 `Config.HeadHTML`。它的說明明訂設定值套用於每一個請求、沒有 per-request 變體，而本專案的 title、description、canonical 每頁都不同。把固定標籤搬到設定檔、動態標籤留在中介軟體，等於把同一段 head 拆成兩個來源，維護成本高於收益。
- 不使用 `Config.Dir`。繁體中文是由左至右，維持屬性不輸出即可。
- 不採用 v0.9.0 的 Page URLs。本專案沒有使用 `sy.AddPage`，三個 app 各自掛在不同路徑。
- 不動公開頁伺服器直出，那是獨立的後續工作。

## Alternatives Considered

在 `main.go` 的 `sy.Config` 設定 `Lang` 欄位。效果相同，但 `syralit.toml` 是專案既有的非機密設定檔，語言屬於設定而非程式邏輯，放檔案裡與 title、theme、i18n 一致。

## Impact

- Affected specs: 修改 `crawler-metadata`
- Affected code:
  - Modified:
    - `syralit.toml`
    - `go.mod`
    - `go.sum`
    - `internal/seo/middleware.go`
    - `internal/seo/seo_test.go`
    - `AGENTS.md`
  - New: (none)
  - Removed: (none)
- 對外行為：後台頁面的語言標示由 `en` 改為 `zh-Hant-TW`；公開頁的語言標示不變，來源改為框架。版面、內容與其他中繼資料皆不變。
- 依賴：syralit v0.7.0 升為 v0.9.0。連帶更新 go-isatty 與 modernc.org/libc 兩個間接依賴，並將 go 指令版本改為 1.25.12。
