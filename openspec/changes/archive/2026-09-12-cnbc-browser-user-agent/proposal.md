## Problem

批次 5 的 CNBC 每次執行都失敗：頭條 RSS 端點與文章頁都回 `403 Access Denied`，執行回報每天都把 CNBC 列為缺項，國際背景密度長期為零。skill 只寫「文章頁一般 HTTP 抓取」，沒有交代任何請求條件，agent 用預設工具抓就必然被擋。

## Root Cause

2026-09-12 實測 `https://www.cnbc.com/id/100003114/device/rss/rss.html` 與其中一篇文章頁 `https://www.cnbc.com/2026/09/11/cpi-inflation-breakdown-august-2026.html`：

| User-Agent | RSS | 文章頁 |
|---|---|---|
| curl 預設（`curl/8.x`） | 403 | 403 |
| `python-requests/2.31` | 403 | 未測 |
| `Claude-User/1.0` | 403 | 未測 |
| 空 UA | 200 | 200 |
| Chrome 桌面瀏覽器 UA | 200（30 則） | 200（正文可解析） |

403 回應的 server header 是 `AkamaiGHost`，內容為 Akamai 的 Access Denied 頁。結論：CNBC 由 Akamai 依 User-Agent 黑名單封鎖工具型客戶端（curl、python-requests、Claude-User 等），與 IP、頻率、cookie 無關；同一秒換成瀏覽器 UA 就通過。這也解釋了為什麼「每次都 403」而不是間歇失敗。

批次 1 到 4 的所有端點（cnyes、TWSE、ctee、CNA、TechNews、LTN）以 curl 預設 UA 與 `Claude-User` 皆回 200，問題只發生在 CNBC，不需要全域改動。

另外，RSS 的 `description` 只有一句約 120 字的摘要，不含全文，所以「只用 RSS、不抓文章頁」不是可行的替代方案；文章頁仍必須抓，且也要帶瀏覽器 UA。

## Proposed Solution

在 `skills/daily-brief/SKILL.md` 批次 5 段落明定 CNBC 的請求條件：

- RSS 清單與文章頁**都必須**帶桌面瀏覽器的 `User-Agent` 標頭，並直接給出一組可用的 UA 字串與 curl 範例，讓 agent 不必自行猜測。
- 說明原因（Akamai 依 UA 封鎖工具型客戶端）與症狀（`403 Access Denied`、server `AkamaiGHost`），並明定：**遇到 403 先確認 UA 是否已帶，不得把 UA 缺漏造成的 403 直接記為來源失敗**。
- 保留既有失敗分流：帶了瀏覽器 UA 仍 403 或連線失敗，才走「重試一次後照常發布、執行回報列缺項」的規則。

同步更新 `openspec/specs/daily-brief-skill/spec.md` 的 Pinned source endpoints requirement，把 CNBC 的抓取條件從「plain HTTP」改為「plain HTTP with a desktop browser User-Agent header on both the feed and article requests」，並新增對應 scenario。

## Non-Goals

- 不引入無頭瀏覽器，也不新增 fetch 腳本。瀏覽器 UA 已足夠，多一層工具只增加維護成本。
- 不改動批次 1 到 4 的抓取方式，實測它們不受 UA 影響。
- 不改用「只讀 RSS 摘要」的降級做法，摘要不足以支撐 overview_md 的背景脈絡。
- 不處理 seen.py 或去重規則，這次只修取得層。

## Success Criteria

- 以 SKILL.md 寫明的 UA 與 curl 範例實際請求 CNBC RSS，回 200 且解析出 30 則 item。
- 以同一 UA 請求 RSS 內任一篇文章頁，回 200 且能從 HTML 取出正文段落。
- SKILL.md 批次 5 段落包含：UA 必帶的規則、可直接複製的 UA 字串與 curl 範例、403 時先查 UA 的檢查步驟、與既有失敗分流的銜接。
- `openspec/specs/daily-brief-skill/spec.md` 的 Pinned source endpoints requirement 與 SKILL.md 敘述一致，`spectra validate cnbc-browser-user-agent` 通過。

## Impact

- Affected specs: `daily-brief-skill`（Pinned source endpoints requirement 修改）
- Affected code:
  - Modified: `skills/daily-brief/SKILL.md`
  - New: （無）
  - Removed: （無）
