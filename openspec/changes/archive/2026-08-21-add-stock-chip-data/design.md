## Context

個股判讀目前的資料脈絡只有新聞與昨收價。昨收價已有一套完整模式可沿用：蒐集期抓全市場資料存檔、日期驗證（民國格式換算）、jq/grep 按代號查值不整份進 context、失敗分流照常發布。籌碼面資料（三大法人買賣超、融資融券）本質相同：全市場、每日一包、帶資料日期，適合複製同一套模式。

四個端點已於 2026-08-21 實測可用：

- 上市三大法人買賣超日報：www.twse.com.tw 的 rwd/zh/fund/T86（帶 date=YYYYMMDD 與 selectType=ALL 參數，回應含 date 與逐檔股數欄位，單位：股）
- 上市融資融券：www.twse.com.tw 的 rwd/zh/marginTrading/MI_MARGN（同參數形式，tables 第二個表為逐檔融資融券前日／今日餘額，單位：交易單位＝張）
- 上櫃三大法人買賣明細：www.tpex.org.tw 的 openapi/v1/tpex_3insti_daily_trading（回應逐檔含民國格式 Date 與各法人 Difference 欄位，單位：股）
- 上櫃融資融券餘額：www.tpex.org.tw 的 openapi/v1/tpex_mainboard_margin_balance（逐檔含前日／今日融資融券餘額，單位：張）

報告以整包 JSON 存於 reports 表的 payload 欄位，schema 擴充不需資料庫 migration。

## Goals / Non-Goals

**Goals:**

- 晨報帶有報告內個股的前一交易日籌碼面數字：外資、投信、自營商買賣超，三大法人合計，融資融券餘額增減。
- 個股新聞詳情條目以水平長條視覺化三類法人買賣超，紅買超綠賣超，遵循台股紅漲綠跌與 FT 系規線風格（零圓角、零陰影、伺服器端渲染、不引入前端圖表庫）。
- 判讀時把籌碼數據納為佐證，引用紀律與昨收價一致：標明資料日期、不得單獨構成 call、缺料整塊省略。
- 舊報告與缺 chips 的新報告完全相容，渲染不出錯、無佔位內容。

**Non-Goals:**

- 不做多日籌碼趨勢圖。單日報告只嵌單日數據；趨勢需要跨日累積資料，屬後續變更。
- 不在個股多空判斷（calls）四卡與首頁其他區塊顯示籌碼數據，只做個股新聞詳情。
- 不抓借券賣出、鉅額交易、股權分散表、主力分點等其他籌碼資料。
- 不新增任何資料表或 migration。

## Decisions

### 沿用昨收價參考資料模式抓取籌碼資料

四個端點在蒐集階段各以一個一般 HTTP 請求整包拉回存檔，與昨收價資料同層。理由：模式已被 spec 驗證過（讀取紀律、日期驗證、失敗分流都有現成條文可對齊），執行 agent 已熟悉同構流程。替代方案是判讀時按個股逐檔查詢證交所個股頁，被否決：請求數量隨個股數線性成長，且觸發 TWSE 鎖 IP 風險高。

兩個 www.twse.com.tw 端點併入既有 TWSE 節流序列（序列執行、間隔至少 5 秒、429 等 60 秒重試一次）。www.twse.com.tw 與 openapi.twse.com.tw 主機名不同但同屬證交所，保守起見一律視為同一節流對象。櫃買兩個端點不受該節流限制。

日期驗證與昨收價一致：T86 與 MI_MARGN 用目標日期（台北時區前一交易日）的 YYYYMMDD 作參數並驗證回應 date；櫃買回應內嵌民國格式 Date，換算後必須等於前一交易日，過舊視為抓取失敗，不得沿用舊數據。

### chips 資料形狀與驗證規則

stock_news 條目新增選填 chips 物件，欄位固定：

- date：資料日期 YYYY-MM-DD（必填於 chips 存在時）
- foreign_net、trust_net、dealer_net、total_net：外資（含外資自營商）、投信、自營商、三大法人合計買賣超，單位：股，整數，可負
- margin_change、short_change：融資、融券餘額增減（今日餘額減前日餘額），單位：張，整數，可負

驗證規則：chips 缺漏或為 null 一律接受（舊報告與抓取失敗場景）；chips 存在時 date 必須是有效 YYYY-MM-DD，違規訊息標明條目索引。不驗證 total_net 等於三項加總——來源另拆外資自營商子欄位，強制等式會因來源口徑誤拒合法資料；total_net 直接取來源的三大法人合計欄位。餘額增減由 skill 從前日／今日餘額相減求得，不另存餘額絕對值，版面要的是流向不是存量。

替代方案「另立頂層 chips 區塊集中放全部個股」被否決：chips 跟著 stock_news 條目走，渲染與驗證都不需要跨區塊對齊 symbol。

### 個股詳情籌碼視覺化版面

帶 chips 的 stock_news 條目在觀察重點區塊之後、來源列之前渲染「籌碼面」區塊：

- 三列水平長條：外資、投信、自營商。長條寬度按三者絕對值的最大者等比縮放，伺服器端算好百分比以行內寬度輸出，純 HTML div 實作。
- 買超用既有多方紅（tag-bull 色系），賣超用既有空方綠（tag-bear 色系），零圓角零陰影，規線風格與現有條目一致。
- 每列右側標數值，股數換算為張顯示（除以 1,000 四捨五入，正值帶加號）。買賣超為零的列照渲染，長條寬度為零。
- 長條列之下一行小字：三大法人合計、融資增減、融券增減、資料日期。margin_change 或 short_change 缺漏時省略該片段。
- chips 不存在時整塊省略，無佔位文字，與觀察重點區塊的既有慣例一致。

不引入 JS 圖表庫的理由：版面是伺服器端字串拼裝的靜態 HTML，三條 bar 用 div 就能表達，引入圖表庫違反現有零依賴前端的架構與 FT 系靜態版式。

### 判讀引用紀律比照昨收價

SKILL.md 判讀規則新增：籌碼數據只作既有新聞判斷的佐證（對齊個股排序依據第 5 項「籌碼資金面只作佐證加分」），不得單獨構成 call 或 stock_news 條目；正文引用籌碼數字必須標明資料日期；讀取一律 jq/grep 按代號查存檔，不整份進 context。組稿時每個 stock_news 條目按其市場別（上市／上櫃）從對應存檔取值填 chips；查無該代號（新上市、ETF 無融資券等）就省略 chips 或省略對應欄位。抓取失敗照常判讀與發布：該市場個股的 chips 整塊省略，缺漏列入執行回報，報告內容不揭露。

## Implementation Contract

- 行為：POST /api/report 接受 stock_news 條目帶或不帶 chips；帶 chips 且 date 非法回 400 並以索引標明；GET /api/report/{date} 原樣回傳 chips。首頁與歷史頁渲染帶 chips 的條目時出現「籌碼面」區塊（三條法人長條、合計與融資券增減、資料日期），不帶 chips 的條目與既有版面完全一致。
- 資料形狀：chips 為上述七欄位物件，Go 端以指標型別掛在 StockNews 上（nil 即缺漏），JSON 欄位名 date、foreign_net、trust_net、dealer_net、total_net、margin_change、short_change；margin_change 與 short_change 亦為指標型別以區分「零」與「缺漏」。
- 失敗模式：chips 驗證違規走既有 400 errors 陣列與 rejected-schema log；渲染端對 nil chips 靜默省略區塊。
- 驗收：internal/report/schema_test.go 覆蓋 chips 合法、缺漏、date 非法三態；internal/api/report_test.go 覆蓋帶 chips 的 ingest 往返；internal/site/site_test.go 覆蓋有 chips 條目渲染出籌碼區塊與長條寬度、無 chips 條目不出現該區塊；SKILL.md 更新後 go test ./... 全綠。
- 範圍邊界：改 skills/daily-brief/SKILL.md、internal/report/schema.go、internal/site/site.go、internal/site/styles.go 與對應測試。不改 internal/store、internal/admin、internal/api 的程式（驗證邏輯在 report 套件內，API 層自動生效）、不改資料庫。

## Risks / Trade-offs

- [TWSE rwd 端點非 openapi 正式文件的一部分，格式可能變動] → 日期驗證與欄位名核對已能擋住格式漂移；漂移發生時視為抓取失敗，報告照常發布，chips 省略。
- [新增兩個 TWSE 請求拉長節流序列（約多 10 秒以上）] → 蒐集期本就以分鐘計，可接受；序列順序由 skill 條文固定，避免 agent 並行觸發鎖 IP。
- [股數換算張的四捨五入與媒體慣用口徑可能有 ±1 張差異] → 版面標示單位「張」，數值僅供快速掃讀，正文引用仍以來源原始數字為準。
- [T86 外資欄位拆「不含外資自營商」與「外資自營商」兩欄，口徑易搞混] → 設計定死 foreign_net ＝ 兩欄相加（全外資口徑），並寫進 SKILL.md 組稿規則與 spec 情境。
