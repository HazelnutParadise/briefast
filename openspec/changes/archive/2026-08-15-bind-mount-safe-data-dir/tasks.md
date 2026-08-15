## 1. Entrypoint 降權

- [x] 1.1 依 Data directory works on any host mount 新增 docker-entrypoint.sh：以 root 執行時建立 BRIEFAST_DB_PATH 所在目錄、把擁有者設為 briefast，再以 su-exec 降權執行傳入的指令；以非 root 執行時略過修正直接執行。驗證方式：檢視腳本確認兩條路徑都存在且不硬編資料目錄路徑。
- [x] 1.2 修改 Dockerfile：安裝 su-exec、以可執行權限複製 entrypoint、移除 USER briefast 並改用 ENTRYPOINT 搭配既有 CMD。驗證方式：docker build 成功，且 docker run 後以 ps 或 id 確認應用程序的執行使用者是 briefast 而非 root。

## 2. 驗證與文件

- [x] 2.1 實測兩種主機掛載情境：掛到不存在的主機路徑、掛到既有且屬於其他使用者的主機路徑，兩者都要能啟動並回應首頁 200，且過程中未在主機執行任何 chown。驗證方式：docker run 指定 bind mount 後從主機 curl 首頁取得 200。
- [x] 2.2 更新 README 的 Docker 部署一節，說明 BRIEFAST_DATA 可直接指向主機目錄、不需事先建立或調整權限。驗證方式為內容審查，確認敘述與 entrypoint 的實際行為一致。
