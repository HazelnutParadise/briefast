## 1. 夜盤資料解析

- [x] 1.1 實作 Verified TX night extraction 與「用期交所 CSV 與純標準函式庫擷取夜盤方向」：命令列只從指定日期、時段、最高成交量的單月 TX 合約及前一交易日同契約日盤計算完整 JSON，錯誤不輸出部分 JSON；以 `python3 skills/daily-brief/scripts/tx_night_test.py` 驗證正常、跨週末、換月、零變動與錯誤案例。

## 2. 每日晨報判讀

- [x] 2.1 落實 Evidence-based daily market outlook 與「用同契約前一日盤收比較且只影響開盤情境」：`skills/daily-brief/SKILL.md` 指引保存兩日官方 CSV、核對已完成夜盤、將有效方向只作開盤佐證並與美股及台股資訊比對；以內容檢查確認盤前總覽、`market_outlook` 欄位及盤中畫圖門檻未改。
- [x] 2.2 落實「失敗時回退既有市場判讀」：下載或解析失敗重試一次後略過 TX，缺漏只進執行回報，晨報仍以其他證據完成；以內容檢查與 parser 錯誤案例確認不沿用舊日或異契約資料。

## 3. 完成驗證

- [x] 3.1 以官方 2026-09-24 與 2026-09-23 TX 原始 CSV 執行命令列，核對 202610 合約差額 -427 與方向 down；執行 `spectra validate add-tx-night-market-signal`、`python3 skills/daily-brief/scripts/tx_night_test.py` 及 `git diff --check`，確認規格、解析器和文件一致。
