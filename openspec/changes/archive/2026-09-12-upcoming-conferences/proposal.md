## Why

讀者開盤前掃報告時，看不到哪幾家公司這幾天要開自辦法說會，也沒有一份在法說會前就寫好、事後可以對答案的方向判斷。現有報告只在個股 `watch_md` 零星提到「追蹤下一次法說」，而報告是同日全量覆寫的單日快照，放不下會跨天累積、還要在會後回填結果的法說會排程。

## What Changes

- 新增法說會資料表與 API：以「代號＋法說會日期」為鍵的獨立資料列，透過現有 Bearer 驗證的 `POST /api/conference` 登記或覆寫一場自辦法說會的會前報告，`POST /api/conference/settle` 回填會前一日與會後第一個交易日的收盤價並由伺服器判定命中或落空，`GET /api/conferences` 讓每日流程查詢已登記的未來場次與待結算場次。所有寫入與拒絕都寫進既有 update_log。
- 首頁在個股多空判斷與廣告版位之間新增「近期法說會」區塊：每張卡片顯示公司名、代號、法說會日期時間、會前預測方向標籤（會後看漲紅、會後看跌綠、無法判斷灰）與一句話標題，點卡片進入會前報告頁。區塊標題列顯示「近 N 場命中 X 場」的預測準確度，卡片下方列出最近結算的場次與實際漲跌。沒有未來場次時整個區塊省略。
- 新增會前報告頁 `/conference/?symbol=代號&date=日期`：預測方向、判斷依據（依本次抓取的月營收與最新季損益算出相對上月、去年同月、去年累計的變好或變壞，加上近期消息與籌碼）、法說會要看的重點、來源與資料日期、固定免責聲明；已結算的場次在頁首顯示會前收盤、會後收盤、實際漲跌與命中結果。頁面帶自己的 title、description、canonical 與 Open Graph 中繼資料，並列入 sitemap。
- 每日流程 skill 新增法說會步驟：從批次 2 已抓的上市與上櫃重大訊息「第 12 款」條目解析法說會日期、時間、地點與擇要訊息，只登記公司自辦場次（受邀券商論壇一律略過），為每場新登記的場次抓取月營收與最新季損益、昨收與籌碼後寫會前報告；對已過會期且尚未結算的場次，用證交所與櫃買的個股日成交資訊取得會前一日與會後第一個交易日收盤價並回填。
- **BREAKING**：首頁固定版面從五個內容區塊改為五個固定區塊加一個條件出現的法說會區塊；歷史報告檢視不受影響。

## Capabilities

### New Capabilities

- `conference-brief-api`: 法說會會前報告的登記、結算回填與查詢端點，含 schema 驗證、持久化、命中判定與 update_log 紀錄。

### Modified Capabilities

- `report-viewing`: 首頁新增條件出現的近期法說會區塊與預測準確度，新增會前報告頁與已結算結果的呈現。
- `daily-brief-skill`: 新增法說會登記、自辦與受邀判斷、基本面抓取與會前報告組稿、會後結算回填的流程條文。
- `crawler-metadata`: 會前報告頁的標題、描述、canonical 與社群分享中繼資料，sitemap 涵蓋每一場法說會報告頁。

## Impact

- Affected specs: `conference-brief-api`（新）、`report-viewing`、`daily-brief-skill`、`crawler-metadata`
- Affected code:
  - New: `internal/store/migration_v2.sql`、`internal/report/conference.go`、`internal/report/conference_test.go`、`internal/api/conference.go`、`internal/api/conference_test.go`
  - Modified: `internal/store/store.go`、`internal/store/store_test.go`、`internal/site/site.go`、`internal/site/styles.go`、`internal/site/site_test.go`、`internal/seo/meta.go`、`internal/seo/endpoints.go`、`internal/seo/seo_test.go`、`main.go`、`main_test.go`、`skills/daily-brief/SKILL.md`、`AGENTS.md`、`DESIGN.md`
  - Removed: （無）
- 資料庫：新增第 2 版 migration 建立 conferences 表，既有 migration 不動；SQLite 單機部署，migration 於啟動時自動套用。
