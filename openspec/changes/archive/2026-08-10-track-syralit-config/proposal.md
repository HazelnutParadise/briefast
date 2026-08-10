## Summary

`syralit.toml` 改為進版控的正式設定檔並移除 `[secrets]` 區段；機密與部署變數集中到 `.env`，該檔同時驅動 Docker Compose 的掛載與埠設定。

## Motivation

目前的設定是「範本 + 各自複製」模式：repo 只放 `syralit.toml.example`，`.gitignore` 排除真檔，使用者自己 copy 一份填密碼。這個模式在本專案沒有帶來對應價值——設定檔內容（標題、主機、埠、主題色票、i18n 字串）全部是專案共用的非機密設定，每個部署環境都一樣；唯一的機密是管理員密碼，而它本來就支援 `BRIEFAST_ADMIN_PASSWORD` 環境變數。

代價則是實際存在的：範本與真檔容易分歧、README 兩處都要教 copy 步驟、compose 預設掛載範本檔（一個名字叫 example 的檔案卻是正式部署的預設值），新加入的人要多一個步驟才能跑起來。

## Proposed Solution

- `syralit.toml.example` 更名為 `syralit.toml` 並納入版控；`.gitignore` 移除 `syralit.toml` 排除（`.env` 的排除保留）。
- 進版控的 `syralit.toml` 移除整個 `[secrets]` 區段，只留非機密設定；`sy.Secrets` 會回退到環境變數，因此管理員密碼改由 `.env` 的 `BRIEFAST_ADMIN_PASSWORD` 提供。
- `.env` 成為機密與部署變數的單一來源，同時驅動 Docker Compose：Compose 原生會讀專案目錄的 `.env` 做變數替換，所以 `BRIEFAST_CONFIG`（掛載來源）、`BRIEFAST_PORT`、`BRIEFAST_ADMIN_PASSWORD` 都寫在 `.env` 即可生效，不需額外參數。
- repo 根目錄新增 `.env.example` 記錄所有鍵名與用途；`.env` 本身維持被 `.gitignore` 排除。
- `docker-compose.yml` 的預設掛載來源由範本檔改為 `syralit.toml`。
- README 的本機開發與 Docker 部署兩節移除複製範本的步驟，改為以 `.env` 提供設定；本機 `go run` 需先把 `.env` 載入環境，README 附上做法。
- AGENTS.md 的設定約定同步更新。

## Non-Goals

- 不改 Go 程式的設定讀取邏輯，`syralit.toml` 的解析與 `BRIEFAST_CONFIG` 覆寫行為不變。
- 不改 Go 程式讀取 `.env` 的能力——程式本身不解析 `.env`，本機開發靠 shell 載入，Compose 靠原生支援。`.env` 仍被 `.gitignore` 排除。
- 不引入額外的秘密管理工具或加密設定檔。
- 不改主題色票、i18n 字串等設定內容本身。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `deployment`: 設定注入方式由「範本加本機副本」改為「設定檔進版控、密碼走環境變數」。

## Impact

- Affected specs: `deployment`
- Affected code:
  - New: syralit.toml, .env.example
  - Modified: docker-compose.yml, README.md, AGENTS.md
  - Removed: syralit.toml.example
