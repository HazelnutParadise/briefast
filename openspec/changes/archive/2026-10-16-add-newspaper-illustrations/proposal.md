## Why

晨報以長段落連續呈現，讀者難以快速找到關注的產業。加入淡彩報紙插圖與新聞分布導覽，並讓完整內文按需展開。

## What Changes

- 頭版增加插圖式產業導覽，依本期產業新聞則數顯示比例，連到各產業。
- 科技、金融、傳產、房地產各有墨線、網點與低彩度灰藍／赭黃插圖。
- 盤前總覽及產業事件正文改成可展開閱讀，保留原文與來源。
- 插圖內嵌於程式並由唯讀靜態網址提供，不新增外部圖片依賴。

## Capabilities

### Modified Capabilities

- `report-viewing`: 報紙插圖、產業新聞分布與漸進閱讀。

## Impact

- Modified: `internal/site/site.go`, `internal/site/styles.go`, `main.go`, `main_test.go`, `DESIGN.md`, `README.md`
- New: `internal/site/editorial.go`, `internal/site/editorial_test.go`, `internal/site/editorial_assets.go`, `internal/site/assets/editorial-sectors-v1.webp`, `docs/design/editorial-illustrations.md`
- 無報告 JSON、資料庫、權限或預測規則變更。不部署或發布正式報告。
