## 1. 文件

- [x] 1.1 改寫 `internal/guide/daily-brief.md` 執行前檢查第 1、2 點與送出段落，依 Endpoint configuration via environment 改成環境變數優先、`.env` 備用；以內容檢查確認不再出現「不從環境變數取值」。
- [x] 1.2 更新 `README.md` 每日流程與 Cowork 一行指示，說明雲端排程只要設定兩個環境變數。

## 2. 驗證

- [x] 2.1 執行 `go test ./...`、`openspec validate read-config-from-env-vars --strict`、`git diff --check`。
