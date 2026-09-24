## Context

`market_outlook.direction` 是相對前一交易日收盤的方向，`trajectory_md` 是盤前預期的開盤、盤中、尾盤敘述。站台由 Go 渲染固定版面，資料庫保存完整報告 JSON。盤前沒有當日實際分時數據，不能把示意圖做成點位預測。

## Goals / Non-Goals

**Goals:** 用一眼可讀的三階段圖呈現有依據的基準情境，同時保留文字中的前提和反向條件，讓舊報告仍可讀。

**Non-Goals:** 指數點位、漲跌幅、發生機率、精確轉折時刻、即時行情、互動縮放、從自由文字自動猜圖。

## Decisions

### Add optional qualitative trajectory chart data

在 `market_outlook` 加選填 `trajectory_chart`，含 `open`、`midday`、`close` 三個欄位，各只能是 `above`、`near`、`below`，表示相對前一交易日收盤的定性位置。圖存在時需有非空白 `trajectory_md`；`up` 的 `close` 須是 `above`、`down` 須是 `below`、`range` 須是 `near`，`uncertain` 不提供圖。用明確資料避免從文字猜測及方向標籤自動畫圖。API 缺欄位或 null 視為舊資料而不畫圖。

### Render one static three-stage SVG

站台在大盤預測內由 Go 直接產生一個 SVG：橫向開盤、盤中、尾盤，縱向只有前收上方、附近、下方三個定性層級。折線與三個點呈現基準走法，末點沿用台股紅漲綠跌；各階段直接顯示階段與層級，並標示「情境示意・非點位」。不加 JavaScript 或圖表依賴。既有文字置於圖後，供條件與細節閱讀；SVG 有完整文字替代描述。首頁與歷史頁使用同一個 render path。

### Preserve uncertainty and legacy reports

新每日流程只有在證據足以描述三階段時才填 `trajectory_chart`，且與 `trajectory_md` 一致；若文字明說盤中路徑無法可靠判斷，就省略圖。舊報告無圖時保留方向、走法文字和判斷依據，沒有預測時整段省略。

## Implementation Contract

- `market_outlook.trajectory_chart` 選填，其 JSON 形狀為 `{ "open": "above|near|below", "midday": "above|near|below", "close": "above|near|below" }`。null 視為缺席。任何非上述值、缺少階段、與 `direction` 收盤方向不符、或有圖但缺少有效 `trajectory_md`，`POST /api/report` 回 400，指出 `market_outlook.trajectory_chart` 且不寫入。有效圖隨報告按日期 GET 原樣讀回。
- 有圖時首頁與歷史報告的大盤預測在方向後、文字走法前顯示同一張定性三階段圖；折線、點、階段與層級標籤、示意非點位聲明可見，SVG 提供非視覺描述。無圖的報告沒有空白圖或佔位內容。圖表不依賴即時資料、外部 JavaScript、hover 或點擊。
- 新每日流程在可推斷基準路徑時填圖及走法文字；證據不足時仍寫出無法判斷與待觀察條件，省略圖。圖的 `close` 與方向一致，三階段與文字互相一致。
- 驗收涵蓋有效、缺欄位、錯值、方向矛盾、舊報告 API 測試；首頁與歷史頁渲染、舊報告缺圖測試；每日範例有效性；`go test ./...`、Spectra 驗證、桌面與手機窄版實際畫面及無橫向溢出檢查。

## Risks / Trade-offs

- [圖看似精確走勢] → 只用三個定性層級，不顯示點位或刻度，顯著標示情境示意與非點位。
- [圖與文字或方向矛盾] → API 檢查收盤方向，每日流程要求逐階段核對敘述；有疑慮就省略圖。
- [手機空間不足] → 單張 SVG 依容器等比例縮放，階段與層級文字使用 HTML 排列，不把小字塞進 SVG。
