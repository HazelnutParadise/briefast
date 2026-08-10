## 1. 設定讀取

- [x] 1.1 依 Endpoint configuration via environment 修改 skills/daily-brief/SKILL.md 執行前檢查：改為從執行工作區根目錄的 .env 檔（KEY=VALUE 格式）讀 BRIEFAST_URL 與 BRIEFAST_API_KEY，檔案不存在或任一鍵缺漏、空白即停止並回報，不蒐集不組稿不 POST；不顯示、不記錄 key 的規則保留。驗證方式為內容審查，確認全文不再出現「環境變數」作為設定介面的敘述。
- [x] 1.2 新增 skills/daily-brief/.env.example：僅含 BRIEFAST_URL 與 BRIEFAST_API_KEY 兩鍵、佔位值與一行用途說明，無真實秘密；並在 repo 根目錄 .gitignore 加入 .env 條目。驗證方式：cat 檔案確認內容、git check-ignore .env 回報被忽略。

## 2. 來源表改版

- [x] 2.1 依 Pinned source endpoints 改寫 SKILL.md 新聞來源表：鉅亨網列出六個分類 API（tw_stock、headline、tech、tw_macro、cnyeshouse、wd_stock）與「列表即含全文與股票標注、不抓詳情頁、wd_stock 只餵 overview_md」；證交所列出 t187ap04_L 與 t187ap05_L 兩資料集與「整包拉回依窗口過濾」；工商時報改為 rss_web/livenews 六分類 RSS（policy、stock、finance、industry、house、tech）加一般 HTTP 抓內文，移除無頭瀏覽器敘述；中央社財經 RSS 與科技新報全站 feed 附上網址；移除聯合新聞網、新增自由時報財經 RSS。驗證方式為內容審查，確認每個來源都有具體網址且無殘留 UDN 與無頭瀏覽器敘述。
- [x] 2.2 依 TWSE request throttling 在 SKILL.md 加入證交所節流規則：序列請求、間隔至少 5 秒、不並行；429 或連線被拒等 60 秒重試一次，仍失敗依來源成敗把關當次要來源失敗（照發並揭露缺 TWSE）。驗證方式為內容審查，確認與既有「來源成敗與完整性把關」段落銜接一致。

## 3. 兩段式判讀與撰寫訣竅

- [x] 3.1 依 Two-pass ranked judgement 重組 SKILL.md 判讀章節為先產業輪、後個股輪兩段，各自載明五項排序依據（產業：價格供需已變動、政策將生效、波及廣度、持續性、確定性；個股：獲利可量化、事實層級、反映時點、競爭地位、籌碼僅佐證），並明文規定排序只影響事件領銜、詳略與先後，不得用來篩選收錄，不設數量上限規則原文保留。驗證方式為內容審查，確認兩組依據完整、附「排序影響詳略但不影響收錄」的明確句子。
- [x] 3.2 依 Analyst summary writing craft 在兩輪各寫入撰寫訣竅（產業四條：數字含方向幅度期間與基準且以新聞載明為限、供需歸因附持續性且驅動不明就明寫不明、影響分佈寫誰受益誰承壓誰能轉嫁、深寫事件標明週期或結構；個股五條：首句寫改變了什麼判斷且比 headline 深一層而 call none 首句寫為何無法判斷、量級換算算得出才寫算不出直接省略不留佔位句、預期差限新聞載明或已登記事件、公關語言還原為事實且無數字合作降級為意向、具名歸因以新聞載明為限），附反例→正例對照表並註明為示意非範本句。驗證方式為內容審查，逐條確認「來源綁定」與「省略式出口」兩道護欄都有寫進去，且「明寫缺料」只出現在供需歸因一條。
- [x] 3.3 整體一致性驗證：以未參與討論的接手者視角重讀 SKILL.md，確認 .env 讀取、來源表、節流、兩段式判讀、撰寫訣竅與既有規則（來源把關、一句話標題、觀察重點、基本面紀律、送出前檢查清單）無矛盾且無殘留舊敘述；執行 spectra validate daily-brief-source-overhaul 確認通過。
