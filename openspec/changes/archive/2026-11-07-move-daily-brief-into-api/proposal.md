## Why

每日晨報的流程現在寫在 `skills/daily-brief/SKILL.md`，還附兩支 Python 腳本。Cowork 的工作區要另外裝這包 skill，伺服器驗證規則一改，skill 也要跟著改一次，兩邊很容易對不上。把流程說明和腳本功能都放進網站 API 之後，agent 只要知道網址和 API key，就能從 API 讀到完整做法並完成整個流程，這包 skill 就可以刪掉。

## What Changes

- 新增 `GET /api/guide`，用 Bearer key 讀取每日晨報的完整流程說明（Markdown）。說明內容從 SKILL.md 搬進 Go 程式並一起部署，原本叫 agent 執行腳本的段落改成呼叫新端點。
- 新增 `GET /api` 端點總覽和 `GET /api/openapi.json`，列出所有端點、方法、驗證方式、參數與錯誤。
- `/api/` 底下每個 JSON 回應都多一個 `guide` 欄位，內容是一句提示：先用同一把 key 讀 `GET /api/guide`，再呼叫其他端點。回應標頭也加上 `Link: </api/guide>; rel="help"`。未知的 `/api/...` 路徑改回 JSON 404，同樣附提示，不會再掉到首頁 HTML。
- 新增去重端點，取代 `seen.py`：`POST /api/seen/peek`、`POST /api/seen/record`、`POST /api/seen/similar`、`POST /api/seen/events/check`、`POST /api/seen/events`。已讀紀錄改存 SQLite，文章保留 90 天、事件保留 30 天。
- 新增 `POST /api/tx-night`，上傳兩份期交所 TX CSV，由伺服器算出夜盤訊號，取代 `tx_night.py`。
- **BREAKING**：`GET /api/report/{date}` 成功時改回 `{"ok":true,"guide":"…","report":{…}}`，報告內容放在 `report` 欄位。目前只有 daily-brief 流程會呼叫這個端點，新說明會同步改成讀 `report`。
- 刪除 `skills/daily-brief/`（SKILL.md、`.env.example`、`seen.py`、`tx_night.py`、`tx_night_test.py`）。Cowork 排程只要留一行指示：讀工作區 `.env` 的 `BRIEFAST_URL` 和 `BRIEFAST_API_KEY`，先呼叫 `GET /api/guide` 再照做。README 提供這行指示和 `.env` 範本。

## Capabilities

### New Capabilities

- `agent-api-guide`：端點總覽、OpenAPI 文件、流程說明端點、每個 API 回應的提示欄位，以及未知 API 路徑的 JSON 404。
- `news-dedup-api`：文章已讀查詢、處理紀錄、相似內文比對、跨天事件登記，資料存在 SQLite 並定期清除過期紀錄。
- `tx-night-api`：上傳兩份 TX 每日行情 CSV，回傳同契約的夜盤漲跌方向。

### Modified Capabilities

- `daily-brief-skill`：流程改由 `GET /api/guide` 提供，不再是 repo 裡的 skill 檔；`.env` 範本移到 README；外電跨語言去重改用 `news-dedup-api` 記錄；TX 夜盤解析改由 `tx-night-api` 負責。
- `report-ingest-api`：讀取端點成功時把報告包在 `report` 欄位，並附上 `guide` 提示。

## Impact

- 程式：`main.go` 路由、`internal/api/`（新增 guide、overview、seen、tx-night 處理器，修改 read 回應與 `writeJSON`）、新增 `internal/guide/`（嵌入流程說明）、`internal/store/`（新增 `migration_v3.sql` 的去重資料表）。
- 依賴：新增 `golang.org/x/text`，用來把期交所 CP950 編碼的 CSV 轉成 UTF-8。
- 測試：`internal/report/skill_example_test.go` 改讀嵌入的說明；`tx_night_test.py` 的案例改寫成 Go 測試；新增各端點的 httptest。
- 營運：Cowork 工作區不再安裝 daily-brief skill，排程 prompt 換成 README 的一行指示。工作區原本的 `seen.jsonl` 不匯入，去重紀錄從部署當天重新累積。
