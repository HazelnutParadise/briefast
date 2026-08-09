## 1. SKILL.md 判讀角色與基本面紀律

- [x] 1.1 在 skills/daily-brief/SKILL.md 新增判讀角色指引（Analyst perspective guidance）：產業區塊以產業分析師視角解讀供需變化、價格與跨公司影響而非複述標題，個股區塊以個股分析師視角把新聞對照公司處境與基本面脈絡解讀。完成標準：兩種視角各有獨立段落與具體寫法示例；驗證方式為內容審查，確認指引能直接指導執行模型下筆而非抽象口號。
- [x] 1.2 在 skills/daily-brief/SKILL.md 新增基本面引用紀律（Fundamentals freshness discipline）：報告中任何基本面數字必須來自當次執行抓取的公開來源並標明資料期間（如「7 月營收」），禁止憑模型記憶填寫，抓不到當期資料就省略基本面脈絡。驗證方式為內容審查，確認規則同時涵蓋來源要求、期間標註、禁止記憶、省略處理四個要素。

## 2. SKILL.md 產業四分類與個股結構

- [x] 2.1 將 industries 改為固定四分類（Fixed industry taxonomy）：科技、金融、傳產、房地產為封閉分類並附各分類涵蓋邊界，區塊內以粗體子題帶當日內容，分類經查證確無新聞才可省略；同步把 SKILL.md 報告 JSON 範例的 industries 條目改為固定名稱與粗體子題格式。驗證方式為內容審查，確認規則、邊界表與 JSON 範例三者一致。
- [x] 2.2 依 Watch points in dedicated fields 更新 SKILL.md：industries 與 stock_news 每個條目都必須填非空白 watch_md（1–2 條可事後驗證的前瞻觀察點，連結該條目引用的新聞或基本面），watch 內容不得在 summary_md 重複、不得用範本化空話；同步更新 JSON 範例，兩處條目都展示 watch_md。驗證方式為內容審查。
- [x] 2.3 明定台股涵蓋邊界（Taiwan-market coverage boundary）：calls 與 stock_news 僅限台股上市櫃個股，國際市場與外國個股動態只寫入 overview_md。驗證方式為內容審查，確認規則出現在判讀原則段落。

## 3. Go 實作：schema 與 API 驗證

- [x] 3.1 依 Report schema validation 與 design 的 watch_md 必填範圍決策修改 internal/report/schema.go：Industry 與 StockNews 型別新增 JSON 欄位 watch_md，驗證每個條目非空白（缺漏或僅空白視為違規），違規訊息逐條列入 errors 陣列；internal/report/schema_test.go 新增缺欄位、空白字串、正常三種案例。驗證方式：go test ./internal/report/... 全綠。
- [x] 3.2 在 internal/api/report_test.go 新增整合案例：缺 watch_md 的 payload 經 POST /api/report 回 400、errors 列出違規、報告不落庫且 update_log 記 rejected-schema。驗證方式：go test ./internal/api/... 全綠。

## 4. Go 實作：版面渲染

- [x] 4.1 依 Watch points display 與 design 的觀察重點渲染方式、歷史報告相容策略決策修改 internal/site：產業與個股條目在摘要內容後渲染標示「觀察重點」的獨立子區塊，樣式遵循 DESIGN.md（零圓角、零陰影、規線只到大區塊）；watch_md 空缺或空白時子區塊完全不渲染、無佔位文字。internal/site/site_test.go 覆蓋含 watch_md 與不含 watch_md（歷史報告）兩種報告的渲染結果。驗證方式：go test ./internal/site/... 全綠。

## 5. 檢查清單與整體驗證

- [x] 5.1 更新 SKILL.md「送出前逐項檢查」清單：新增三項檢核——industries 名稱僅限四個固定值、industries 與 stock_news 每個條目的 watch_md 非空白、calls 與 stock_news 無外國個股 symbol。驗證方式為內容審查，確認檢查清單與第 1、2 群組新增的規則一一對應。
- [x] 5.2 整體一致性驗證：以未參與討論的接手者視角重讀 SKILL.md，確認新增規則彼此之間、與既有 calls 條目規則及 seen.py 流程無矛盾且無殘留「summary_md 尾段觀察重點」舊敘述；執行 go test ./... 全綠，並執行 spectra validate daily-brief-analyst-guidance 確認通過。

## 6. 發布把關規則

- [x] 6.1 依 Source collection completeness gate 與 design 的「蒐集完整性把關放在 skill 而非 API」決策，在 SKILL.md 蒐集章節加入來源成敗記錄與重試規則：每個來源記錄成功或失敗、失敗重試一次；鉅亨網重試後仍失敗即停止流程、不組稿不 POST 並回報失敗來源與原因；僅次要來源失敗則照常發布並在執行回報列出所有缺漏來源。驗證方式為內容審查，確認記錄、重試、主次分流、缺漏揭露四個要素齊備且與既有來源表一致。
- [x] 6.2 依 Same-day overwrite guard 與 design 的「同日重跑的覆寫防護」決策，在 SKILL.md POST 章節加入送出前比對：該日期已有報告時，先比對新報告與已發布版本的產業數與個股數，任一較少就停止並說明差異原因、等待指示，不得逕行 POST。驗證方式為內容審查，確認規則與既有 HTTP status 處理及「重試前不得改變 date 或省略原有內容」無矛盾。
- [x] 6.3 把兩條把關規則納入 SKILL.md 送出前檢查清單，並以未參與討論的接手者視角重讀全文確認流程順序正確（蒐集把關在組稿前、覆寫比對在 POST 前）；執行 spectra validate daily-brief-analyst-guidance 確認通過。

## 7. 唯讀端點

- [x] 7.1 依 Authenticated report read endpoint 與 design 的「唯讀端點的驗證與記錄方式」決策實作 GET /api/report/{date}：沿用 ingest 的 Bearer 驗證並每次重查 key 狀態，200 回傳該日完整報告 JSON、401 記 read_rejected_auth、日期格式不合法回 400、查無報告回 404，非 200 一律帶 ok=false 與 errors，成功讀取不寫 update_log；於 main.go 註冊路由且不影響既有 POST /api/report。驗證方式：go test ./... 全綠。
- [x] 7.2 在 internal/api/report_test.go 補涵蓋五種情況的測試：有效 key 讀到完整報告（含 watch_md）、撤銷後的 key 回 401 且寫入 read_rejected_auth、無 Authorization 標頭回 401、日期格式錯誤回 400、日期正確但無報告回 404 且不寫 log。驗證方式：go test ./internal/api/... 全綠。
- [x] 7.3 更新 SKILL.md，把同日覆寫比對與時間窗口規則改為實際呼叫 GET /api/report/{date} 取得已發布報告，附與既有 POST 相同風格的 curl 範例並沿用環境變數，不寫死網址或 key。驗證方式為內容審查，確認兩處規則都指向此端點且與第 4 步流程順序一致。
