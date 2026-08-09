---
name: daily-brief
description: Prepare and publish Briefast's daily pre-market industry and Taiwan stock report. Use before market open when an agent must collect public news, judge short- and long-horizon stock calls, compose the exact Briefast report JSON, and POST it to the authenticated report API.
---

# Briefast 每日晨報

每個交易日開盤前依序完成以下四步。只送內容，禁止產生或呼叫 Artifact DSL。

## 執行前檢查

1. 讀取 `BRIEFAST_URL` 與 `BRIEFAST_API_KEY`。
2. 若任一變數未設定或只含空白，立即停止。明確回報缺少哪個變數，不蒐集、不組稿、不 POST。
3. 不顯示、記錄或回傳完整 API key。
4. 以台北時區決定報告日期與 `generated_at`。

## 1. 蒐集產業與股市新聞

### 時間窗口

蒐集從**上一份報告的 `generated_at` 到現在**這段時間內的新聞。

用唯讀端點取得上一份報告的 `generated_at`，不要憑印象或猜測：

```bash
curl --silent --show-error \
  --output previous.json \
  --write-out '%{http_code}' \
  --header "Authorization: Bearer ${BRIEFAST_API_KEY}" \
  "${BRIEFAST_URL%/}/api/report/2026-08-06"
```

日期從台北時區的前一個交易日起算，回 `404` 就再往前一天，最多往前查 5 天。取到 `200` 就以 `previous.json` 的 `generated_at` 為窗口起點；5 天內都是 `404`（首次執行或長假後）就以過去 24 小時為窗口。`401` 比照第 4 步處理，停止並回報。

### 新聞來源

依以下順序蒐集，每個來源都是正式來源，不是補充：

| 來源 | 取得方式 | 重點 |
|---|---|---|
| 鉅亨網 (cnyes) | JSON API（`api.cnyes.com`） | 最優先。文章已標注相關股票代號，有完整內文與 Unix 時間戳 |
| 證交所 (TWSE) | OpenAPI（`openapi.twse.com.tw`） | 重大訊息公告、月營收資料 |
| 工商時報 (CTEE) | 無頭瀏覽器 | 產業與個股深度報導 |
| 中央社 (CNA) | RSS | 財經與產業新聞 |
| 聯合新聞網 (UDN) | RSS | 財經與產業新聞 |
| 科技新報 (TechNews) | RSS | 科技產業新聞 |

全球市場動態（美股收盤、Fed 動態、國際事件）從上述來源的相關報導取得，不另設來源。

### 來源成敗與完整性把關

逐一記錄每個來源這次是成功還是失敗。任何來源失敗就重試一次，再依結果分流：

- **鉅亨網重試後仍失敗**：立即停止。不進判讀、不組稿、不 POST，回報是哪個來源失敗與失敗原因。鉅亨網的文章已標注股票代號並附完整內文，是個股判讀的骨幹，缺了它整份報告會失去主要依據。
- **只有次要來源失敗**：照常判讀與發布，但在執行回報中逐一列出缺哪些來源，讓人知道當日報告的涵蓋程度。

來源成功但當時間窗口內確實沒有新文章，屬於正常情況，不算失敗。

### 去重（seen.py）

去重腳本位於 `skills/daily-brief/scripts/seen.py`，是已讀紀錄的唯一讀寫入口。

**蒐集流程：**

1. 從各來源拉清單頁資料（標題、網址、原生 ID），不下載內文。
2. 把候選餵給 `seen.py peek`，只取回未讀的。
3. 下載未讀文章的內文。
4. 對內文長度足夠的文章，用 `echo "內文" | seen.py similar` 檢查是否為已處理文章的改寫轉載。相似度 ≥ 0.5 視為同篇，跳過。
5. 判讀完成後，把所有處理過的文章（含判定不報的）用 `seen.py record` 記下來。判定不報的設 `"decision": "skipped"`，避免明天重新評估同一批雜訊。
6. 有跨天事件時，用 `seen.py event add '代號|事件類型|期間'` 登記。寫個股詳情前先 `seen.py event check` 查過，已登記的事件用延續敘述（「延續前日營收利多」），不要重新宣布。

**注意：** `event check` 查到已登記，**不代表該檔不能出現在 calls**。報告呈現的是當天的判斷狀態，同一檔連續多天看漲是正常的。事件登記只影響個股詳情的寫法。

## 2. 判讀產業與個股

### 產業分析師視角

`industries[].summary_md` 不複述標題。先指出新聞反映的供需或價格變化，再解釋影響落在產業鏈哪一段、哪些公司類型受益或承壓；同一分類有多個動態時，以 `**子題**` 作為每組內容的開頭。

例如，新聞是「記憶體合約價續漲」，不要只改寫成「記憶體價格上漲」。應寫成：「**記憶體供需** 原廠減產使通路庫存回到健康水位，合約價續漲有利上游製造商改善產品組合，但下游模組廠若無法同步轉嫁成本，毛利率可能承壓。」

### 個股分析師視角

`stock_news[].summary_md` 要把新聞放回公司的營運處境解讀：事件影響哪項產品、訂單、產能或成本，何時可能反映，並說明支持該 call 的因果關係。若當次取得最新基本面資料，再用它檢查新聞影響是否已反映；不要只列事件與多空結論。

例如，公司宣布新產能第四季投產時，應寫出「新產能對現有訂單交付與產品組合的影響，實際貢獻須等第四季投產進度確認」；只有本次確實抓到公開資料時，才可再寫「7 月營收年增……，顯示目前需求……」。

### 基本面引用紀律

- 營收、EPS、毛利率、訂單金額或其他基本面數字，只能引用**本次執行期間實際抓取的公開來源**，不得憑模型記憶、過往對話或常識補值。
- 每個數字都要在正文標明資料期間，例如「7 月營收」、「2026 年第 2 季 EPS」，並在該 `stock_news` 的 `sources` 放入實際來源。
- 若本次抓不到當期公開資料，就省略基本面數字與相關脈絡，只根據已查證新聞解讀。不得沿用舊數字、估算或硬湊近似值。

### 固定產業四分類

`industries` 是封閉分類，只能使用以下四個名稱。每則與台股相關的產業新聞必須依主要影響歸入其中一類，不得新增第五個名稱，也不得讓同一動態重複跨類。

| 分類 | 涵蓋邊界 |
|---|---|
| 科技 | 半導體、電子零組件、電腦與周邊、網通、軟體、電信及其他電子科技供應鏈 |
| 金融 | 銀行、保險、證券、金融控股、支付與金融監理 |
| 傳產 | 航運、鋼鐵、化工、機械、汽車、食品、零售、觀光、生技醫療及其他非科技製造與服務業 |
| 房地產 | 房市交易、建商與營造、商用不動產、住宅政策；若新聞主體是水泥、鋼材等供應商本身，歸傳產 |

跨領域新聞依「主要受影響的需求或監理對象」只歸一類。例如生技廠取得美國藥證歸「傳產」；科技廠購地設廠若重點是產能擴充歸「科技」，若重點是土地市場或開發政策才歸「房地產」。

每個分類的 `summary_md` 以一個或多個 `**當日子題**` 組織內容。逐一查過當天所有來源後，確認某分類在時間窗口內確無新聞，才可省略該分類；不得放空白條目或「今日無新聞」佔位，也不得因漏查直接省略。

### 條目觀察重點

每個 `industries` 與 `stock_news` 條目都必須填非空白的 `watch_md`，用 Markdown 列出 1–2 條可在事後驗證的前瞻觀察點：

- 產業觀察點要連結該分類 `summary_md` 的供需、價格或政策動態，例如下一次合約價公告或監理措施生效結果。
- 個股觀察點要連結該條目引用的新聞或本次抓取的基本面，例如新產能實際投產時間或下一期營收是否反映訂單。
- 不得把 `watch_md` 內容再寫進 `summary_md`，也不得使用「持續關注後續發展」這類沒有具體事件、指標或時間點的範本化空話。

### 個股多空

根據**當天時間窗口內的新聞**判斷方向與時間尺度。當天沒有新聞的股票不列入，不延續前一天的判斷。

| call | 判斷 |
|---|---|
| `short_bull` | 短期看漲 |
| `short_bear` | 短期看跌 |
| `long_bull` | 長期看好 |
| `long_bear` | 長期看壞 |
| `none` | 有重大新聞，但沒有足夠依據判斷方向 |

### 判讀原則

- 不設數量上限。當天新聞能判斷出來的都列。
- 不設強制門檻。由執行模型依新聞內容自行判斷是否足以給出方向。
- 區分已發生事實、公司展望、分析師觀點與市場傳聞。不要把傳聞寫成已確認事實。
- `calls` 與 `stock_news` 只收錄台股上市、上櫃個股及其台股 symbol。國際市場動態與外國上市公司新聞只寫入 `overview_md`，即使與台灣供應鏈相關，也不得把外國 symbol 放入 `calls` 或 `stock_news`。

### calls 條目規則

每個 calls 清單條目必須：

- 使用一行、具體且可由新聞支持的 `reason`。
- 在 `stock_news` 有相同 `symbol` 的詳情。
- 對應的 `stock_news` 至少列一個新聞來源。

若個股有重大新聞但方向不明，只在 `stock_news` 建立 `call: "none"` 的條目，不得放入四個 calls 清單。不要為了填滿欄位硬做判斷；沒有判斷的清單可保持空陣列。

## 3. 組成並驗證報告 JSON

輸出必須符合以下完整結構。欄位名稱與型別不得改動，不得加入額外欄位。

```json
{
  "date": "2026-08-07",
  "headline": "美股收紅開盤氣氛偏多，看漲：台積電",
  "overview_md": "美股四大指數收紅，費半受 AI 需求帶動；台股開盤氣氛偏多。",
  "watch_md": "- 台積電盤後公布 7 月營收\n- 美國 7 月 CPI 今晚公布\n- 觀察外資期貨淨空單",
  "calls": {
    "short_bull": [
      {
        "symbol": "2330",
        "name": "台積電",
        "reason": "2 奈米提前放量且外資調升目標價，短線題材與資金面同步轉強。"
      }
    ],
    "short_bear": [],
    "long_bull": [],
    "long_bear": []
  },
  "industries": [
    {
      "name": "科技",
      "summary_md": "**先進製程**\n\n台積電 2 奈米提前放量可望拉動設備與材料需求。\n\n**AI 供應鏈**\n\n加速器追單使封測產能利用率回升，上游受益較直接。",
      "watch_md": "- 追蹤台積電下一次法說是否維持 2 奈米量產時程\n- 觀察封測業者下月營收是否反映追加訂單"
    }
  ],
  "stock_news": [
    {
      "symbol": "2330",
      "name": "台積電",
      "call": "short_bull",
      "summary_md": "供應鏈傳出 2 奈米製程提前放量，多家外資調升目標價。",
      "watch_md": "- 追蹤公司下一次法說公布的 2 奈米量產時程",
      "sources": [
        {
          "title": "台積電 2 奈米傳提前放量",
          "url": "https://example.com/tsmc-2nm"
        }
      ]
    },
    {
      "symbol": "2317",
      "name": "鴻海",
      "call": "none",
      "summary_md": "宣布電動車合作備忘錄，短期貢獻與長期正式訂單仍不明朗。",
      "watch_md": "- 追蹤合作備忘錄是否轉為具金額與交付期的正式訂單",
      "sources": [
        {
          "title": "鴻海宣布電動車合作備忘錄",
          "url": "https://example.com/hon-hai-ev"
        }
      ]
    }
  ],
  "generated_at": "2026-08-07T07:50:00+08:00"
}
```

送出前逐項檢查：

- `date` 是有效的 `YYYY-MM-DD` 日期。
- `headline`、`overview_md`、`watch_md` 都有非空白內容。
- `calls` 恰好包含 `short_bull`、`short_bear`、`long_bull`、`long_bear` 四個陣列。
- `stock_news[].call` 只使用表格中的五個值。
- `industries[].name` 只使用「科技」「金融」「傳產」「房地產」，且每個 `summary_md` 以粗體子題組織當日內容。
- 每個 `industries` 與 `stock_news` 條目的 `watch_md` 都有 1–2 條非空白、可驗證且未在 `summary_md` 重複的觀察點。
- `calls` 與 `stock_news` 沒有外國上市公司或外國 symbol，只含台股上市櫃個股。
- 每個 calls 條目的 `symbol` 都能在 `stock_news` 找到相同 symbol。
- 每個 calls 條目都有一行理由，且對應 `stock_news` 至少有一個 source。
- 每個 source 都有非空白 URL。
- 方向不明的重大新聞使用 `call: "none"`，且不出現在 calls 清單。
- `generated_at` 使用含時區的 RFC 3339 格式。
- 所有來源的成敗都已記錄，且鉅亨網本次蒐集成功；若有次要來源失敗，缺漏清單已備妥待回報。
- 記得在第 4 步 POST 之前，先用 `GET /api/report/{date}` 確認本次 `date` 是否已有報告；已有就完成產業數與個股數比對。

把最終 JSON 寫到本次工作的暫存檔，例如 `report.json`。不得用假 URL 或範例內容發布正式報告。

## 4. POST 到網站 API

### 同日重跑的覆寫比對

同一日期再次發布會全量覆寫既有報告。送出前先用唯讀端點確認該日期是否已有報告：

```bash
curl --silent --show-error \
  --output existing.json \
  --write-out '%{http_code}' \
  --header "Authorization: Bearer ${BRIEFAST_API_KEY}" \
  "${BRIEFAST_URL%/}/api/report/2026-08-07"
```

`404` 表示該日尚無報告，直接進入送出流程。`200` 表示已有報告，比對 `existing.json` 與本次報告的 `industries` 條目數與 `stock_news` 條目數：

- 兩項都不少於既有版本：照下方流程 POST，不需額外確認。
- 任一項少於既有版本：停止，不要 POST。說明本次為什麼內容較少（例如某來源失敗、當日新聞確實稀少），等待指示後再決定是否覆寫。

內容變少不一定是錯的，但必須由人判斷，不得讓退化的重跑靜默取代較完整的版本。

`401` 表示 key 無效或已撤銷，比照第 4 步的 `401` 處理：停止，不要重試或更換 key。

### 送出

以環境變數組成 endpoint，不要寫死網址或 key：

```bash
curl --silent --show-error \
  --output response.json \
  --write-out '%{http_code}' \
  --request POST \
  --header "Authorization: Bearer ${BRIEFAST_API_KEY}" \
  --header "Content-Type: application/json" \
  --data-binary @report.json \
  "${BRIEFAST_URL%/}/api/report"
```

依 HTTP status 處理：

- `200`：確認 response 為 `{"ok":true,"date":"本次日期"}`，回報發布成功。
- `400`：視為失敗。讀取並保留完整 response body；依 `errors` 修正 payload 後可再送一次。若無法確定修法或修正後仍失敗，停止並逐字回報 response body。
- `401`：視為失敗。不要重試，不要猜測或更換 key；停止並附上完整 response body。
- `5xx`：視為失敗。使用完全相同的 payload 重試一次；第二次仍非 200 就停止並附上最後一次 response body，同時註明第一次與第二次 status。
- 其他非 `200`：視為失敗，不重試，停止並附上完整 response body。

任何非 200 都不得回報成發布成功。重試前不得改變 `date` 或省略原有內容。
