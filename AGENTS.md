# AGENTS.md

Briefast 的專案操作約定。任何 agent 動工前先讀完這份文件。

## 專案是什麼

每日產業與股市新聞報告網站。兩個組成部分都放在這個 repo：

1. **Syralit 網站**（Go）：顯示每日報告，內嵌 Artifact Canvas，透過帶驗證的 Artifact API 接受 agent 更新。以 Docker Compose 部署。
2. **每日流程 skills**：給 Claude Cowork 每天執行的 skill，流程是「爬新聞 → 整理成報告 → 產生 Artifact DSL → POST 到網站的 Artifact API」。

## 目前狀態

Repo 剛初始化。網站程式（`go.mod`、`main.go`、`syralit.toml`）、`docker-compose.yml`、每日流程 skill 都還沒建立。動工前先確認現況，不要假設檔案已存在。

## 必用 skills

- **`syralit-dev`**（`.claude/skills/syralit-dev/`）：寫任何 Syralit 程式碼之前必讀。涵蓋 rerun 模型、widget API、ArtifactCanvas、`HandleArtifactAPI`、`syralit.toml` 設定與 `sy.AppTest` 測試。
- **`syralit-artifact-dsl`**（`.claude/skills/syralit-artifact-dsl/`）：產生或驗證 Artifact DSL JSON payload 時必讀。每日流程 skill 產出的報告 payload 必須符合這份 DSL 規格。

兩個 skills 由 `skills-lock.json` 追蹤、安裝自 `HazelnutParadise/syralit`（`.agents/skills/` 是同步副本）。更新方式：`npx skills add HazelnutParadise/syralit`，不要手改 SKILL.md。

## 開發約定

- Import 慣例：`import sy "github.com/HazelnutParadise/syralit"`。
- Artifact 更新一律走認證 API（`sy.HandleArtifactAPI` + agent key）；不開放無驗證的更新端點。
- Agent key 與其他機密放 `syralit.toml` 的 `[secrets]` 或環境變數。含機密的 `syralit.toml` 不得進版控——建 repo 時提供 example 檔並把實際檔案加入 `.gitignore`。
- 每日報告的 Artifact node `id` 要穩定（例如固定為 `headline`、`market-summary`），前端動畫依賴穩定 id。
- Artifact 更新必帶 `expected_revision`，收到 409 時重新取得 spec 再合併，不得直接重送。

## 常用指令

網站建立後適用（目前還跑不了）：

```bash
syralit dev              # 開發模式，熱重載
syralit run              # 生產模式
go test ./...            # 測試（UI 用 sy.NewAppTest headless 測）
docker compose up -d --build   # 部署
```

## Follow-ups

（無）
