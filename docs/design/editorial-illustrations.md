# 報紙淡彩插圖

四格素材用於科技、金融、傳產與房地產的產業導覽，並作為產業標題的圖像識別。圖像只表示主題，不表示行情、數據或特定新聞事件。新聞則數與比例由報告資料計算。

## 素材

- 正式素材：internal/site/assets/editorial-sectors-v1.webp
- 內建 ImageGen 生成，透明背景，1254 × 1254。
- 上左科技、上右金融、下左傳產、下右房地產。
- 淡灰藍、赭黃、灰褐配墨線與網點，紅綠保留給多空。
- 原始 PNG 2,459,311 bytes，WebP 361,878 bytes，保留透明度。
- WebP 編碼：cwebp -q 80 -alpha_q 90 -m 6。
- 圖片嵌入 Go binary，GET /assets/editorial-sectors-v1.webp 可快取。更換圖片時使用新的版本檔名。

## 生成提示

Create a production asset for a serious Taiwanese financial newspaper website. ONE square 2-by-2 SPRITE ATLAS containing exactly FOUR separate editorial engraved illustrations on a truly transparent background. All four equal quadrants must have equal subject scale and generous inset margins, absolutely no overlap or crossing quadrant boundaries, no separator lines, no text or symbols. Top-left quadrant: a silicon wafer partially behind a single detailed integrated circuit package, representing technology. Top-right quadrant: a classical stone bank facade with a small stack of plain unbranded coins, representing finance. Bottom-left quadrant: a factory building with chimneys beside port shipping containers and a small crane, representing traditional industry and trade; no smoke clouds. Bottom-right quadrant: a small cluster of residential buildings under construction with one tower crane, representing real estate. Style: sophisticated newspaper copperplate etching and fine pen crosshatching, small halftone / stipple shadow texture like ink absorbed into newsprint, realistic architectural proportions, editorial object vignettes. Apply sparse very low-saturation translucent ink washes, muted slate blue #6F7E86, weathered ochre #AE956B and warm gray brown #94877B; main etched outlines charcoal #33302E. Colors should feel slightly faded as printed on a newspaper, not bright, not glossy, not cartoon, not flat vector icons. Avoid red and green entirely because these are reserved for financial up/down semantics in the consuming website. All four vignettes must be equally detailed and visually balanced, with a similar baseline and comfortable clear space around each object. No background scenery, no paper sheet rectangle, no white mat, no borders, no shadows from floating cards, no logos, no company names, no flags, no people, no labels, no numbers, no price graph, no arrows. The background must be transparent alpha, including clear gaps around every vignette. This is a reusable four-sector editorial illustration sprite, not an infographic with facts.

## 邊緣修整提示

Clean up this exact 2x2 newspaper illustration sprite atlas. Keep all four subjects, engraved detail, low saturation blue ochre gray palette, object positions and equal quadrant layout unchanged. Remove EVERY bright red, fluorescent yellow, orange fringe, halo, and chroma-key artifact around object silhouettes. Silhouette edges should be perfectly clean, neutral charcoal with smooth antialiased alpha against a truly transparent background. Absolutely no bright red, yellow, orange, or green pixels anywhere. Preserve the muted pale ochre washes inside coins and crane, but remove all high-saturation edge pixels. Do not add background, text, labels, divider lines, logos, or new subjects. Match all composition, inset spacing and scale exactly. This is a production-ready transparent PNG cleanup, not a redesign.

## 驗證

- 首頁與歷史報告頁：桌面及 390px 手機版均可顯示插圖與實際則數比例，手機無橫向溢出。
- 原生展開閱讀：點選及 Enter 可切換完整內容，產業連結可跳到對應錨點。
- 深色插圖可讀，列印時收合的全文仍顯示。
- go test ./...、go vet ./...、git diff --check 通過。
- Lighthouse 本機預覽的慢速手機測試：accessibility 100、performance 55、CLS 0.002、TBT 0 ms。這是測試內容及本機伺服器的結果，不是正式站測量。
- Lighthouse 指出現有 Google Fonts CSS 為阻擋顯示的資源，列為後續優化。本次不變更字型載入方式。
