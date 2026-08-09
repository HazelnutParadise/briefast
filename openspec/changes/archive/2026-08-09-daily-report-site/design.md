## Context

Repo 目前只有文件與兩個 Syralit skills，沒有任何程式。規格已在討論中收斂：首頁即當日報告（盤前總覽 → 今日觀察 → 個股多空判斷 → 產業新聞摘要 → 個股新聞詳情 → 免責聲明），歷史頁逐日回看；報告由 Cowork agent 每個交易日開盤前產生並推送；需要後台管理 API key（可重複檢視明文）與報告更新紀錄；以 docker compose 部署。版型結構有使用者確認過的 HTML 示意稿（v5），視覺識別已定案為 FT 系規線報紙版式並寫成專案根目錄的 DESIGN.md（含淺深雙主題 token、字型、版面規則與 syralit.toml 對映），定案視覺稿在 docs/design/design-demos/benchmark.html。

## Goals / Non-Goals

**Goals:**

- 讀者開站即看到最新報告，版面固定不受 agent 輸出品質影響。
- 一次 POST 完成整份報告更新，伺服器負責驗證、持久化、即時推播。
- 後台可完整管理 API key 生命週期並追蹤每次推送。
- 單一 docker compose 指令可部署，資料放 volume。
- Cowork 每日流程有明確、可獨立執行的 skill 與輸出契約。

**Non-Goals:**

- 不整合即時行情或報價 API，所有內容皆來自新聞判讀。
- 不做判斷績效回測或戰績頁（未來可另立 change）。
- 不做讀者帳號、訂閱、RSS、多語言。
- 不使用 Artifact Canvas 與 Artifact DSL。
- HTTPS 憑證與網域由部署環境的反向代理處理，不在本 change 內。
- Cowork 端的排程機制（何時觸發每日流程）由 Cowork 自身設定，本 change 只提供 skill。

## Decisions

### 版面寫死於程式，agent 只送內容

Syralit 原生元件在 Go 內渲染五個區塊與樣式（多空卡片、紅漲綠跌標籤）。曾考慮 Artifact Canvas（單一或多 store），但版面固定後 Canvas 的動態 UI 能力沒有用處，且其 DSL 元件集做不出客製樣式；agent 送壞資料時也不應能弄壞版面。

### 單一報告 API 全量替換

POST /api/report 一次接收整份報告 JSON。同一日期重複推送即覆蓋（最後一次為權威版本）。曾考慮分區塊多次更新，但會出現新舊區塊混合的中間狀態，且需要 revision 協商；全量替換天然原子，agent 端也只需組一份 JSON。

### SQLite 持久化

單一 SQLite 資料庫 data/briefast.db（driver 用 modernc.org/sqlite，純 Go 免 CGO，docker 多階段建置不需 gcc；啟用 WAL 與 busy_timeout）。三張表：

- reports：date 為主鍵（YYYY-MM-DD）、headline、payload（整份報告 JSON 原文）、generated_at、created_at、updated_at。同日重推為 upsert。
- api_keys：id、name、token（明文、UNIQUE）、created_at、revoked_at（NULL＝有效；撤銷即軟刪除，不做硬刪）。
- update_log：id、at、api_key_id（FK → api_keys，ON DELETE RESTRICT，建索引）、key_name（寫入當下 snapshot，不靠 join 回查）、report_date、action（ingest_ok／ingest_rejected_auth／ingest_rejected_schema）、detail。拒絕的請求也記錄，不靜默。

Schema 以 migration 檔管理（internal/store/schema.sql 起，schema_migrations 表記版本，啟動時自動套用；migration 檔進 git、視為不可變）。開發與部署天然分離：dev 用本機 ./data/，prod 用 docker volume，兩個獨立資料庫檔案。曾考慮檔案系統（每日 JSON 檔），改用 SQLite 是使用者決策，同時換得跨日查詢、之後做績效回測與統計的擴充性。本專案為單機單寫者、無多使用者權限面，全表觸發器式 audit log 與 request log 依比例簡化為 update_log 一張表，此為記錄在案的取捨。

### 明文 API key 與後台管理

使用者明確要求已建立的 key 能重複檢視完整內容，因此 key 以明文存 data/keys.json（一般雜湊做法無法重看）。風險由後台登入（管理員密碼）與 volume 檔案權限承擔。撤銷即時生效：API 每次請求都重讀有效 key 清單。

### sy.Shared 即時推播

伺服器收到新報告後更新一個 sy.Shared 版本號，所有在線 session 觸發 rerun 重新載入最新報告。曾考慮輪詢（sy.RunEvery），但 Shared 推播即時且成本更低。

### 自訂 mux 掛載 API 與 Syralit app

main.go 建立自己的 http mux：/api/report 走自訂 handler，其餘路徑交給 sy.Handler 嵌入的 Syralit app（多頁：首頁、歷史、admin）。實作時需驗證 sy.AddPage 與 sy.Handler 併用的行為，若不相容則改用 sy.Navigation 在單一 app 函式內分頁。

### Cowork skill 輸出契約

skills/daily-brief/SKILL.md（專案根目錄 skills/，與 syralit repo 發佈 skills 的慣例一致）定義每日流程：蒐集新聞 → 判讀多空 → 組報告 JSON → POST。skill 內含完整 JSON schema 範例與 curl 指令模板；網站 URL 與 API key 由環境變數（BRIEFAST_URL、BRIEFAST_API_KEY）提供，不寫死在 skill 裡。

## Implementation Contract

**報告 JSON schema（POST /api/report 的 request body）：**

```json
{
  "date": "2026-08-07",
  "headline": "美股收紅開盤氣氛偏多，看漲：台積電、緯創",
  "overview_md": "盤前總覽 markdown",
  "watch_md": "今日觀察 markdown 條列",
  "calls": {
    "short_bull": [{"symbol": "2330", "name": "台積電", "reason": "一句理由"}],
    "short_bear": [],
    "long_bull": [],
    "long_bear": []
  },
  "industries": [{"name": "半導體", "summary_md": "該產業要聞 markdown 條列"}],
  "stock_news": [
    {
      "symbol": "2330", "name": "台積電",
      "call": "short_bull",
      "summary_md": "摘要段落，可含列點",
      "sources": [{"title": "新聞標題", "url": "https://..."}]
    }
  ],
  "generated_at": "2026-08-07T07:50:00+08:00"
}
```

- `call` 允許值：short_bull、short_bear、long_bull、long_bear、none（none = 有大消息但無明確判斷，前台不顯示標籤）。
- 驗證規則：date 必為 YYYY-MM-DD；calls 四清單中每一檔（symbol）必須在 stock_news 有對應條目，違反即拒收；headline、overview_md、watch_md 不得為空；sources 每條需有 url。
- 驗證失敗回 400，body 帶 `{"ok": false, "errors": ["每條錯誤描述"]}`，整份不落盤。

**API 行為：**

- 驗證：HTTP header Authorization: Bearer <key>，比對 data/keys.json 中未撤銷的 key；無效或已撤銷回 401。
- 成功：upsert 進 reports 表（以 date 為鍵）、插入一列 update_log（時間、key 名稱 snapshot、報告日期、action=ingest_ok）、更新 sy.Shared 版本號，回 200 `{"ok": true, "date": "..."}`。
- 拒絕也留痕：401 記 action=ingest_rejected_auth、400 記 action=ingest_rejected_schema（detail 帶錯誤摘要），報告內容不入庫。

**前台行為：**

- 首頁渲染 reports 表中日期最新的一份，masthead 標示報告日期與 generated_at；資料庫無任何報告時顯示「尚無報告」空狀態，不得當機。
- 歷史頁：日期新到舊列出（每頁 10 筆，headline 為摘要），點日期以與首頁相同的渲染函式顯示該日報告。
- 免責聲明頁尾寫死於程式，每頁皆顯示。

**後台行為（/admin）：**

- 未登入只見登入表單；管理員密碼來自 syralit.toml secrets 或環境變數 BRIEFAST_ADMIN_PASSWORD。
- Key 管理：建立（輸入名稱，產生隨機 token）、列表（名稱、建立時間、完整 token 隨時可見）、撤銷（標記 revoked_at，立即生效）。
- 更新紀錄：由新到舊顯示 update_log 內容（含被拒絕的嘗試）。

**驗收方式：**

- go test ./... 全綠：API 用 httptest 覆蓋（合法推送入庫＋200、壞 schema 400 且不入庫、無效與已撤銷 key 401、拒絕留痕）；前台與後台用 sy.AppTest 覆蓋（空狀態、有報告渲染、歷史列表、登入閘門、key 建立與撤銷）。
- docker compose up -d --build 後，瀏覽器開站可見報告頁，重啟容器後歷史與 key 仍在（volume 生效）。

**範圍界線：**

- In scope：上述五個 capability 的全部行為、測試、部署設定、AGENTS.md 與 README 更新。
- Out of scope：Non-Goals 列出的所有項目、視覺細節微調（結構以示意稿 v5 為準、視覺 token 與版面語言以 DESIGN.md 為準，像素級調整不擋驗收）。

## Risks / Trade-offs

- [明文 key 落盤，拿到 data/ 即拿到所有 key] → 後台密碼強制設定（未設定則 admin 頁拒絕服務）、文件明示 volume 權限建議；接受此風險為使用者明確決策。
- [agent 判讀屬投資意見，公開網站有誤導疑慮] → 免責聲明頁尾寫死每頁顯示；判斷必附理由與新聞來源。
- [sy.AddPage 與 sy.Handler 併用行為未經驗證] → 實作首日先做 spike，不相容即改 sy.Navigation 單 app 分頁，兩方案 UI 相同。
- [同日覆蓋無版本歷程，agent 推錯即遺失前一版] → 接受（agent 可重推正確版）；updates.jsonl 至少留下推送軌跡。
- [非交易日不推送，首頁日期停留在上個交易日] → 屬預期行為，masthead 明確標示報告日期避免誤會。

## Migration Plan

全新部署，無既有資料可遷移。上線步驟：主機放 docker-compose.yml 與 syralit.toml（自 example 複製填入密碼）→ docker compose up -d --build → 後台建立第一把 key → 提供給 Cowork 環境變數。回滾即停容器，data/ 不受影響。

## Open Questions

（無——規格已在討論階段收斂完畢）
