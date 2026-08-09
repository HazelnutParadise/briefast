## Why

需要一個每個交易日開盤前自動更新的產業與股市新聞報告網站：Cowork agent 從公開新聞判讀個股多空、整理產業摘要，讀者打開網頁就看到當天報告，也能回看歷史。目前 repo 只有文件與 skills，還沒有任何程式。

## What Changes

- 新增 Syralit（Go）前台：首頁即當日報告，依序為盤前總覽、今日觀察、個股多空判斷（短多／短空／長多／長空，每檔附一句理由）、產業新聞摘要、個股新聞詳情（多空判斷的每一檔＋有大消息但無明確判斷者，後者不標標籤），加上固定免責聲明頁尾；歷史頁列出過去報告，點日期以同版面檢視。版面由程式寫死，agent 只提供內容。
- 新增報告接收 API：POST /api/report，bearer key 驗證，一次接收整份報告 JSON（自訂 schema），驗證後寫入 SQLite 資料庫（data/briefast.db，同日覆蓋），並即時推播給所有在線瀏覽器。首頁永遠顯示最新一份報告並標明日期。
- 新增後台 /admin：管理員密碼登入。功能一：API key 管理——建立（可命名）、列表、隨時可重看完整 key、撤銷即時生效；key 以明文存於資料庫 api_keys 表。功能二：報告更新紀錄——記錄哪個 key 在何時推送了哪一天的報告。
- 新增 docker compose 部署設定，data/ 目錄掛 volume 保存 SQLite 資料庫。
- 新增 Cowork 每日流程 skill：定義每個交易日開盤前「爬新聞 → 判讀整理 → 產生報告 JSON → POST 到 API」的步驟與產出格式。
- 明確不使用 Artifact Canvas 與 Artifact DSL：版面固定後 Canvas 沒有存在理由，且客製樣式（多空卡片、紅漲綠跌標籤）需要原生元件才能實作。

## Capabilities

### New Capabilities

- `report-viewing`: 前台首頁（當日報告五區塊＋免責聲明）與歷史頁（日期列表、逐日檢視）的顯示行為。
- `report-ingest-api`: POST /api/report 的驗證、報告 JSON schema、SQLite 持久化與即時推播。
- `admin-panel`: 後台登入、API key 全生命週期管理（含重複檢視明文 key）、報告更新紀錄。
- `daily-brief-skill`: Cowork 每日流程 skill 的步驟、判讀規則與輸出契約。
- `deployment`: docker compose 部署、資料持久化與設定注入方式。

### Modified Capabilities

（無）

## Impact

- Affected specs: 全部為新增——report-viewing、report-ingest-api、admin-panel、daily-brief-skill、deployment。
- Affected code:
  - New: go.mod、main.go、internal/report/schema.go、internal/store/store.go、internal/store/schema.sql、internal/api/report.go、internal/admin/keys.go、internal/admin/log.go、Dockerfile、docker-compose.yml、syralit.toml.example、.gitignore、skills/daily-brief/SKILL.md
  - Modified: AGENTS.md、README.md
  - Removed: 無
- 依賴：新增 Go module 依賴 github.com/HazelnutParadise/syralit 與 modernc.org/sqlite（純 Go SQLite 驅動，免 CGO）。
- 系統：部署主機需可跑 docker compose；Cowork 端需能對外連到網站的 API。
