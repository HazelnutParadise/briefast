## 1. 專案骨架與持久層

- [x] 1.1 建立 Go module 與 Syralit 進入點（go.mod、main.go），並依 design「自訂 mux 掛載 API 與 Syralit app」完成 spike：確認 sy.AddPage 多頁與 sy.Handler 併用可行，不可行則改用 sy.Navigation 單 app 分頁。驗證：go run . 可啟動，瀏覽器開站看到空首頁骨架。
- [x] 1.2 實作報告 schema 與驗證（internal/report/schema.go），交付 Report schema validation 與 Calls and stock news consistency 的全部規則（date 格式、必填欄位、call 允許值、calls 與 stock_news 對應）。驗證：schema 單元測試以表格覆蓋合法與各類非法 payload，go test 全綠。
- [x] 1.3 依 design「SQLite 持久化」實作資料層（internal/store/store.go、internal/store/schema.sql）：modernc.org/sqlite 驅動、WAL 與 busy_timeout、schema_migrations 啟動時自動套用 migration；建 reports（date 主鍵、payload、標準時間欄位）、api_keys（明文 token、revoked_at 軟刪除）、update_log（api_key_id FK 加索引、key_name snapshot、action）三表；查詢：最新報告、日期列表分頁、同日 upsert。驗證：store 單元測試涵蓋空庫、寫入、同日覆寫、日期排序、撤銷後查詢，go test 全綠。

## 2. 報告接收 API

- [x] 2.1 依 design「單一報告 API 全量替換」實作 POST /api/report（internal/api/report.go）：Bearer key authentication（缺 key、未知 key、已撤銷 key 一律 401 且不入庫，並記 update_log action=ingest_rejected_auth）、驗證失敗回 400 帶完整 errors 陣列且整份不入庫（記 action=ingest_rejected_schema）、Successful ingestion（reports upsert＋update_log 一列＋回 200 ok，同日重推覆蓋）。驗證：httptest 覆蓋 401、400、200、同日覆寫與拒絕留痕五類案例，go test 全綠。
- [x] 2.2 依 design「sy.Shared 即時推播」實作 Live update push：ingestion 成功後更新 sy.Shared 版本號，所有在線 session 自動 rerun 顯示新報告。驗證：AppTest 模擬版本號變更後，首頁渲染內容切換為新報告。

## 3. 前台頁面

- [x] 3.1 依 design「版面寫死於程式，agent 只送內容」實作首頁：Homepage shows the latest report（渲染最新日期報告、masthead 標日期與 generated_at、無報告時顯示空狀態不當機）與 Fixed report layout（五區塊固定順序，與 payload 欄位順序無關）。驗證：AppTest 空狀態與有報告兩案例斷言區塊順序與日期。
- [x] 3.2 實作 Stock calls display：短多、短空、長多、長空四張卡片，紅漲綠跌配色，每檔顯示股名、代號、一句理由。驗證：AppTest 以含四類判斷的報告斷言各卡內容。
- [x] 3.3 實作 Stock news detail display：每檔一條（股名代號、call 標籤、摘要 markdown、來源連結），call 為 none 的條目不顯示標籤。驗證：AppTest 含 none 條目案例，斷言標籤有無。
- [x] 3.4 實作產業新聞摘要卡片與 History browsing：歷史頁日期新到舊、每頁 10 筆、每列顯示 headline，點日期以與首頁相同的渲染函式顯示該日報告。驗證：AppTest 建 25 份報告斷言排序與分頁，並斷言點選歷史日期後渲染該日內容。
- [x] 3.5 實作 Disclaimer footer：免責文字寫死於程式，首頁與歷史檢視皆顯示。驗證：AppTest 斷言兩種頁面都含免責文字。

## 4. 後台管理

- [x] 4.1 依 design「明文 API key 與後台管理」實作 Admin login gate：/admin 未登入僅顯示登入表單；管理員密碼取自 syralit.toml secrets 或 BRIEFAST_ADMIN_PASSWORD；未設定密碼時顯示設定錯誤並拒絕登入。驗證：AppTest 未登入、密碼未設定、登入成功三案例。
- [x] 4.2 實作 API key creation 與 API keys remain fully viewable：輸入名稱建立 key、產生隨機 token、列表隨時顯示完整明文 token（含日後重新登入）。驗證：AppTest 建立後模擬重新進入頁面仍能讀到完整 token。
- [x] 4.3 實作 API key revocation：撤銷寫入 revoked_at、列表標示已撤銷、API 端立即 401。驗證：AppTest 撤銷後列表狀態改變，httptest 以該 key 請求得 401。
- [x] 4.4 實作 Report update log view：後台顯示 update_log 表內容（含被拒絕的嘗試），新到舊，每筆含時間、key 名稱、報告日期與動作。驗證：AppTest 推送兩筆報告後斷言列表順序與欄位。

## 5. 部署與文件

- [x] 5.1 實作 Docker compose deployment：多階段 Dockerfile 與 docker-compose.yml（data/ 掛 volume、port 可設定）。驗證：docker compose up -d --build 後開站看到報告頁；重建並重啟容器後歷史報告與 key 仍在。
- [x] 5.2 實作 Configuration injection：提供 syralit.toml.example 註解所有設定；新增 .gitignore 排除 syralit.toml 與 data/。驗證：複製 example 填入密碼可啟動；git status 確認 syralit.toml 與 data/ 不被追蹤。
- [x] 5.3 更新 AGENTS.md 與 README.md 反映定案架構（報告 JSON schema、API、前後台、部署步驟、skill 用法），移除已否決的 Artifact Canvas 與 Artifact DSL 約定。驗證：逐段審閱兩份文件與實際程式行為一致，無過時敘述。

## 6. Cowork 每日流程 skill

- [x] 6.1 依 design「Cowork skill 輸出契約」建立 skills/daily-brief/SKILL.md（專案根目錄 skills/），滿足 Daily workflow skill exists in the repo（蒐集新聞 → 判讀多空 → 組報告 JSON → POST 四步驟）與 Skill output contract（內嵌完整 schema 與填好的範例、每檔判斷必附一句理由與至少一個來源、方向不明者以 call none 進 stock_news）。驗證：內容審閱對照 report-ingest-api spec 的 schema 與驗證規則逐項一致。
- [x] 6.2 在 skill 中落實 Endpoint configuration via environment 與 Failure reporting：讀 BRIEFAST_URL 與 BRIEFAST_API_KEY、缺任一即停止並回報；非 200 視為失敗並附回應內容，5xx 重試一次，400 依 errors 修正或原文回報。驗證：以本機測試伺服器實際演練一次成功 POST 與一次 400 流程。
