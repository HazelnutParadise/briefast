## Context

Briefast 的報告是「一天一份、同日全量覆寫」的 JSON 快照，存在 reports 表的 payload。法說會排程與現況不同：公告只提前 1 到 5 天、每天只有當日快照，所以排程要跨天累積；會後還要回填收盤價與命中結果，同一列會被兩個不同時點的流程寫入。硬塞進 report payload 會被隔天的全量覆寫抹掉，因此需要獨立資料列與獨立端點。

資料來源在 2026-09-12 實測確認：

- 未來法說會排程沒有現成資料集。證交所與櫃買 OpenAPI 目錄只有 ESG 年度法說會次數統計；公開資訊觀測站的法人說明會一覽表舊版 ajax 回空頁、新版 api 一律 302。
- 可用來源是每日流程批次 2 已抓的重大訊息（上市 opendata/t187ap04_L、上櫃 openapi/v1/mopsfin_t187ap04_O），符合條款「第12款」即法說會公告，說明欄有固定格式的日期、時間、地點、擇要訊息。單日樣本上市 92 則中 17 則、上櫃 46 則中 2 則，其中八成是受邀券商論壇。
- 基本面：月營收（opendata/t187ap05_L、mopsfin_t187ap05_O）內含上月、去年同月、累計與去年累計；季損益（opendata/t187ap06_L_ci 與各產業變體、mopsfin_t187ap06_O_ci）只有最新一季的累計數，沒有前一季或去年同季。
- 結算用個股日成交：證交所 rwd/zh/afterTrading/STOCK_DAY（date 與 stockNo 參數，回一個月）、櫃買 www/zh-tw/afterTrading/tradingStock（code 與 date 參數，回一個月），兩者皆實測可用。

使用者已鎖定的決策：只列公司自辦場次；材料額外抓月營收、季損益、籌碼、昨收並寫出相對上月、去年同期是變好或變壞；會後驗證一起做並在區塊顯示準確度；準確度以會後第 1 個交易日收盤相對會前一日收盤的方向對照預測，「無法判斷」不計入分母，顯示「近 N 場命中 X 場」而非百分比。

## Goals / Non-Goals

**Goals:**

- 首頁多一個條件出現的「近期法說會」區塊：卡片、預測標籤、準確度、最近結算；點進去有完整會前報告頁。
- 法說會資料獨立存放、獨立端點，登記與結算互不覆蓋。
- 每日流程能自動偵測自辦法說會、抓基本面、組稿、POST，並在會後回填結算。
- 舊資料庫啟動時自動升到第 2 版 migration，既有報告與頁面行為不變。

**Non-Goals:**

- 不列受邀券商論壇場次，也不為它們寫報告。
- 不抓法說會逐字稿、簡報檔（MOPS 不可程式化取得）。
- 不做後台手動增刪法說會。
- 歷史報告檢視與歷史頁不顯示法說會區塊；不做法說會的歷史列表頁。
- 不做前一季或去年同季的季損益結構化比較（資料集沒有），只在新聞載明時以文字帶出。
- 不引入前端圖表庫或新的 Go 依賴。

## Decisions

### 獨立 conferences 表與第 2 版 migration

新增 `internal/store/migration_v2.sql` 建立 conferences 表，主鍵 (symbol, held_on)，欄位：symbol、name、market、held_on、held_at、venue、announced_on、prediction、payload（完整 JSON blob）、pre_close_date、pre_close、post_close_date、post_close、outcome、settled_at、created_at、updated_at；索引 (held_on)、(settled_at, held_on)。migrate 函式改為依序套用版本 1 與版本 2，已在版本 1 的資料庫只補套版本 2，全部在同一交易中完成。

沿用 reports 表「索引欄位 + payload blob」的形狀，理由：首頁列表與準確度統計只需要索引欄位，報告頁才解 payload；報告內容欄位（summary_md、fundamentals、chips、sources）隨 skill 演進時不必再改表。替代方案「全部拆成關聯欄位」被否決：沒有跨欄查詢需求，只增加 migration 次數。

不加 deleted_at：專案沒有刪除路徑，reports 表也沒有；日期過了的場次自然退出首頁，報告頁仍可連結。

寫入鍵與時間戳沿用 store 套件現有的固定寬度 storedTimeLayout。

### 登記與結算分成兩個寫入路徑

`POST /api/conference` 只覆寫 brief 相關欄位（name、market、held_at、prediction、payload、updated_at），用 ON CONFLICT DO UPDATE 明列欄位、不動結算欄位；`POST /api/conference/settle` 只寫 pre_close_date、pre_close、post_close_date、post_close、outcome、settled_at。理由：同一列被兩個時點寫入，用兩個端點各自明列欄位，就不會出現「重送 brief 把結算抹掉」或「結算把 brief 蓋回舊版」。

outcome 由伺服器依存好的 prediction 計算：bull 且 post > pre 為 hit；bear 且 post < pre 為 hit；bull 或 bear 其餘皆 miss（含持平）；none 存空字串。理由：判定規則寫死在站台，skill 不能自己改判。

`GET /api/conferences` 回 upcoming 與 pending_settlement 兩個陣列，只帶索引欄位不帶 brief 內文，讓 skill 一次拿到「不必重登記的」與「該結算的」。

所有動作沿用 update_log：conference_ok、conference_rejected_auth、conference_rejected_schema、settle_ok、settle_rejected_auth、settle_rejected_schema、list_rejected_auth；report_date 欄位放 held_on。

### conference brief 的資料形狀與驗證

新增 `internal/report/conference.go` 定義 Conference 型別與 Validate，JSON 欄位：symbol、name、market、held_on、held_at、venue、announced_on、prediction、headline、summary_md、watch_md、fundamentals（選填）、chips（選填，重用既有 Chips 型別）、sources、generated_at。fundamentals 內含選填 revenue（month、current、prev_month、last_year_month、ytd、last_year_ytd，int64 千元）與選填 quarter（label、revenue、operating_income、net_income、eps 字串）。

百分比變化由站台計算，payload 不存百分比。理由：來源資料集本身帶百分比，但站台自算能保證顯示口徑一致，也避免 skill 抄錯欄位。

prediction 用 bull、bear、none 而非沿用 short_bull 等四值：法說會預測只有一個時間尺度（會後第 1 個交易日），四值裡的長短期區分沒有意義。

### 首頁近期法說會區塊與會前報告頁

區塊插在個股多空判斷之後、廣告版位之前，只在首頁且有未來場次時渲染；站台以 Asia/Taipei 的今天判斷「未來」（held_on >= today）。歷史檢視不渲染，因為它描述的是「現在」不是某一天的報告。

卡片沿用 calls 四直欄規線版式（call-col 的 grid 與斷點），標籤沿用 tag 樣式：bull 紅、bear 綠、none 描邊灰。準確度取最近 20 場已結算且 prediction 非 none 的場次，文案「近 N 場命中 X 場」；最近結算最多列 5 場，漲跌百分比 (post - pre) / pre 四捨五入到小數兩位，正紅負綠。

報告頁掛在 `/conference/`（自訂 mux 以 StripPrefix 掛第四個 Syralit app），query 參數 symbol 與 date；頁面結構：masthead（回首頁導覽）、已結算時的結果帶、標題區（股名代號、日期時間地點、預測標籤、headline）、判斷依據、基本面表格、觀察重點、籌碼面（重用 renderChips）、來源、產生時間與公告日期、免責頁尾。找不到或參數不合法一律 not-found 狀態，不 panic。

SEO：新增 PageConference 種類的 DocumentFunc，標題「{name} 法說會前預測｜Briefast」，描述取 summary_md 純文字截 150 字；sitemap 加入每一場 brief 頁，lastmod 用 updated_at。seo.Reports 介面加一個 ListConferences 方法（回代號、日期、updated_at）。

### skill 的法說會流程放在判讀之後、組稿之前

SKILL.md 新增「法說會」一節，位置在第 2 步判讀與第 3 步組報告之間，包含四件事：

1. 偵測：從批次 2 存檔用 jq 篩 符合條款 為 第12款 的條目，解析說明欄四個固定標籤；自辦判斷規則（主旨與擇要訊息不含 受邀、應…之邀、邀請、參加、論壇、Summit、Conference、投資論壇，且日期是單一日期）；先 GET /api/conferences，已在 upcoming 且日期時間地點相同者不重送。
2. 基本面：月營收與季損益資料集整包存檔、jq 按代號取值、TWSE 請求併入節流序列；季損益是累計數，標籤寫「2026Q2 累計」；前一季或去年同季比較只在新聞載明時寫並引用。
3. 組稿與 POST：brief JSON 完整範例；判斷先行；依據不足就 none；與當日 calls 的時間尺度差異要點明；每筆 POST 回報結果。
4. 結算：pending_settlement 每筆抓個股當月（跨月抓兩月）日成交，取會前最後一個交易日與會後第一個交易日的收盤，POST settle；會後交易日尚未出現就留待下次。

放在判讀之後的理由：偵測要用批次 2 存檔，組稿要用判讀階段已整理的個股新聞；放在組報告之前是讓報告 watch_md 可以提到「某公司某日法說」而不倒過來依賴 brief。

## Implementation Contract

- 行為：`POST /api/conference` 帶合法 brief 回 200 {ok, symbol, held_on}，違規回 400 errors 陣列，未驗證回 401；`POST /api/conference/settle` 回 200 {ok, symbol, held_on, outcome}，查無回 404，日期不在會期兩側回 400；`GET /api/conferences` 回 {upcoming: [], pending_settlement: []}。首頁在有未來場次時於 calls 與 ad 之間出現 data-section="conferences" 區塊；`/conference/?symbol=&date=` 渲染 brief 或 not-found；sitemap 含 brief 頁；三個端點與頁面的中繼資料由 DocumentFunc 產生。
- 資料形狀：Conference 型別與 fundamentals 形狀如「conference brief 的資料形狀與驗證」一節；conferences 表欄位如「獨立 conferences 表與第 2 版 migration」一節；store 提供 UpsertConference、SettleConference、ConferenceByKey、ListUpcomingConferences(today)、ListPendingSettlement(today)、ListSettledConferences(limit)、ListConferences(全部索引欄位) 方法。
- 失敗模式：驗證違規走既有 400 與 rejected-schema log；settle 對未登記場次 404 不落庫；站台查詢失敗顯示「無法載入」不 panic；migration 失敗阻止啟動（與現況一致）。
- 驗收：`go test ./...` 全綠、`go vet ./...` 無輸出；store 測試覆蓋版本 1 資料庫升到版本 2 與 upsert 不覆蓋結算；report 測試覆蓋 brief 驗證矩陣；api 測試覆蓋三端點的 200、400、401、404 與 log 動作；site 測試覆蓋區塊出現與省略、卡片順序與標籤色相、準確度計數、報告頁各區塊、not-found；seo 測試覆蓋 brief 頁標題與 sitemap 條目；瀏覽器實際檢視首頁區塊與報告頁深淺兩色；`spectra validate` 通過。
- 範圍邊界：改 internal/store、internal/report、internal/api、internal/site、internal/seo、main.go、SKILL.md、AGENTS.md（架構段落補一句）、DESIGN.md（元件段補法說會卡片與結果帶）。不改 internal/admin、不改既有 migration、不改 report schema、不改歷史頁。

## Risks / Trade-offs

- [自辦與受邀的判斷靠關鍵字，會有漏判或誤判] → 規則寫進 SKILL.md 並附分類表；誤判成受邀只會少列，誤判成自辦會多一張 none 卡片，兩者都不破壞資料；上線兩週後用實際公告校正關鍵字。
- [季損益只有累計數，「相對上季」做不到結構化] → 設計明講只做月營收三項比較加最新季累計數，季比較交給新聞引用；不硬算。
- [結算依賴 skill 每天跑；連假或漏跑會延後結算] → pending_settlement 不設過期，下次跑就補；個股日成交端點回一整月，跨月抓兩月即可補到。
- [首頁多一個 DB 查詢與時區換算] → 只查索引欄位、最多幾十列；today 用 time.LoadLocation("Asia/Taipei")，容器映像已含 tzdata 則可用，否則 fallback 固定 +08:00。
- [TWSE 個股日成交端點在 rwd 路徑，非 OpenAPI 正式文件] → 結算失敗只是延後，不影響報告發布；併入節流序列避免鎖 IP。
- [首頁固定五區塊的 spec 要改，屬 BREAKING] → 歷史頁不變，只有首頁在有場次時多一區，測試斷言五區順序仍成立。

## Migration Plan

1. 部署新版映像；啟動時 migrate 自動套版本 2，reports 與 api_keys 不動。
2. 在 Cowork 更新 daily-brief skill 內容，下一個交易日開始登記。
3. 回滾：換回舊映像即可，conferences 表留在 SQLite 不影響舊版；舊版 migrate 只檢查版本 >= 1。

## Open Questions

（無，四項決策已由使用者於 docs/plans/2026-09-12-main-upcoming-conferences-plan.md 鎖定。）
