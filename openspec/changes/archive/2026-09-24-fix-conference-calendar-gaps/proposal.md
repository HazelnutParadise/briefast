## Why

網站漏登已公告的公司自辦法說會：星亞（7753，2026-10-14）的第 12 款公告已在每日蒐集檔，更新後清單卻未收錄。僅從近期重大訊息蒐集、且以「參加」字串排除場次，無法確保自辦場次完整。

## What Changes

- 每次處理法說會時，以上市與上櫃官方月曆核對未來場次，將漏登的自辦場次納入會前報告流程。
- 將排除條件改為辨識受邀或券商主辦的語意，避免把「參加方式」誤判為受邀。
- 月曆與重大訊息證據不足時，記錄待查，不憑地點或模糊文字猜測。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `daily-brief-skill`: 法說會自辦分類與官方月曆補漏。

## Impact

- Affected code:
  - Modified: `skills/daily-brief/SKILL.md`
  - Modified: `openspec/specs/daily-brief-skill/spec.md`
- No site rendering, database schema, or API change.
