## Context

使用者確認以報紙插圖減少文字密度，並允許低彩度彩色。沿用暖紙底、襯線標題、規線及紅漲綠跌。

## Goals / Non-Goals

納入：頭版圖像導覽、產業插圖、總覽與事件展開閱讀、深淺色與行動版。
排除：API JSON、資料庫、新聞內容改寫、自動產生每日圖片、權限與正式環境變更。插圖代表產業，不宣稱特定公司、事件或市場走勢。

## Decisions

用一張透明四格插圖素材供四類產業重用，減少請求數與維護成本。圖片打包進 Go binary，公開靜態端點僅支援 GET/HEAD 與快取。數值仍由 Go 依報告計算，圖片不含數字或文案。用原生 details/summary 保留完整內文與鍵盤操作，不新增 JS 或依賴。

## Implementation Contract

renderEditionGuide(industries []report.Industry) string：略過零事件類別，總則數為各 len(Events) 加總，比例為該類別則數除總則數。HTML nav 的 aria-label 為「本期產業新聞導覽」，標題「本期新聞分布」，各 a 的 href 為 #industry-N，N 為原始 slice 索引，包含類別、N 則與 aria-hidden 的比例條。空輸入回空字串。未知名稱以 html.EscapeString 輸出，不對它挑一張不相干插圖。

renderIndustryIllustration(name string) string：科技、金融、傳產、房地產各對應一個固定 sprite class sector-tech、sector-finance、sector-industry、sector-property，輸出具 aria-hidden 的 sector-art div。未知名稱回空字串。不得串接未轉義輸入到 class 或 style。兩個函式只放 internal/site/editorial.go，不修改測試。

整合格子：頭條後加導覽，再用初始收合的 overview-details 包住完整 overview_md，再接現有 market outlook。產業 article id 與 guide href 同步，各產業保留名稱、圖與全部事件標題；每個正文在 native details 內。原文與 watch、來源、個股、法說及免責聲明不變。

資產格子：internal/site/assets/editorial-sectors-v1.webp 是四格透明點陣圖。internal/site/editorial_assets.go 提供 EditorialIllustrationHandler() http.Handler，GET /assets/editorial-sectors-v1.webp 傳 WebP、ETag、Cache-Control 及 nosniff，支援 HEAD/304。main.go 只增加此 GET 靜態掛載。

驗收：先跑新增導覽/收合/靜態端點測試確認失敗，完成後 go test ./... 與 git diff --check 通過。實際查看首頁和歷史頁，確認圖與文字、比例、錨點與展開操作，並測桌面、390px 行動版、深色與列印。DESIGN.md 允許插圖用低彩度灰藍 #6F7E86、赭黃 #AE956B、灰褐 #94877B，UI 紅綠仍僅多空。圖像失敗時仍能讀文字及則數，字體與重要資訊對比維持原樣。

## Risks / Trade-offs

圖像僅表示產業，不呈現資料衍生的供應鏈因果。分布圖只表示報導數，不當成產業強弱。收合減少首屏文字，使用者可展開原文；列印全部展開。深色以局部濾鏡適配插圖，須實際檢視。版本化資產有快取，後續改圖需換檔名。
