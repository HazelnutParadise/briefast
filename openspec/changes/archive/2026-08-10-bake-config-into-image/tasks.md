## 1. 映像自帶設定

- [x] 1.0 從 .dockerignore 移除 syralit.toml 排除條目，使建置脈絡包含該檔；同時加入 .env 排除，避免機密進入建置脈絡。驗證方式：docker build 成功且容器內存在 /app/syralit.toml，.dockerignore 含 .env 條目。

- [x] 1.1 依 Configuration injection 與 design 的「設定檔包進映像」決策，在 Dockerfile 的 runtime stage 加入 syralit.toml 的 COPY，並確保執行使用者可讀取。驗證方式：建置映像後在容器內確認 /app/syralit.toml 存在且內容與 repo 版本相同。
- [x] 1.2 依 design 的「容器只持久化資料庫」決策，從 docker-compose.yml 移除設定檔的 bind mount，並把資料目錄的掛載來源改為 ${BRIEFAST_DATA:-briefast-data}，讓 .env 決定資料落地位置；容器內資料庫路徑維持固定。驗證方式：docker compose config 在未設 BRIEFAST_DATA 時展開為 named volume、設為主機路徑時展開為該路徑的 bind mount，兩者皆無指向 /app/syralit.toml 的掛載，埠與 BRIEFAST_ADMIN_PASSWORD 仍取自 .env。
- [x] 1.3 以不帶任何掛載的方式啟動建置好的映像，確認服務綁在 0.0.0.0:8600、從主機取得首頁 200——即修正「映像單獨執行只綁 loopback、外部連不進來」的缺陷。驗證方式：docker run 映像並從主機 curl 首頁取得 200。

## 2. 移除 BRIEFAST_CONFIG 與文件同步

- [x] 2.1 依 design 的「移除 BRIEFAST_CONFIG」決策，從 .env.example、README 與 AGENTS.md 移除該變數的所有敘述與範例，改為說明 .env 控制的是埠、機密與 BRIEFAST_DATA 資料落地位置；並把 syralit.toml 檔頭註解改為說明它隨程式一起發佈、不是部署期注入點。驗證方式：grep 確認活躍文件（排除 openspec 與封存目錄）無 BRIEFAST_CONFIG 殘留，並逐處審查改寫後敘述正確。
- [x] 2.2 整體驗證：執行 go test ./... 全綠；本機不經容器啟動確認仍讀得到 repo 內的 syralit.toml 並綁 0.0.0.0:8600；執行 spectra validate bake-config-into-image 確認通過。
