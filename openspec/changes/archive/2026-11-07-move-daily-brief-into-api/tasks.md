## 1. 提示欄位與端點總覽

- [x] 1.1 先寫失敗測試：`/api/` 底下每個 JSON 回應（含 400、401、404、405、5xx）都有 `guide` 欄位與 `Link: </api/guide>; rel="help"` 標頭；未知的 `/api/x` 回 JSON 404，不回首頁 HTML。
- [x] 1.2 修改 `internal/api` 的 `writeJSON`，統一加上提示欄位與標頭；新增 `/api/` 預設處理器並掛到 `main.go`，讓 1.1 通過。
- [x] 1.3 先寫失敗測試，再修改 `GET /api/report/{date}`：成功時回 `{"ok":true,"guide":…,"report":…}`，並驗證把 `report` 原樣 POST 回 `/api/report` 會被接受。

## 2. 流程說明端點

- [x] 2.1 把 `skills/daily-brief/SKILL.md` 搬到 `internal/guide/daily-brief.md`，刪掉 frontmatter，用 `go:embed` 嵌入。
- [x] 2.2 改寫說明內容：`seen.py` 的指令換成 `/api/seen/*` 的 curl 範例；`tx_night.py` 換成 `POST /api/tx-night` 的 multipart 範例；報告讀取段落改讀 `report` 欄位；`.env.example` 的指引改成指向 README。
- [x] 2.3 先寫失敗測試，再實作 `GET /api/guide`：有效 key 回 Markdown，401 寫入 `guide_rejected_auth` log，成功不寫 log。
- [x] 2.4 把 `internal/report/skill_example_test.go` 改成讀嵌入的說明，確認報告與法說會 JSON 範例仍能通過驗證。

## 3. 去重 API

- [x] 3.1 新增 `internal/store/migration_v3.sql`（`seen_articles`、`seen_events` 與索引），接到啟動 migration 流程，並寫 store 層測試。
- [x] 3.2 把網址正規化、文字正規化、二連文 MinHash 與相似度估計移植成 Go，用 `seen.py` 的既有行為寫對照測試（追蹤參數、署名、NFKC、照登轉載高於 0.9、不相關文章低於 0.25）。
- [x] 3.3 先寫失敗測試，再實作 `POST /api/seen/peek`、`/record`、`/similar`、`/events`、`/events/check`，包含 401／400 寫 log、成功不寫 log，以及每次請求先清除過期紀錄。

## 4. 夜盤 API

- [x] 4.1 加入 `golang.org/x/text` 依賴，把 `tx_night_test.py` 的案例改寫成 Go 測試（正常、跨週末、換月、零變動、各種錯誤），測試先失敗。
- [x] 4.2 移植 `tx_night.py` 為 Go 解析器並實作 `POST /api/tx-night`，讓 4.1 通過；再用過去晨報執行留下的期交所原始 CSV（2026-10-01、10-16、10-20、10-27、10-29 五次）比對，伺服器結果與當時 `tx_night.py` 的輸出逐欄一致。

## 5. 總覽與 OpenAPI

- [x] 5.1 實作公開的 `GET /api` 與 `GET /api/openapi.json`，從同一份路由表產生，涵蓋所有端點的參數、請求與回應格式、驗證方式與錯誤。
- [x] 5.2 寫測試比對 `main.go` 實際掛載的 `/api` 路由與總覽、OpenAPI 內容，任何一邊多或少都要失敗。

## 6. 移除 skill 與文件更新

- [x] 6.1 刪除 `skills/daily-brief/`，確認 repo 內沒有檔案再引用 SKILL.md、`seen.py`、`tx_night.py`（`openspec/changes/archive/` 的歷史紀錄除外）。
- [x] 6.2 更新 README：Cowork 排程的一行指示、工作區 `.env` 範本（`BRIEFAST_URL`、`BRIEFAST_API_KEY`）、新端點清單；更新 AGENTS.md 的專案說明與架構段落。

## 7. 驗證與收尾

- [x] 7.1 執行 `go test ./...`、`go vet ./...`、`openspec validate move-daily-brief-into-api`、`git diff --check`，全部通過。
- [x] 7.2 本機啟動網站，用 curl 走一次完整流程：讀總覽、讀說明、peek／record／similar／events、上傳夜盤 CSV、讀報告、POST 報告，確認每個回應都有提示。
- [x] 7.3 archive 後把 `openspec/specs/daily-brief-skill/` 改名為 `daily-brief-workflow`，內文的 skill 改成 guide，再跑一次 `openspec validate --specs`。
