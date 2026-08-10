## Context

Briefast 的容器化目前由兩個檔案描述：Dockerfile 以多階段建置產出靜態 binary，runtime stage 只 COPY 該 binary；docker-compose.yml 掛載一份主機上的 `syralit.toml` 到容器的 `/app/syralit.toml`，並提供 `briefast-data` named volume 給 SQLite。

實測確認的行為：工作目錄有 `syralit.toml` 時，程式以檔內的 `host = "0.0.0.0"` 綁定；沒有該檔時，Syralit 退回內建預設 `127.0.0.1:8600`。因此目前的映像單獨執行會綁在 loopback，容器外無法連線；服務可用完全依賴 compose 補上的掛載。本次把設定改為建置期包進映像。

## Goals / Non-Goals

**Goals:**

- 映像自帶執行所需的完整設定，單獨執行即可對外服務。
- 部署期的注入面只剩真正屬於環境的東西：機密與埠。
- 容器的持久化範圍只有資料庫。

**Non-Goals:**

- 不改設定檔的內容值與 Go 的讀取邏輯。
- 不改 `briefast-data` volume 與資料庫路徑。
- 不新增環境差異的設定覆寫機制。
- 不改 `.env` 承載機密與 compose 變數的角色。

## Decisions

### 設定檔包進映像

`syralit.toml` 決定綁定位址、主題與 i18n，這些隨程式版本一起變動、不隨部署環境變動，屬於程式碼而非使用者設定。把它 COPY 進 runtime stage，映像因此自成完整可執行單元；同一個映像在任何編排工具下的行為都一致，不再取決於外部有沒有補掛載。捨棄「維持掛載但在 Dockerfile 也放一份預設」的折衷做法，理由是那會讓同一份設定有兩個來源，載入順序與覆蓋關係要另外解釋，換不到對應價值。

### 移除 BRIEFAST_CONFIG

程式從不讀這個變數，它只是 compose 決定掛載來源的參數。設定檔進映像後掛載本身消失，這個變數也失去作用；留著只會讓文件繼續暗示「設定檔是可替換的部署旋鈕」，而實際上本機執行時它完全無效。一併從 `.env.example`、README 與 AGENTS.md 移除，避免留下無效的操作指示。日後若真需要環境差異設定，應設計明確的覆寫機制而不是復活這個半可用的變數。

### 容器只持久化資料庫

移除設定檔掛載後，compose 唯一的 volume 是資料目錄。這符合持久化的判準：只有跨容器生命週期需要保留的可變狀態才需要 volume，SQLite 資料庫是唯一符合的項目，設定檔隨映像重建即可。

資料落地位置本身則交給 `.env`：掛載來源寫成 `${BRIEFAST_DATA:-briefast-data}`，不設定時沿用 named volume，填入主機路徑時 Compose 自動視為 bind mount。這讓「資料放哪」成為部署者的決定而不必改 compose 檔，與埠、機密同屬 `.env` 的職責。容器內的資料庫路徑維持固定，不隨之外露成變數——路徑是程式與映像的內部約定，可變的只有它對應到主機上的哪裡。

## Implementation Contract

- **映像內容**：建置後的映像在 `/app/syralit.toml` 存在該檔，內容與 repo 中的版本相同；runtime stage 的檔案擁有者與既有 `/app` 一致，執行使用者可讀取。
- **執行行為**：不掛載任何設定檔的情況下啟動容器，程式以 `0.0.0.0:8600` 綁定，容器外可連線並取得首頁 200。
- **compose 行為**：`docker compose config` 展開後只有一個對應 `/app/data` 的掛載，沒有任何指向 `/app/syralit.toml` 的 bind mount。未設 `BRIEFAST_DATA` 時該掛載是 named volume `briefast-data`；設為主機路徑時展開為該路徑的 bind mount。埠與 `BRIEFAST_ADMIN_PASSWORD` 仍取自 `.env`。
- **文件一致性**：`.env.example`、README 與 AGENTS.md 不再出現 `BRIEFAST_CONFIG`，並說明 `BRIEFAST_DATA` 的用途；`syralit.toml` 檔頭註解說明它隨程式發佈、不是部署期注入點。
- **驗收標準**：建置映像後以不帶任何掛載的方式啟動容器，從主機取得首頁 200，並確認容器內存在 `/app/syralit.toml`；`docker compose config` 在有無 `BRIEFAST_DATA` 兩種情況下分別展開為主機路徑與 named volume、且皆無設定檔掛載；grep 確認活躍文件無 `BRIEFAST_CONFIG` 殘留；go test 全綠。
- **範圍邊界**：in scope 為 Dockerfile、`.dockerignore`、docker-compose.yml、`.env.example`、README.md、AGENTS.md、`syralit.toml` 註解與 `deployment` spec；out of scope 為設定內容值、Go 讀取邏輯、資料庫 volume 與 `.env` 的機密角色。

## Risks / Trade-offs

- [設定變更必須重建映像才生效，不能改主機檔重啟] → 這正是把設定視為程式碼的預期結果；設定與程式版本綁定反而讓部署可重現。真正屬於環境的埠與機密仍走 `.env`，不需要重建。
- [既有以 compose 部署的環境若曾用 `BRIEFAST_CONFIG` 指向自訂檔，升級後該檔不再生效] → 目前唯一的部署路徑是 repo 內的 compose，預設值一直指向 repo 中的設定檔，沒有已知的自訂檔使用者；README 與 AGENTS.md 會明說這個變數已移除。
- [映像體積增加] → 增加的是一個約 1 KB 的文字檔，可忽略。
