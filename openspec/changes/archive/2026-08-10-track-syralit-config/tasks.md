## 1. 設定檔進版控

- [x] 1.1 依 Configuration injection 與 design 的「設定檔進版控」決策，把 syralit.toml.example 更名為 syralit.toml 納入版控，並從 .gitignore 移除 syralit.toml 排除條目（保留 .env 排除）。驗證方式：git ls-files 顯示 syralit.toml 被追蹤且 syralit.toml.example 不存在，git check-ignore 對 syralit.toml 無輸出、對 .env 仍有輸出。
- [x] 1.2 依 design 的「機密與部署變數集中在 .env」決策，從進版控的 syralit.toml 移除整個 secrets 區段，只留非機密設定；確認移除後 sy.Secrets 仍能由環境變數取得管理員密碼。驗證方式：grep 確認檔內無 secrets 區段，並以設定環境變數的方式啟動確認後台可登入。

## 2. .env 作為機密與部署變數來源

- [x] 2.1 依 design 的「機密與部署變數集中在 .env」決策，於 repo 根目錄新增 .env.example，列出 BRIEFAST_ADMIN_PASSWORD、BRIEFAST_CONFIG、BRIEFAST_PORT、BRIEFAST_DB_PATH 的鍵名、用途與佔位值，無任何真值，並註明本機直接跑 Go 時需先把 .env 載入環境。驗證方式：cat 檔案確認內容，git check-ignore 確認 .env.example 未被忽略而 .env 仍被忽略。
- [x] 2.2 依 design 的「compose 預設掛載」決策，把 docker-compose.yml 的設定檔預設掛載來源由範本改為 syralit.toml，容器內掛載路徑與唯讀模式不變。驗證方式：在 .env 填入測試值後執行 docker compose config，確認展開後的掛載來源與埠取自 .env，未設定時退回預設值。

## 3. 文件與驗證

- [x] 3.1 更新 README 的本機開發與 Docker 部署兩節：移除複製 syralit.toml 範本的步驟，改為複製 .env.example 為 .env 並填入密碼；本機開發一節附上先把 .env 載入環境再啟動的做法。同步更新 AGENTS.md 的設定約定，使其與「設定檔進版控、機密與部署變數集中在 .env」一致。驗證方式為內容審查，並以 grep 確認全 repo（排除封存目錄）無殘留的 syralit.toml.example 引用。
- [x] 3.2 端到端驗證：執行 go test ./... 全綠；以載入 .env 的方式啟動確認讀得到 repo 內的 syralit.toml、站台可服務、後台可用密碼登入；未提供密碼時後台顯示設定錯誤而公開頁正常；執行 spectra validate track-syralit-config 確認通過。
