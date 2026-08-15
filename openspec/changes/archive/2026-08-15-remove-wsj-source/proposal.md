## Summary

移除 WSJ Markets 來源，批次 5 外電只保留 CNBC 頭條。

## Motivation

實測 WSJ 文章頁有付費牆（回 401），只能取得 RSS 內附的標題與一句摘要。摘要層級的資訊不足以支撐事件判讀與跨語言去重比對，留著只是增加抓取與比對成本，還可能因摘要語焉不詳造成 overview 誤讀。

## Proposed Solution

從 spec 與 SKILL.md 全面移除 WSJ Markets：釘選來源集刪除該端點與其摘要例外條款；批次 5 改為只涵蓋 CNBC 頭條 feed 並回報其計數；外電範圍與跨語言去重規則的行為不變，適用對象改為單一 CNBC 來源；去重章節刪除 WSJ 摘要比對的特例說明。

## Non-Goals

- 不尋找 WSJ 的替代外電來源；CNBC 可正常抓全文，先以單一外電來源運行，涵蓋不足再另議。
- 不改動外電的用途限制與去重行為本身，只縮減來源集。
- 不改動批次 1 至 4 與昨收價資料的任何規則。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `daily-brief-skill`: 釘選來源集移除 WSJ Markets；批次 5 涵蓋範圍改為僅 CNBC；外電範圍與去重 requirement 的來源指涉同步更新。

## Impact

- Affected specs: `daily-brief-skill`
- Affected code:
  - Modified: skills/daily-brief/SKILL.md
  - New: （無）
  - Removed: （無）
