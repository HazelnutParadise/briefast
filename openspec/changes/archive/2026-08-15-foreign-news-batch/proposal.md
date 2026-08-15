## Why

daily-brief 目前的國際新聞涵蓋完全依賴台灣媒體的轉譯報導（鉅亨 `wd_stock` 分類與各台媒的美股、Fed 報導）。台媒轉譯有一到三小時的時差，且中型國際事件（出口管制細節、供應鏈訊息、海外客戶動態）可能不被台媒跟進，造成盤前判讀的國際背景有缺口。

## What Changes

- 在 SKILL.md 的分批蒐集新增批次 5：國外財經新聞，來源為 WSJ Markets RSS（`https://feeds.a.dj.com/rss/RSSMarketsMain.xml`）與 CNBC 頭條 RSS（`https://www.cnbc.com/id/100003114/device/rss/rss.html`）。
- 批次 5 的定位比照鉅亨 `wd_stock` 分類：內容只寫入 `overview_md` 與產業事件的背景脈絡，不得進 `calls` 與 `stock_news`。
- 明定跨語言去重規則：`seen.py similar` 無法比對中英兩篇同一事件，因此外電只補台媒未報的事件；台媒已報導的事件以台媒版本為準，外電版本判定不報並記為 skipped。
- 批次 5 納入既有的四批彙整表與來源成敗把關（改為五批），失敗處理比照次要來源：不停發、只在執行回報列出。
- 每批回報計數的要求同樣適用：批次 5 要寫下兩個來源各取得幾則、其中未讀幾則。

## Non-Goals

- 不把外電新聞納入 `calls` 或 `stock_news`，即使內容與台灣供應鏈直接相關；個股層級的判讀仍以台媒與證交所來源為準。
- 不修改 seen.py 以支援跨語言相似度比對；去重靠流程規則（外電只補台媒未報者）處理，不靠工具。
- 不新增日經等其他外電來源；日經中文網已實測回 403，本次只收 WSJ Markets 與 CNBC 頭條兩個已驗證可用的 RSS。
- 不改動報告 JSON 結構與 report-ingest-api；外電內容落在既有的 `overview_md` 與產業事件欄位。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `daily-brief-skill`: 蒐集流程從四批擴為五批，新增國外財經新聞批次的來源、用途限制（僅 overview 與產業背景）、跨語言去重規則與計數回報要求。

## Impact

- 封存順序約束：本變更與 primary-failure-spec-fix 都修改 Batched collection with per-batch counts requirement。必須先封存 primary-failure-spec-fix，再套用與封存本變更，否則本變更的五批版本會被四批版本覆蓋（或反向覆蓋掉修正）。
- Affected specs: `daily-brief-skill`
- Affected code:
  - Modified: skills/daily-brief/SKILL.md
  - New: （無）
  - Removed: （無）
