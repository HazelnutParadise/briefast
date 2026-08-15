## Problem

把 `BRIEFAST_DATA` 指向主機目錄後容器起不來，日誌反覆出現 `ping sqlite: unable to open database file (14)`，程式在開資料庫時就失敗，從未進入監聽。改用 named volume 才能運作，等於新加的主機掛載能力實際不可用。

## Root Cause

容器以非 root 的 `briefast` 使用者執行，而 Docker 在建立 bind mount 的目標目錄時一律以 **root** 建立。掛載發生在容器啟動時，會直接覆蓋映像中 `/app/data` 原本已 chown 給 `briefast` 的擁有者，因此應用使用者對該目錄沒有寫入權限，SQLite 回 `SQLITE_CANTOPEN`。

主機目錄若事先存在，情況同樣不會自動變好：它的擁有者是建立它的主機使用者，UID 幾乎不會與容器內的 `briefast` 相同。

`user:` 映射解不了這件事——即使把容器使用者對成主機 UID，Docker 仍會以 root 建立不存在的掛載目錄，首次啟動照樣失敗。

## Proposed Solution

改用官方映像常見的 entrypoint 降權模式：

- 新增 `docker-entrypoint.sh`：容器以 root 進入，建立並把資料目錄的擁有者設為 `briefast`，再以 `su-exec` 降權執行應用。若容器被指定以非 root 啟動（例如 compose 設了 `user:`），entrypoint 略過修正步驟直接執行，不強制要求 root。
- Dockerfile 安裝 `su-exec`、加入 entrypoint 並移除 `USER briefast`，改由 entrypoint 負責降權；應用本身仍以 `briefast` 執行，不以 root 執行。
- 資料目錄從既有的 `BRIEFAST_DB_PATH` 推導，不新增設定項。

如此不論主機目錄是否存在、屬於誰，首次啟動即可運作，且不需要任何手動 chown。

## Non-Goals

- 不讓應用程式以 root 執行；root 僅存在於 entrypoint 修正權限的瞬間。
- 不新增設定變數；資料目錄仍由 `BRIEFAST_DB_PATH` 決定。
- 不改 named volume 的行為，未設 `BRIEFAST_DATA` 時維持現狀。
- 不改應用程式碼與資料庫結構。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `deployment`: 容器啟動時自我修正資料目錄權限，使主機目錄掛載免手動設定即可運作。

## Impact

- Affected specs: `deployment`
- Affected code:
  - New: docker-entrypoint.sh
  - Modified: Dockerfile, README.md
  - Removed: （無）
