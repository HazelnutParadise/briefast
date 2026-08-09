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

蒐集從**上一份報告的 `generated_at` 到現在**這段時間內的新聞。第一次執行且沒有上一份報告時，蒐集過去 24 小時的新聞。

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

## 2. 判讀個股多空

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
  "headline": "美股收紅開盤氣氛偏多，看漲：台積電、緯創",
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
      "name": "半導體",
      "summary_md": "- 台積電傳 2 奈米提前放量\n- AI 加速器追單帶動封測需求"
    }
  ],
  "stock_news": [
    {
      "symbol": "2330",
      "name": "台積電",
      "call": "short_bull",
      "summary_md": "供應鏈傳出 2 奈米製程提前放量，多家外資調升目標價。",
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
- 每個 calls 條目的 `symbol` 都能在 `stock_news` 找到相同 symbol。
- 每個 calls 條目都有一行理由，且對應 `stock_news` 至少有一個 source。
- 每個 source 都有非空白 URL。
- 方向不明的重大新聞使用 `call: "none"`，且不出現在 calls 清單。
- `generated_at` 使用含時區的 RFC 3339 格式。

把最終 JSON 寫到本次工作的暫存檔，例如 `report.json`。不得用假 URL 或範例內容發布正式報告。

## 4. POST 到網站 API

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
