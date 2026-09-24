## 1. 資料契約

- [x] 1.1 完成 Daily market outlook payload 與 Add optional qualitative trajectory chart data：API 有效圖原樣讀回；缺階段、錯值、收盤方向矛盾、有圖卻無走法文字時回 400 且不寫入；以 schema 與 API 測試確認，舊報告仍可讀。

## 2. 畫面呈現

- [x] 2.1 完成 Market outlook display 與 Render one static three-stage SVG：首頁及歷史頁顯示有階段、層級與非點位聲明的示意圖及文字替代描述，無圖不留佔位；以站台測試及桌面、手機實際畫面確認。

## 3. 每日內容

- [x] 3.1 完成 Evidence-based daily market outlook 與 Preserve uncertainty and legacy reports：每日範例有與方向及文字一致的三階段圖；不確定時省略圖且保留待觀察條件；以範例驗證測試及內容核對確認。

## 4. 整體驗證

- [x] 4.1 對照 Implementation Contract 確認新舊相容、圖表語意與無遷移，執行 `go test ./...`、`spectra validate add-market-trajectory-diagram`、`spectra analyze add-market-trajectory-diagram --json` 及 `git diff --check`。
