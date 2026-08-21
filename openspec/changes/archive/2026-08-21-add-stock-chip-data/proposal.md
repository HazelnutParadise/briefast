## Why

目前個股多空判讀只依據新聞與昨收價，缺少籌碼面佐證：三大法人買賣超與融資融券變化是市場已公開、可查證的資金動向，能幫助判斷利多利空是否已被大戶反映。報告頁也缺少任何數據視覺化，讀者無法快速看出法人動向。判讀規範中「籌碼資金面只作佐證加分」目前有規則卻沒有資料來源支撐。

## What Changes

- 每日晨報蒐集流程新增籌碼面資料抓取：上市（證交所 T86 三大法人買賣超、MI_MARGN 融資融券）與上櫃（櫃買 openapi 三大法人買賣明細、融資融券餘額），各端點已驗證可用且回應含資料日期欄位。
- 報告 JSON 的 stock_news 條目新增選填 chips 區塊，內含資料日期、外資／投信／自營商買賣超股數、三大法人合計買賣超股數、融資餘額增減與融券餘額增減；API 驗證接受並檢核該區塊，舊報告與缺漏該欄位的新報告照常接受與渲染。
- 個股新聞詳情版面為帶有 chips 資料的條目渲染籌碼面視覺化：以水平長條呈現三類法人買賣超的相對量（買超紅、賣超綠，遵循台股紅漲綠跌），並列出融資融券增減數字與資料日期。
- 判讀規範明定籌碼面引用紀律：籌碼數據只能作為既有新聞判斷的佐證，不得單獨構成 call；引用須標明資料日期；抓取失敗時照常判讀與發布，chips 整塊省略。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `daily-brief-skill`: 新增籌碼面參考資料的抓取、日期驗證、節流歸屬與判讀引用紀律，以及 chips 區塊的組稿規則。
- `report-ingest-api`: 報告 schema 新增 stock_news 選填 chips 區塊與其驗證規則。
- `report-viewing`: 個股新聞詳情條目新增籌碼面視覺化區塊，無資料時整塊省略。

## Impact

- Affected specs: `daily-brief-skill`、`report-ingest-api`、`report-viewing`
- Affected code:
  - Modified: `skills/daily-brief/SKILL.md`、`internal/report/schema.go`、`internal/report/schema_test.go`、`internal/site/site.go`、`internal/site/styles.go`、`internal/site/site_test.go`、`internal/api/report_test.go`
  - New: （無）
  - Removed: （無）
- 不需資料庫 migration：報告以整包 JSON 存於 reports.payload，schema 擴充不動資料表。
