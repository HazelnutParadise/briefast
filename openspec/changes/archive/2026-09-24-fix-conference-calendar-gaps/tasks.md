## 1. 官方月曆補漏

- [x] 1.1 完成 Official conference calendar reconciliation 與「官方月曆補漏」：每日流程取得當月、次月的上市與上櫃四份官方月曆，檢查 myTable、回報各份筆數與失敗，依市場、代號、日期去重並與 upcoming 比對；用已存的官方月曆檔人工核對 7753/2026-10-14 會列為缺漏，以及任一表格缺失不被算成零場。

## 2. 分類與流程驗收

- [x] 2.1 完成 Conference announcement detection 與「受邀場次判斷」：調整每日流程的文字分類，使券商受邀場次仍排除、「參加方式」不誤排、模糊場次會查原公告或回報待查；以 2408、7753 及受邀參加永豐證券的官方樣例逐項核對分類結果。
- [x] 2.2 核對「Implementation Contract」及「Risks / Trade-offs」：確認漏登自辦場次進入既有 brief 與 POST 流程、失敗逐份回報、未加入空白報告或變更網站/API；執行 spectra validate fix-conference-calendar-gaps 與 git diff --check 驗證文件和改動。
