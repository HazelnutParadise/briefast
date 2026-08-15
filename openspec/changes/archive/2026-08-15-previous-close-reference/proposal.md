## Why

判讀個股新聞時缺少價格脈絡：agent 無法判斷利多是否已反映在股價上（例如已連漲或漲停後的追高風險），overview 的大盤與個股敘述也只能轉述新聞裡的零散數字。既有的基本面引用紀律禁止憑記憶補值，因此需要一個本次執行實際抓取、日期可驗證的昨收價來源。

## What Changes

- 在 SKILL.md 第 1 步新增「昨收價參考資料」小節：以一般 HTTP 抓取證交所上市個股日成交資訊（openapi.twse.com.tw 的 exchangeReport/STOCK_DAY_ALL）與櫃買中心上櫃日收盤行情（www.tpex.org.tw 的 openapi/v1/tpex_mainboard_daily_close_quotes），各一個請求取回全市場昨收。
- 日期驗證：兩個回應都內嵌民國格式交易日期，agent 必須確認該日期為台北時區的前一個交易日；不符（過舊）時視為抓取失敗，不得沿用舊日期資料而不標示。
- 證交所這個請求納入既有的 TWSE 節流規則（序列、間隔至少 5 秒）；櫃買為不同主機，不受該節流限制。
- 用途規則：昨收價是判讀脈絡（新聞是否已反映、量級對照）與報告引用來源；引用時必須標明價格日期。多空 call 仍以新聞為依據，不得只憑價格漲跌在沒有新聞的情況下建立 call 或 stock_news 條目。
- 失敗分流：重試一次仍失敗就照常判讀與發布（改用純新聞脈絡），缺漏列在執行回報；報告內容不揭露價格資料缺漏，比照既有的蒐集機制不外露規則。
- 計數回報：昨收抓取要回報兩個來源各取得幾筆與資料日期，納入蒐集完成度的回報。

## Non-Goals

- 不用 yfinance 或其他第三方套件：官方 OpenAPI 免金鑰、內嵌日期戳、一次取回全市場，已否決引入非官方依賴的做法。
- 不抓大盤指數與國際指數收盤值：美股與大盤數字維持從新聞報導取得，國際內容依既有規則只進 overview_md。
- 不抓即時行情或盤中價格：即時行情屬證交所收費授權範圍，且 daily-brief 只在開盤前執行，用不到。
- 不抓興櫃與債券行情，也不做歷史價格序列：只取前一交易日單日收盤。
- 不修改報告 JSON 結構：價格數字寫進既有的 overview_md、summary_md 等 Markdown 欄位，不新增價格欄位。
- 不修改既有的 Batched collection with per-batch counts requirement：昨收抓取是獨立的資料步驟，不併入新聞批次，避免與停放中的 foreign-news-batch、primary-failure-spec-fix 對同一 requirement 產生封存順序衝突。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `daily-brief-skill`: 新增昨收價參考資料 requirement（來源端點、日期驗證、TWSE 節流適用、用途與引用規則、失敗不停發），以 ADDED 方式加入，不改動既有 requirement。

## Impact

- Affected specs: `daily-brief-skill`
- Affected code:
  - Modified: skills/daily-brief/SKILL.md
  - New: （無）
  - Removed: （無）
