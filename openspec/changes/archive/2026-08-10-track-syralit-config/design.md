## Context

Briefast 的執行期設定放在 `syralit.toml`，內容包含站台標題、主機與埠、主題色票、i18n 字串與一個 `[secrets]` 區段。目前 repo 只放 `syralit.toml.example`，真檔被 `.gitignore` 排除，README 教使用者複製一份填密碼，compose 的預設掛載來源則是那個範本檔。實體的 `syralit.toml` 目前並不存在於工作區——也就是說「範本加本機副本」這套流程實際上沒有被使用，大家跑的都是範本本身。本次把設定檔改為正式進版控。

## Goals / Non-Goals

**Goals:**

- 設定檔在 repo 裡只有一份，不再有範本與真檔分歧的可能。
- 新加入的人 clone 後不需要任何複製步驟就能啟動。
- 機密不因設定檔進版控而外洩。

**Non-Goals:**

- 不改 Go 的設定讀取邏輯與 `BRIEFAST_CONFIG` 覆寫行為。
- 不改 `.env` 的排除規則。
- 不引入秘密管理工具或加密設定檔。
- 不調整主題色票與 i18n 字串的內容。

## Decisions

### 設定檔進版控

改為單一份進版控的 `syralit.toml`，捨棄「範本加本機副本」。理由是本專案的設定內容除了密碼之外全部是跨環境共用的非機密值，分成兩份只換來分歧風險與一個多餘的啟動步驟，沒有換到隔離價值。既有的 `BRIEFAST_CONFIG` 環境變數仍可指向自訂設定檔，需要環境差異時走那條路，不需要靠 gitignore 製造差異。

### 機密與部署變數集中在 .env

設定檔一旦進版控，`[secrets]` 就成為誘導犯錯的欄位——只要有人照著填一次密碼並提交，機密就永久留在 git 歷史。因此整個區段從進版控的設定檔移除，而不是留著註解：留註解仍是把「可以填在這裡」的暗示留在檔案裡。`sy.Secrets` 找不到設定檔區段時會回退到同名環境變數，所以移除不影響取值。

機密改由 repo 根目錄的 `.env` 提供，該檔同時承擔部署變數的角色。這個組合有一個現成的便利：Docker Compose 原生會讀專案目錄的 `.env` 做變數替換，因此 `BRIEFAST_CONFIG`（設定檔掛載來源）、`BRIEFAST_PORT` 與 `BRIEFAST_ADMIN_PASSWORD` 寫在同一份 `.env` 就同時餵給容器與 compose 本身，不需要 env_file 設定或額外參數。代價是本機直接跑 Go 時不會自動載入——程式不解析 `.env`，README 因此附上先把 `.env` 載入環境再啟動的做法。捨棄在 Go 端引入 dotenv 套件的理由是：那會讓程式多一個只在開發期有用的依賴，而 shell 一行就能解決。

`.env` 無法進版控，鍵名的文件因此需要另一個載體，repo 根目錄新增 `.env.example` 記錄所有鍵與用途。這與 `syralit.toml` 取消範本並不矛盾：`syralit.toml` 取消範本是因為真檔本身已進版控，`.env` 的真檔永遠不能進版控，範本是唯一的文件位置。

### compose 預設掛載

預設掛載來源由範本檔改為 `syralit.toml`。現況「名為 example 的檔案是正式部署預設值」本身就是設計氣味；改名後預設值指向唯一那份設定。`BRIEFAST_CONFIG` 的覆寫行為不變，只是它的值現在通常來自 `.env`，由 Compose 自動替換。

## Implementation Contract

- **版控狀態**：`syralit.toml` 存在於 repo 且被 git 追蹤；`syralit.toml.example` 不再存在；`.gitignore` 不再排除 `syralit.toml`，但仍排除 `.env`。
- **內容約束**：進版控的 `syralit.toml` 不含 `[secrets]` 區段；`BRIEFAST_ADMIN_PASSWORD` 不出現在任何進版控檔案的實值位置。repo 根目錄的 `.env.example` 列出所有鍵名與用途且無真值。
- **執行行為**：不帶 `BRIEFAST_CONFIG` 啟動時讀取 repo 內的 `syralit.toml`；帶 `BRIEFAST_CONFIG` 時讀取指定檔案，行為與現況相同。Compose 啟動時自動從 `.env` 取得 `BRIEFAST_CONFIG`、`BRIEFAST_PORT` 與 `BRIEFAST_ADMIN_PASSWORD`，不需額外參數。管理員密碼由環境提供，未設定時公開頁與 API 照常服務、後台拒絕登入並顯示設定錯誤——此行為不變。
- **文件一致性**：README 的本機開發與 Docker 部署兩節不再出現複製範本的步驟；AGENTS.md 的設定約定與新規則一致。
- **驗收標準**：以 git ls-files 確認 `syralit.toml` 與 `.env.example` 被追蹤且 `syralit.toml.example` 不存在；以 git check-ignore 確認 `syralit.toml` 與 `.env.example` 未被忽略而 `.env` 仍被忽略；以 grep 確認全 repo 無殘留的 `syralit.toml.example` 引用、且進版控檔案無 `[secrets]` 區段；以 docker compose config 確認 `.env` 的值有被替換進掛載來源與埠；執行 go test 全綠並實際啟動確認讀得到設定。
- **範圍邊界**：in scope 為 `syralit.toml`、`.env.example`、`.gitignore`、`docker-compose.yml`、README.md、AGENTS.md 與 `deployment` spec；out of scope 為 Go 設定讀取程式（不加 dotenv 解析）、主題與 i18n 內容。

## Risks / Trade-offs

- [有人日後在進版控的設定檔填入密碼並提交，機密永久留在 git 歷史] → 整個 `[secrets]` 區段從進版控檔案移除，檔案裡沒有可填的位置；spec 以需求形式固定此約束，讓後續變更必須正面推翻它才能改。
- [本機直接跑 Go 時 `.env` 不會自動生效，密碼看似沒設定] → README 在本機開發一節明確附上先載入 `.env` 再啟動的做法；未載入時的表現是後台顯示設定錯誤，訊息本身已指向 `BRIEFAST_ADMIN_PASSWORD`。
- [不同環境需要不同設定時，單一檔案不夠用] → `BRIEFAST_CONFIG` 覆寫機制原封保留，環境差異走覆寫而非改進版控的檔案。
- [compose 預設值改變，既有以預設值部署的環境行為變動] → 掛載目標路徑與檔案內容都不變，只是來源檔名由範本改為正式檔，容器內看到的設定完全相同。
