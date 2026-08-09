## Why

daily-brief skill 目前對 `industries` 區塊完全沒有撰寫指引，產業涵蓋隨當天抓到的新聞浮動，重要產業可能整天消失；個股判讀也沒有分析視角與基本面引用紀律，執行模型可能憑訓練記憶寫出過時的營收、EPS 等數字。需要把判讀角色、資料時效規則與固定產業分類寫進 skill，讓每日報告的深度與涵蓋穩定一致。另外，觀察重點（每個條目的前瞻追蹤點）若值得存在，就應成為版面上的固定區塊與 API 可驗證的欄位，而不是藏在摘要尾段的文字慣例。

報告發布的傳輸面已是原子的——API 只收完整 JSON，驗證失敗整份不落庫，upsert 與 log 同一 transaction——但內容完整性沒有任何規範：新聞來源掛掉時 agent 仍能組出欄位全部合法卻內容殘缺的報告並成功發布，且同日全量覆寫會讓內容較少的重跑靜默蓋掉先前完整的版本，事後無從察覺。

## What Changes

- `skills/daily-brief/SKILL.md` 新增判讀角色指引：產業區塊以產業分析師視角撰寫，個股區塊以個股分析師視角撰寫。
- 新增基本面引用紀律：基本面數字只能來自當次執行抓取且附資料期間的公開來源（如 TWSE OpenAPI 月營收），禁止憑模型記憶填寫；抓不到就不寫基本面，不硬湊。
- `industries` 改為固定四分類「科技、金融、傳產、房地產」，定義為封閉分類（任何台股相關產業新聞必落入其一），區塊內以粗體子題帶當日動態內容；當天查證後確無新聞的分類可省略，但省略必須是查證後的決定。
- **BREAKING** report JSON 的 `industries[]` 與 `stock_news[]` 各新增 `watch_md` 欄位（觀察重點，每項 1–2 條可事後驗證的前瞻追蹤點）。API 驗證新報告兩處皆必填非空白，缺漏回 400；既有產出端只有 daily-brief skill，於同一 change 內同步更新。
- `internal/site` 於產業與個股條目渲染獨立「觀察重點」子區塊；`watch_md` 空缺（欄位存在前的歷史報告）時不渲染、無佔位。
- 明確定義涵蓋範圍：個股判讀限台股，國際市場與個股動態只進 `overview_md`。
- 新增蒐集完整性把關：蒐集階段逐一記錄每個來源的成敗並對失敗來源重試；主來源鉅亨網重試後仍失敗就停止不發報告，次要來源失敗則照常發布但在執行回報中明列缺哪些來源。
- 新增覆寫防護：同一日期重跑時，若新報告的產業數或個股數少於已發布版本，必須停下來說明差異原因並取得指示，不得直接覆寫。
- 新增唯讀端點 `GET /api/report/{date}`，以與 ingest 相同的 Bearer key 驗證回傳該日完整報告 JSON，讓覆寫比對與既有「從上一份報告的 generated_at 起算時間窗口」規則實際可執行。
- 同步更新 `daily-brief-skill`、`report-ingest-api`、`report-viewing` 三份 spec 對應需求。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `daily-brief-skill`: 新增判讀角色、基本面時效紀律、固定四分類、watch_md 觀察重點與台股範圍的 skill 內容要求。
- `report-ingest-api`: Report schema validation 擴充——`industries[].watch_md` 與 `stock_news[].watch_md` 必填非空白。
- `report-viewing`: 新增觀察重點區塊渲染需求，含歷史報告空缺相容。

## Impact

- Affected specs: `daily-brief-skill`, `report-ingest-api`, `report-viewing`
- Affected code:
  - New: （無）
  - Modified: skills/daily-brief/SKILL.md, internal/report/schema.go, internal/report/schema_test.go, internal/api/report.go, internal/api/report_test.go, internal/site/site.go, internal/site/site_test.go, internal/site/styles.go, main.go, main_test.go
  - Removed: （無）
- 資料庫結構不變：報告以完整 JSON 儲存，新欄位不需要 migration。
