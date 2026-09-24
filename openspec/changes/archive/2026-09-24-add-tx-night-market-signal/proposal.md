## Why

目前每日晨報已根據美股、新聞、籌碼和法說判斷大盤走勢，卻沒有使用盤前已結束的臺指期夜盤。以期交所及證交所歷史資料回測，夜盤方向對開盤相對昨收較有參考價值，因此需要納入開盤情境，同時限制它對全天路徑的影響。

## What Changes

- 每日晨報在判讀前取得期交所 TX 夜盤與前一交易日同契約日盤資料，核對日期、時段與合約後形成可追溯的盤前訊號。
- 將已核對的夜盤訊號用於開盤傾向；與費半、那斯達克及台股證據背離時說明取捨，不能單靠夜盤推論收盤方向或盤中轉折。
- 夜盤資料缺漏、尚未完整公布或無法核對時，繼續依其他已驗證資料產出晨報，並在執行回報記錄缺漏。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `daily-brief-skill`：每日市場走勢判讀加入已核對的 TX 夜盤開盤訊號、背離處理與缺漏回退。

## Impact

- Affected specs: `daily-brief-skill`
- Affected code:
  - New: `skills/daily-brief/scripts/tx_night.py`, `skills/daily-brief/scripts/tx_night_test.py`
  - Modified: `skills/daily-brief/SKILL.md`
- 既有報告 JSON 結構、網站版面、盤前總覽和歷史報告不變。
