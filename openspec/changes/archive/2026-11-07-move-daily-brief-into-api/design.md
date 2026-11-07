## Context

每日晨報由 Claude Cowork 排程執行。agent 目前靠三樣東西運作：

- `skills/daily-brief/SKILL.md`（約 820 行）：交易日檢查、五批新聞來源、判讀與撰寫規則、報告 JSON 範例、POST 與錯誤處理。
- `seen.py`：去重工具，紀錄存在 Cowork 工作區的 `seen.jsonl`。
- `tx_night.py`：讀期交所 TX 每日行情 CSV，算出夜盤相對前一交易日同契約日盤的漲跌。

網站已經有 `POST /api/report`、`GET /api/report/{date}`、法說會三個端點，全部用 Bearer key 驗證，401 和 400 都會寫進 `update_log`。agent 本來就要打這些端點，所以把說明和工具功能放到同一個伺服器上，就不需要另外安裝 skill。

## Goals / Non-Goals

**Goals:**

- agent 只拿到網址和 API key，打任何一個 `/api/` 端點都會看到提示，知道要先讀 `GET /api/guide`。
- 流程說明、驗證規則、去重和夜盤計算都隨網站一起部署，只有一份來源。
- 刪除 `skills/daily-brief/` 後，每日晨報與法說會流程的行為不變。

**Non-Goals:**

- 不做 MCP server。
- 交易日檢查、新聞蒐集、TWSE／TPEx／TAIFEX 下載仍由 agent 執行，伺服器不主動對外抓資料。
- 不匯入舊的 `seen.jsonl`。
- 不改報告與法說會 brief 的 JSON 結構，也不改公開頁面。

## Decisions

### 說明文件嵌入 Go 程式，以 Markdown 回傳

SKILL.md 的內容搬到 `internal/guide/daily-brief.md`，用 `go:embed` 嵌入，`GET /api/guide` 回傳 `text/markdown; charset=utf-8`。前置的 YAML frontmatter 刪掉。原本 `$DAILY_BRIEF_SKILL_DIR/scripts/...` 的命令改寫成 curl 呼叫新端點，報告讀取段落改成讀 `report` 欄位。

替代方案是把說明拆成多個章節端點，但 agent 一次讀完全文最簡單，約 35 KB 也不會太大，所以不拆。

### `guide` 需要 Bearer key，總覽與 OpenAPI 公開

說明裡有來源清單、判讀規則這類流程細節，不必公開，所以 `GET /api/guide` 跟其他端點一樣驗證 key，401 寫進 `update_log`（action `guide_rejected_auth`）。成功讀取不寫 log，比照報告讀取端點。

`GET /api` 和 `GET /api/openapi.json` 只描述端點形狀，不含機密，公開讀取。key 打錯的 agent 仍然能從總覽知道該怎麼呼叫。

### 每個回應都帶提示

`writeJSON` 統一在回應物件加上 `guide` 欄位，並設 `Link: </api/guide>; rel="help"` 標頭。提示文字固定為：

> 先用同一把 API key 呼叫 GET /api/guide，讀完每日晨報流程說明再呼叫其他端點。

`GET /api/guide` 本身回 Markdown，只設 `Link` 標頭。

`GET /api/report/{date}` 原本直接回傳報告物件，加欄位會讓報告 JSON 多出 `guide`。agent 如果把讀回的內容改一改再 POST，會被 `DisallowUnknownFields` 擋下。所以改成 `{"ok":true,"guide":"…","report":{…}}`。目前只有 daily-brief 流程使用這個端點，說明會同步改寫。

`/api/` 底下沒有對應路由的請求，由新的 `/api/` 處理器回 JSON 404 附提示，不再落到首頁的 Syralit app。

### 去重改存 SQLite

新增 `migration_v3.sql`：

- `seen_articles`：`source_id`、`url`、`url_hash`、`title`、`title_hash`、`content_hash`、`lead_hash`、`sketch`（JSON 陣列）、`source_domain`、`symbols`（JSON 陣列）、`event_key`、`published_at`、`decision`、`call`、`first_seen_at`。`source_id`、`url_hash`、`title_hash` 各建索引。
- `seen_events`：`event_key`（主鍵）、`note`、`first_seen_at`。

每次呼叫去重端點時，先刪除超過保留期限的紀錄（文章 90 天、事件 30 天）。請求與回應的欄位沿用 `seen.py` 的格式，例如 peek 回 `{"new":[…],"seen":[…]}`，seen 裡的每筆帶 `seen_by`。

比對規則照搬 `seen.py`：網址去掉追蹤參數與片段，文字做 NFKC 正規化並去除 HTML 標籤與署名，相似度用字元二連文的 bottom-64 MinHash 估計 Jaccard。`similar` 預設列出 0.25 以上的候選，說明文件仍規定 0.5 以上才算同一篇。雜湊函式從 blake2b 換成 SHA-256 前 8 bytes，因為舊資料不匯入，換雜湊不影響相容性，也不用加依賴。NFKC 需要 `golang.org/x/text/unicode/norm`，與 CP950 解碼共用同一個依賴。

去重端點的寫入不記 `update_log`，401 和 400 照專案慣例記錄。

### 夜盤由伺服器解析上傳的 CSV

`POST /api/tx-night` 收 multipart：`report_date`、`previous_date`、`report_csv`、`previous_csv`。伺服器用 `golang.org/x/text/encoding/traditionalchinese` 解碼 CP950，邏輯照搬 `tx_night.py`，成功回傳同樣欄位的 JSON。失敗時回 400，`errors` 寫明原因，不回傳部分結果。

替代方案是讓伺服器自己去期交所下載。這會讓網站多一條對外連線和失敗來源，也和其他資料由 agent 抓取的分工不一致，所以不採用。

### capability 名稱暫時保留

`daily-brief-skill` 的需求大多描述流程內容，換成 API 提供後仍然適用，所以這次只修改受影響的需求。archive 後把資料夾改名為 `daily-brief-workflow`，內文的 "skill" 改成 "guide"。

## Risks / Trade-offs

- [舊的已讀紀錄消失，部署後頭幾天可能重報延續中的事件] → 新聞窗口從上一份報告的 `generated_at` 算起，大部分舊文章本來就在窗口外。事件登記只保留 30 天，影響最多一個月。部署當天在執行回報裡註明去重紀錄重新開始。
- [Cowork 的 agent 沒讀說明就直接 POST] → 每個回應都有 `guide` 提示。400 的錯誤訊息也會指出缺少哪些欄位，agent 讀了提示就會回頭讀說明。
- [說明文件變長後 agent 漏讀段落] → 文件結構和現在的 SKILL.md 相同，只把腳本段落換成端點呼叫，不增加內容。
- [讀取端點格式改變] → 部署和說明更新在同一次發布，舊 skill 同時刪除，不會有新舊版本並存。
- [Big5 解碼器與 CP950 的少數字元不同] → TX CSV 的欄名與值（日期、契約、價量、`一般`、`盤後`）都在共通字集內。以期交所 2026-09-24 與 2026-09-23 的原始 CSV 做測試確認。

## Migration Plan

1. 部署新版網站，啟動時套用 `migration_v3.sql`。
2. 用 curl 確認 `GET /api`、`GET /api/guide`、未知 `/api/x` 的回應。
3. 把 Cowork 排程 prompt 換成 README 的一行指示，移除工作區的 daily-brief skill。
4. 隔天晨報執行後，檢查 `update_log` 有 `ingest_ok`、去重資料表有新紀錄。

回復方式：重新部署前一版映像，並把 Cowork 排程換回 skill。`migration_v3.sql` 只新增資料表，舊版程式不會讀取，可以留著。

## Open Questions

（無）
