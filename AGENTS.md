<!-- SPECTRA:START v1.0.2 -->

# Spectra Instructions

This project uses Spectra for Spec-Driven Development(SDD). Specs live in `openspec/specs/`, change proposals in `openspec/changes/`.

## Use `$spectra-*` skills when:

- A discussion needs structure before coding → `$spectra-discuss`
- User wants to plan, propose, or design a change → `$spectra-propose`
- Tasks are ready to implement → `$spectra-apply`
- There's an in-progress change to continue → `$spectra-ingest`
- User asks about specs or how something works → `$spectra-ask`
- Implementation is done → `$spectra-archive`
- Commit only files related to a specific change → `$spectra-commit`

## Workflow

discuss? → propose → apply ⇄ ingest → archive

- `discuss` is optional — skip if requirements are clear
- Requirements change mid-work? `ingest` → resume `apply`

## Parked Changes

Changes can be parked（暫存）— temporarily moved out of `openspec/changes/`. Parked changes won't appear in `spectra list` but can be found with `spectra list --parked`. To restore: `spectra unpark <name>`. The `$spectra-apply` and `$spectra-ingest` skills handle parked changes automatically.

<!-- SPECTRA:END -->

# AGENTS.md

Briefast 的專案操作約定。任何 agent 動工前先讀完這份文件。

## 專案是什麼

每日產業與股市新聞報告網站。兩個組成部分都放在這個 repo：

1. **Syralit 網站**（Go）：以固定版面顯示最新與歷史報告，提供 `/admin/` 管理 API key 與更新紀錄，透過 `POST /api/report` 接收完整報告 JSON。
2. **每日流程 skill**：`skills/daily-brief/SKILL.md` 給 Claude Cowork 執行「蒐集新聞 → 判讀多空 → 組報告 JSON → POST」流程。

報告版面寫死在 Go 程式，agent 只能提供內容。專案不使用 Artifact Canvas 或 Artifact DSL。

## 架構

- `main.go`：開啟 `data/briefast.db`，用自訂 `http.ServeMux` 掛載 API 與三個 Syralit app。
- `internal/report/`：報告型別與驗證。
- `internal/store/`：modernc.org/sqlite 資料層、啟動 migration、reports／api_keys／update_log。
- `internal/api/`：Bearer 驗證、原子 ingest、拒絕留痕、live update 通知。
- `internal/site/`：首頁、歷史頁、固定五區塊與 FT 系規線版面。
- `internal/admin/`：後台登入、明文 key 建立／檢視／撤銷、更新紀錄。

## 必用 skill

- **`syralit-dev`**（`.claude/skills/syralit-dev/`）：寫任何 Syralit 程式碼之前必讀，依實際 v0.7.0 API 使用 rerun、widget、`sy.Handler`、`sy.Shared`、設定與 `sy.AppTest`。
- `syralit-dev` 由 `skills-lock.json` 追蹤，來源是 `HazelnutParadise/syralit`；不要手改安裝的 `SKILL.md`。

## 開發約定

- Import 慣例：`import sy "github.com/HazelnutParadise/syralit"`。
- 報告更新一律走 `POST /api/report` 與有效 Bearer key，不新增未驗證的更新端點。
- 同日報告採全量覆寫；report upsert 與 `ingest_ok` log 必須在同一 transaction。
- API 每次都從 SQLite 查 key 狀態，撤銷必須立即生效；401 與 400 也要寫 update_log。
- `syralit.toml` 是進版控的非機密設定，不含 `[secrets]` 區段。機密與部署變數（`BRIEFAST_ADMIN_PASSWORD`、`BRIEFAST_CONFIG`、`BRIEFAST_PORT`、`BRIEFAST_DB_PATH`）放 repo 根目錄的 `.env`，該檔不進版控，鍵名見 `.env.example`。Compose 會自動讀 `.env`；本機跑 Go 要先自行載入。
- API key 依已定案需求明文保存並可重複檢視。不得在 log、錯誤訊息或文件範例洩漏真實 token。
- 視覺遵循 `DESIGN.md` 與 benchmark：零圓角、零陰影、規線只到大區塊、台股紅漲綠跌永不反轉、每個公開報告頁都有固定免責聲明。
- migration 已進版控後視為不可變；資料結構變更新增下一版 migration，不修改既有 migration。

## 常用指令

```bash
go run .                       # 本機啟動，自訂 mux 會讀版控中的 syralit.toml
go test ./...                  # API 用 httptest，UI 用 sy.NewAppTest
docker compose up -d --build   # 建置與部署，SQLite 放 named volume
```

本機資料庫預設為 `data/briefast.db`，可用 `BRIEFAST_DB_PATH` 覆寫。Compose 的主機 port 用 `BRIEFAST_PORT` 設定，私密設定檔路徑用 `BRIEFAST_CONFIG` 設定。

## Follow-ups

（無）
