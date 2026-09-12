## 1. 資料層：獨立 conferences 表與第 2 版 migration

- [x] 1.1 依 design.md「獨立 conferences 表與第 2 版 migration」在 `internal/store/store_test.go` 先寫失敗測試，覆蓋 Conference storage：新資料庫開啟後 schema_migrations 最高版本為 2 且 conferences 表存在；先以只套版本 1 的資料庫再開啟，版本補到 2 且 reports 資料不變。驗證：`go test ./internal/store/` 出現預期失敗。
- [x] 1.2 新增 `internal/store/migration_v2.sql` 並讓 migrate 依序套用版本 1 與 2（同一交易）。行為：任一版本的既有資料庫啟動後都到版本 2，既有 migration 檔一字不改。驗證：1.1 測試全綠。
- [x] 1.3 在 `internal/store/store.go` 新增 Conference 列型別與 UpsertConference、SettleConference、ConferenceByKey、ListUpcomingConferences、ListPendingSettlement、ListSettledConferences、ListConferences 方法，並補測試。行為：UpsertConference 對同一 (symbol, held_on) 覆寫 brief 欄位但保留 pre_close、post_close、outcome、settled_at；SettleConference 只寫結算欄位並回傳 ErrNotFound 給未登記場次；三個列表依 held_on 排序且以傳入的 today 字串切分未來與待結算。驗證：`go test ./internal/store/` 全綠，含「先結算再重送 brief 結算仍在」的斷言。

## 2. Brief schema：conference brief 的資料形狀與驗證

- [x] 2.1 依 design.md「conference brief 的資料形狀與驗證」在 `internal/report/conference_test.go` 先寫失敗測試，覆蓋 Conference brief ingestion endpoint 的驗證矩陣：合法 brief 無錯誤；prediction 非 bull/bear/none、market 非 twse/tpex、held_at 非 HH:MM、headline 空白、source 缺 url、chips.date 非法、generated_at 非 RFC 3339 各回對應訊息；fundamentals 缺漏接受。驗證：`go test ./internal/report/` 出現預期失敗。
- [x] 2.2 新增 `internal/report/conference.go` 定義 Conference、Fundamentals、RevenueFigures、QuarterFigures 型別與 Validate，重用既有 Chips 與 Source 型別。行為：JSON 欄位名為 symbol、name、market、held_on、held_at、venue、announced_on、prediction、headline、summary_md、watch_md、fundamentals、chips、sources、generated_at；違規訊息以欄位名開頭。驗證：2.1 測試全綠。

## 3. API：登記與結算分成兩個寫入路徑

- [x] 3.1 在 `internal/api/conference_test.go` 先寫失敗測試，覆蓋 Conference brief ingestion endpoint（200 回 symbol 與 held_on、update_log 出現 conference_ok 且 report_date 為 held_on、notifier 被呼叫；400 列出全部違規並記 conference_rejected_schema；401 記 conference_rejected_auth；重送已結算場次後結算欄位不變）、Conference settlement endpoint（hit、miss、持平為 miss、none 為空 outcome；pre_close_date 等於 held_on 回 400；未登記回 404；settle_ok log）、Conference listing endpoint（upcoming 與 pending_settlement 依今天切分且不含 brief 內文；空庫回兩個空陣列；401 記 list_rejected_auth）。驗證：`go test ./internal/api/` 出現預期失敗。
- [x] 3.2 新增 `internal/api/conference.go` 實作三個 handler，沿用既有 authenticate、decode 限制 2 MiB 與 DisallowUnknownFields、writeJSON；listing 的今天以 Asia/Taipei 計算並可由測試注入。行為：如 3.1 所列。驗證：3.1 測試全綠。
- [x] 3.3 在 `main.go` 掛載 POST /api/conference、POST /api/conference/settle、GET /api/conferences，並在 `main_test.go` 補路由測試。行為：三個路徑各回對應 handler 的 405 或 401 而非首頁 HTML。驗證：`go test .` 全綠。

## 4. 首頁近期法說會區塊與會前報告頁

- [x] 4.1 依 design.md「首頁近期法說會區塊與會前報告頁」在 `internal/site/site_test.go` 先寫失敗測試，覆蓋 Fixed report layout（有未來場次時 data-section="conferences" 位於 calls 之後、ad 之前，五個固定區塊順序不變；無未來場次時整段不出現也無佔位；HistoryPage 日期檢視不出現該區塊）、Upcoming conferences section（卡片依 held_on 再 symbol 排序、bull 標籤紅色相、bear 綠色相、none 描邊、連結指向 /conference/?symbol=&date=；準確度只計 bull/bear 且取最近 20 場，文案「近 N 場命中 X 場」，無已結算時省略；最近結算最多 5 場、百分比正紅負綠、命中／落空）、Live update on conference changes（Notify 後已開啟的 Home 重新渲染出新卡片）。色彩斷言解析色值色相，不只比對選擇器字串。驗證：`go test ./internal/site/` 出現預期失敗。
- [x] 4.2 在 `internal/site/site.go` 實作 renderConferences 並接進 Home（不接進 HistoryPage），在 `internal/site/styles.go` 補卡片、標籤 neutral、準確度與結算列樣式（零圓角零陰影、沿用 call-col 斷點）。行為：如 4.1 所列。驗證：4.1 測試全綠。
- [x] 4.3 在 `internal/site/site_test.go` 先寫失敗測試，覆蓋 Conference brief page：完整 brief 渲染出預測標籤、判斷依據、基本面表格（當月、上月、去年同月、累計、去年累計，三個百分比正紅負綠，數值千分位）、季損益列、觀察重點、籌碼面、來源、產生時間與公告日期、免責頁尾；已結算時頁首結果帶含兩個收盤與日期、百分比、命中／落空／未列入統計；fundamentals 缺漏時無基本面區塊；查無或參數不合法時 not-found 且有回首頁連結。驗證：`go test ./internal/site/` 出現預期失敗。
- [x] 4.4 在 `internal/site/site.go` 新增 Conference page function 與 ConferencePage(symbol, date) 可測入口，重用 masthead、renderChips、renderEntryWatch、footer；在 `main.go` 以 StripPrefix 掛載 `/conference/` 並加 GET /conference 導向。行為：如 4.3 所列，且 `/conference` 無斜線時 307 到 `/conference/`。驗證：4.3 測試全綠，`go test .` 路由測試全綠。
- [x] 4.5 依 DESIGN.md 啟動本機服務，以 sqlite 塞入含未來場次、已結算場次與完整 brief 的測試資料，用瀏覽器檢視首頁區塊與報告頁的深淺兩色。行為：卡片規線與 calls 對齊、紅漲綠跌未反轉、結果帶可讀、767px 單欄不破版。驗證：深淺色各一張首頁與報告頁截圖並逐項比對，驗完刪除測試資料還原資料庫。

## 5. 爬蟲中繼資料

- [x] 5.1 在 `internal/seo/seo_test.go` 先寫失敗測試，覆蓋 Each public page carries its own title and description 的法說會情境（brief 存在時標題含公司名、法說會前預測與站名，描述取 summary_md 純文字截 150 字；查無時退回站台預設且 canonical 指首頁）與 Site serves a sitemap covering every report 的 brief 條目（每場一筆 /conference/?symbol=&date= 且 lastmod 為 updated_at；查詢失敗回 500 無部分 XML）。驗證：`go test ./internal/seo/` 出現預期失敗。
- [x] 5.2 在 `internal/seo/meta.go` 新增 PageConference 種類與對應 metaFor 分支，`internal/seo/endpoints.go` 的 sitemap 加入 brief 頁，Reports 介面加 ListConferences；在 `main.go` 為 conference app 掛上 DocumentFunc。行為：如 5.1 所列。驗證：5.1 測試全綠。

## 6. 每日流程 skill

- [x] 6.1 依 design.md「skill 的法說會流程放在判讀之後、組稿之前」在 `skills/daily-brief/SKILL.md` 新增法說會一節的前兩段：Conference announcement detection（jq 篩第12款、解析四個固定標籤、自辦判斷關鍵字與分類表、先 GET /api/conferences 不重送、每市場回報筆數）與 Conference fundamentals reference data（月營收與季損益端點含產業變體、存檔按代號查值、TWSE 併入節流序列、季損益為累計數的標籤規則、前一季或去年同季只在新聞載明時寫並引用、缺漏列入執行回報）。驗證：對照 `openspec/changes/upcoming-conferences/specs/daily-brief-skill/spec.md` 前兩個 requirement 的每個 scenario 逐條核對條文涵蓋。
- [x] 6.2 在 `skills/daily-brief/SKILL.md` 同一節補後兩段：Pre-conference brief composition（brief JSON 完整範例與欄位名對齊 `internal/report/conference.go`、判斷先行、依據不足即 none、月營收三項比較方向必寫、與當日 calls 的時間尺度差異、POST /api/conference 與回報）與 Conference settlement after the event（pending_settlement 逐筆抓個股當月日成交、跨月抓兩月、取會前最後與會後第一個交易日收盤、POST /api/conference/settle、會後交易日未出現就留待、結算機制不進報告內容）。並在送出前檢查清單加入法說會相關項目。驗證：對照 spec 後兩個 requirement 的每個 scenario 逐條核對，並確認 JSON 範例能通過 `internal/report` 的 Validate（以測試或臨時程式餵入）。

## 7. 文件與收斂驗證

- [x] 7.1 更新 `AGENTS.md` 架構段（conferences 表與 /conference/ 路由、三個 API 端點）與 `DESIGN.md` 元件段（法說會卡片、neutral 標籤、結果帶）。行為：文件描述與實作一致。驗證：內容審閱對照 main.go 路由與 styles.go 類別名。
- [x] 7.2 全專案測試與交付前檢查。行為：`go test ./...` 全綠，`go vet ./...` 無輸出，舊報告渲染與 ingest 行為不變。驗證：兩道指令的實際輸出。
- [x] 7.3 執行 `spectra validate upcoming-conferences` 與 `spectra analyze upcoming-conferences`，確認 artifacts 一致無缺口。驗證：兩道指令回報 valid 且無未覆蓋的 requirement。
