## Context

目前 API 已允許省略路徑欄位，畫面也會省略空缺區塊。問題來自每日 skill 與文件強制要求路徑。

## Goals / Non-Goals

納入：開盤與收盤獨立判斷、產出規則、JSON 範例、文件及相容性驗證。
排除：API、資料庫、歷史報告、部署、移除既有歷史圖表渲染功能。

## Decisions

沿用 direction 與 summary_md，不新增開盤欄位。新報告省略兩個 trajectory 欄位，舊報告照原有方式呈現。

## Implementation Contract

輸入為盤前已查證的新聞、美股、夜盤與台股資料。輸出 market_outlook 僅含 direction 與非空白 summary_md。開盤傾向在理由內描述，收盤方向由 direction 表達；任一階段資料不足都不影響另一階段有依據的判斷。兩者皆不足則使用 uncertain 並交代關鍵觀察。禁止以夜盤推定收盤，禁止新增盤中預測或路徑無法判斷的固定句。

文件格子：README.md 的 JSON 範例與說明、DESIGN.md 的大盤呈現約定須與此輸出一致，說明歷史欄位仍相容。不得修改 API 或樣式。

驗證格子：skill 與 README JSON 範例須通過正式 Report.Validate，兩個 trajectory 欄位皆不存在，理由明確提及開盤與收盤。首頁與歷史頁用現有 renderer 呈現新格式，沒有路徑標示或圖，且舊格式測試持續通過。

## Risks / Trade-offs

理由文字需清楚區分開盤與收盤，避免將方向標籤誤認為開盤預測。執行中的 agent 若已載入舊版 skill，必須重新載入才會採用新規則。
