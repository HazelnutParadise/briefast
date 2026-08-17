## Why

昨收價蒐集是一次取回全市場的大 JSON（上市約 1,300 筆、上櫃約 10,000 筆），但 SKILL.md 只規定抓取、日期驗證與計數回報，沒有規定讀取方式。執行的 agent 可能把整份回應讀進對話 context，浪費大量 context 且對判讀毫無幫助——判讀時只需要新聞提及的那幾十檔。

## What Changes

- 在 SKILL.md「昨收價參考資料」一節明定讀取紀律：兩個回應抓回後存成工作資料夾檔案，判讀時只按新聞提及的代號查詢（如 jq、grep），不得把整份回應讀進對話 context。
- 日期驗證與計數回報同樣以查詢方式取得（讀 `Date` 欄位與筆數即可），不需要整份讀入。
- 同步修改 `daily-brief-skill` spec 的 Previous-close reference data requirement，加入上述讀取紀律與對應 scenario。

## Non-Goals

- 不改變抓取端點、節流規則、日期驗證邏輯與失敗分流，維持既有 requirement 的其餘內容。
- 不新增歷史價格累積與技術分析（本次討論已決定維持現狀，沿用 2026-08-15-previous-close-reference 的 Non-Goal）。
- 不規定檔名與查詢工具的具體實作，只規定「存檔後按代號查詢、不整份讀入」的行為約束。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `daily-brief-skill`: Previous-close reference data requirement 加入讀取紀律——回應存檔後按代號查詢，禁止整份讀入對話 context。

## Impact

- Affected specs: `daily-brief-skill`
- Affected code:
  - Modified: skills/daily-brief/SKILL.md
  - New: （無）
  - Removed: （無）
