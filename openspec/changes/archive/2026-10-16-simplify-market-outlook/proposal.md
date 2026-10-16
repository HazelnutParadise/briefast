## Why

每日報告常無法支持盤中路徑，卻仍被要求填寫三階段走法，形成重複的無法判斷文字。保留有依據的開盤與收盤判斷即可。

## What Changes

- 新報告只輸出大盤方向與判斷依據，開盤與收盤分別依證據判斷。
- 省略路徑文字與三階段圖，不再要求推測盤中轉折。
- 保留舊資料及 API 相容性，不新增欄位。

## Non-Goals

不修改歷史報告、API、資料庫、正式站或其他預測。

## Capabilities

### Modified Capabilities

- `daily-brief-skill`: 開盤與收盤獨立判斷，新報告省略路徑欄位。

## Impact

- Modified: `skills/daily-brief/SKILL.md`, `README.md`, `DESIGN.md`, `internal/report/skill_example_test.go`, `internal/site/site_test.go`
- Affected specs: `daily-brief-skill`
- 無 API、依賴或資料遷移變更。
