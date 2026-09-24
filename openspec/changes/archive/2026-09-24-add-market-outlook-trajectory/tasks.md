## 1. 資料契約

- [x] 1.1 實作 Daily market outlook payload 與 Add an optional trajectory_md field：API 接受並讀回非空白路徑，拒絕明確提供的空白路徑，舊版預測及無預測報告仍可讀；以報告 schema 與 API 測試驗證。

## 2. 畫面呈現

- [x] 2.1 實作 Market outlook display 與 Promote trajectory in the overview：首頁和歷史頁先顯示「預期走法」及內容，再顯示依據；舊報告不出現空白標籤；以站台測試和實際桌面、窄版畫面驗證。

## 3. 每日內容

- [x] 3.1 實作 Evidence-based daily market outlook 與 Write a baseline path with conditional pivots：每日流程範例及檢查點包含開盤、盤中、尾盤與改變走法的條件，證據不足時明示無法判斷；以 skill 範例驗證測試與內容核對確認。

## 4. 整體驗證

- [x] 4.1 對照 Implementation Contract 確認相容性、渲染及無資料遷移；執行 `go test ./...`、`spectra validate add-market-outlook-trajectory`、`spectra analyze add-market-outlook-trajectory --json` 和 `git diff --check`。
