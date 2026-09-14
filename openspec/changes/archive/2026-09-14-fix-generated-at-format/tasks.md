## 1. 上傳驗證（Report schema validation）

- [x] 1.1 在 internal/report/schema_test.go 的 TestReportValidateIndustryEventsAndStockHeadline 表格加入 generated_at 為 2026-09-14T07:50:27+0800 與空字串的案例，預期錯誤含 generated_at，先確認測試失敗
- [x] 1.2 在 internal/report/schema.go 的 Report.Validate 以 time.RFC3339 解析 generated_at，失敗時加入「generated_at 必須是含時區的 RFC 3339 時間」，讓 1.1 通過

## 2. 顯示容錯

- [x] 2.1 在 internal/site/site_test.go 加測試：首頁渲染 generated_at 為 2026-08-07T07:50:27+0800 的報告時含「07:50 更新」且不含「+0800」，先確認失敗
- [x] 2.2 在 internal/site/site.go 新增 parseGeneratedAt，依序嘗試 time.RFC3339 與 2006-01-02T15:04:05-0700，displayGeneratedAt 與 internal/site/conference.go 的 displayGeneratedFull 改用它，讓 2.1 通過
- [x] 2.3 執行 go test ./... 與 go vet ./... 全數通過，並起本機站台以 sqlite 塞入 +0800 報告截圖確認首頁顯示 07:50 更新，驗完刪除測試資料
