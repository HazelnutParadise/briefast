## Why

每日晨報要改成在雲端排程執行。雲端排程通常用環境變數或 secrets 提供網址與 API key，不會有預先放好的 `.env` 檔，但流程說明規定設定只能從工作區的 `.env` 讀，agent 會在蒐集前停下來。

## What Changes

- 流程說明改成先讀環境變數 `BRIEFAST_URL`、`BRIEFAST_API_KEY`，某一項沒設定或只有空白時，再從工作區根目錄的 `.env` 讀那一項。兩個來源都拿不到才停止並回報缺哪一項。
- README 的每日流程說明與 Cowork 一行指示改成同樣的讀取順序，並說明雲端排程只要設定這兩個環境變數。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `daily-brief-workflow`：Endpoint configuration via environment 改為環境變數優先、`.env` 備用。

## Impact

- 文件：`internal/guide/daily-brief.md` 執行前檢查與送出段落、`README.md` 每日流程。
- 程式與 API 不變。
