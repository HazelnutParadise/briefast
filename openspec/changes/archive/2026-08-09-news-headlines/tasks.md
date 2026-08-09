## 1. 報告型別與驗證

- [x] 1.1 依 Report schema validation 與 design 的「產業改為事件陣列而非 Markdown 慣例」決策修改 internal/report/schema.go：Industry 型別移除 SummaryMD、改為 JSON 名稱 events 的事件切片，新增含 headline 與 summary_md 的事件型別，StockNews 型別新增 JSON 欄位 headline。驗證方式：go build ./... 通過且型別的 JSON 標籤與 design 的資料形狀一致。
- [x] 1.2 依 design 的「驗證嚴格度與錯誤訊息形狀」決策擴充 Report 驗證：industries 條目的 events 為空、事件的 headline 或 summary_md 空白、stock_news 條目的 headline 空白，皆逐項寫入 errors 並標明是第幾個產業與第幾個事件，不在第一個違規就中止。internal/report/schema_test.go 涵蓋空 events、事件缺 headline、事件缺 summary_md、個股缺 headline、正常五種案例。驗證方式：go test ./internal/report/... 全綠。
- [x] 1.3 更新 internal/api/report_test.go 與 internal/api/read_test.go 的測試報告資料改用新結構，並新增整合案例：缺標題的 payload 經 POST /api/report 回 400、errors 標明索引、報告不落庫且 update_log 記 rejected-schema。驗證方式：go test ./internal/api/... 全綠。

## 2. 版面渲染

- [x] 2.1 依 Industry event headlines display 與 design 的渲染行為修改 internal/site：產業分類內逐一渲染事件，事件標題為獨立於內文的小標題、內文接在其下，同一分類的多個事件在視覺上可分辨；樣式定義在程式內並遵循 DESIGN.md（零圓角、零陰影、規線只到大區塊）。驗證方式：go test ./internal/site/... 全綠，且測試斷言兩個事件的標題與內文都出現在產業區塊。
- [x] 2.2 依 Stock news detail display 與 design 的「個股標題採欄位而非寫進 summary_md」決策，把個股標題渲染在左側識別欄，位於股票名稱與多空標籤之間；internal/site/site_test.go 斷言標題出現在識別欄而非內文欄。驗證方式：go test ./internal/site/... 全綠。
- [x] 2.3 在 internal/site/styles.go 補上事件標題與個股標題的樣式，確認標題與內文在字級、字重或字體上可分辨，且不新增圓角與陰影。驗證方式：本機啟動網站以示範報告檢視產業與個股區塊，確認標題與內文有明顯層級差異。

## 3. SKILL.md 與整體驗證

- [x] 3.1 依 One-sentence headlines 更新 skills/daily-brief/SKILL.md：規範事件與個股標題必須寫出發生什麼事與影響、不得使用短語標籤、不得複述內文首句，並規範同一分類的事件依主題合併而非一則新聞一個事件。驗證方式為內容審查，確認四項規則齊備且附具體正反例。
- [x] 3.2 依 Fixed industry taxonomy 更新 SKILL.md 的固定四分類段落：把「以粗體子題組織內容」改為「拆成事件，每個事件有自己的標題與內文」，並明訂保留的分類不得留空事件清單；同時依 Watch points in dedicated fields 把產業觀察重點的敘述改為連結該分類各事件的動態，並確認禁止重複的對象改為事件內文。驗證方式為內容審查，確認兩處敘述都不再提及舊的 summary_md 單欄結構。
- [x] 3.3 更新 SKILL.md 的報告 JSON 範例為新結構（industries 含 events、事件有 headline 與 summary_md、stock_news 有 headline），並在送出前檢查清單加入三項檢核：每個產業至少一個事件、每個事件的 headline 與 summary_md 非空白、每檔個股 headline 非空白。驗證方式為內容審查，確認範例、規則與檢查清單三者一致。
- [x] 3.4 整體驗證：執行 go test ./... 全綠；本機啟動網站送出一份新格式示範報告確認 200 並正常渲染，另送一份缺標題的 payload 確認回 400 且不落庫；執行 spectra validate news-headlines 確認通過。
