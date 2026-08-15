#!/bin/sh
# 資料目錄可能是 bind mount，而 Docker 一律以 root 建立掛載目錄，覆蓋映像中原本的擁有者。
# 因此以 root 進入時先把該目錄交給應用使用者，再降權執行；容器若被指定以非 root 啟動就直接執行。
set -e

data_dir="$(dirname "${BRIEFAST_DB_PATH:-/app/data/briefast.db}")"

if [ "$(id -u)" = "0" ]; then
    mkdir -p "$data_dir"
    chown -R briefast:briefast "$data_dir"
    exec su-exec briefast "$@"
fi

exec "$@"
