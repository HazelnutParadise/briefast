## 1. 資料契約

- [x] 1.1 實作 Daily market outlook payload 與 Add optional market_outlook to the report JSON：API 接受並原樣讀回有效預測，拒絕無效方向及空白理由，舊報告仍可讀；以 `internal/report` 與 `internal/api` 測試驗證。

## 2. 報告呈現

- [x] 2.1 實作 Market outlook display 與 Render forecast within the existing overview：首頁與歷史頁顯示方向及理由，舊報告不留空白，紅漲綠跌且五區塊順序不變；以 `internal/site` 測試及實際桌面、窄版畫面驗證。

## 3. 每日產出

- [x] 3.1 實作 Evidence-based daily market outlook 與 Synthesize the forecast from report evidence：每日 skill 範例與檢查點涵蓋整份報告、對立訊號及資料不足情況；以內容核對及 skill 範例驗證測試確認。

## 4. 整體驗證

- [x] 4.1 確認所有預測情境及舊資料相容性符合 Implementation Contract，且未引入移轉；執行 `go test ./...`、`spectra validate add-daily-market-outlook`、`spectra analyze add-daily-market-outlook --json` 與 `git diff --check`。
