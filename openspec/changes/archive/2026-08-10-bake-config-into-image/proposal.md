## Summary

`syralit.toml` 改為在建置時包進映像；Docker Compose 不再掛載設定檔，改由 `.env` 控制資料庫的落地位置。

## Motivation

`syralit.toml` 是 Syralit 框架的設定檔——它決定綁定位址、主題色票與 i18n 字串，是這個程式的一部分，不是使用者設定。目前卻用部署期掛載的方式提供：Dockerfile 的 runtime stage 只 COPY binary，容器內那份設定完全來自 `docker-compose.yml` 的 bind mount。

這導致一個實際的缺陷：映像單獨執行時容器內沒有 `syralit.toml`，Syralit 退回內建預設值 `127.0.0.1:8600`，服務只綁在 loopback，外部完全連不進來。今天能運作純粹是因為 compose 補上了那個掛載。任何不經 compose 的執行方式（直接 docker run、其他編排工具、平台部署）都會拿到一個看似啟動成功、實際無法連線的容器。

`BRIEFAST_CONFIG` 這個變數也強化了錯誤印象：它只是 compose 決定掛哪個主機檔的參數，程式本身從不讀它，卻讓人以為設定檔是可替換的使用者旋鈕。

容器真正需要持久化的只有 SQLite 資料庫，那已經由 named volume 處理。

## Proposed Solution

- Dockerfile 的 runtime stage 加入 `syralit.toml` 的 COPY，讓映像自帶設定。
- `docker-compose.yml` 移除設定檔的 bind mount；資料目錄的掛載來源改為 `${BRIEFAST_DATA:-briefast-data}`，預設維持 named volume，在 `.env` 填主機路徑即改為 bind mount。
- `.env.example`、README、AGENTS.md 移除 `BRIEFAST_CONFIG`，改為說明 `.env` 控制的是埠、機密與資料庫落地位置。
- `syralit.toml` 檔頭註解改為說明它是隨程式一起發佈的設定，不是部署期注入點。
- `deployment` spec 的設定注入需求同步改寫。

## Non-Goals

- 不改設定檔的內容值（標題、主機、埠、主題、i18n 全部維持現狀）。
- 不改 Go 程式讀取設定的邏輯。
- 不改容器內的資料庫路徑（固定 `/app/data/briefast.db`），也不改預設仍走 named volume 的行為。
- 不為環境差異新增設定覆寫機制；有需要時另案處理。
- 不改 `.env` 承載機密與 compose 變數的角色。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `deployment`: 設定檔由部署期掛載改為建置期包進映像，容器只持久化資料庫。

## Impact

- Affected specs: `deployment`
- Affected code:
  - New: （無）
  - Modified: Dockerfile, .dockerignore, docker-compose.yml, .env.example, README.md, AGENTS.md, syralit.toml
  - Removed: （無）
