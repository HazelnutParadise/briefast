## Why

網站要開始承載廣告營收，需要一個文字驅動廣告版位。這類廣告的整合方式是「掛載點加上 loader script」，過去在本專案無法運作：版面全部由 `sy.HTML` 產生，而該元件以 innerHTML 插入內容，其中的 `<script>` 依規範不會執行，因此廣告的 div 會出現、loader 永遠不載入。

syralit v0.11.0 新增的 `sy.Embed` 會把內容插入主文件並執行其中的 script，且以 key 識別節點，內容不變時重用既有元素、不重跑 script。這讓版位得以成立。

## What Changes

- 報告頁在個股多空判斷之後新增一個廣告版位，位置是全頁唯一的閱讀斷點：上方是當日結論，下方是佐證細節。
- 版位以髮絲規線框出並加上「廣告」標示，讓第三方渲染的內容不與內文混淆。第三方內容的外觀不受本專案控制，框線與標示是把它在版面上界定為獨立區塊的方式。
- 廣告的掛載點寫在版面字串既有位置，`sy.Embed` 只承載 loader script。實測確認 `sy.Embed` 的 script 執行時，前一個 `sy.HTML` 節點的 DOM 已進入 document，因此 loader 取得得到掛載點，版面結構完全不需要拆開。
- 版位只出現在有報告內容的頁面。歷史列表沒有內文流，找不到報告的頁面也不適合承載廣告。

## Non-Goals

- 不在歷史列表頁、找不到報告頁與後台放置廣告。
- 不做伺服器直出，爬蟲讀到的 body 仍是應用外殼，廣告與此無關。
- 不調整廣告本身的外觀。第三方內容由廣告服務渲染，本專案只負責界定版位邊界。
- 不新增多個版位。

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `report-viewing`: 報告頁新增廣告版位，位置、標示與出現條件納入規格。

## Impact

- Affected specs: 修改 `report-viewing`
- Affected code:
  - Modified:
    - `internal/site/site.go`
    - `internal/site/styles.go`
    - `internal/site/site_test.go`
    - `DESIGN.md`
  - New: (none)
  - Removed: (none)
- 對外行為：報告頁在個股多空判斷與產業新聞摘要之間多出一個標示為廣告的區塊，並載入第三方廣告服務的 script。其他頁面不變。
- 第三方相依：頁面會向廣告服務的網域請求 loader script。
