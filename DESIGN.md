---
version: alpha
name: Briefast
description: 台灣 AI 每日股市晨報——FT 系規線報紙版式，暖紙面、襯線大標、紅漲綠跌雙品牌色
colors:
  primary: "#33302E"
  secondary: "#F9F1E7"
  tertiary: "#990F3D"
  neutral: "#66605A"
typography:
  masthead:
    fontFamily: Noto Serif TC
    fontSize: 58px
    fontWeight: 900
    lineHeight: 1.1
    letterSpacing: 0.015em
  h1:
    fontFamily: Noto Serif TC
    fontSize: 33px
    fontWeight: 700
    lineHeight: 1.45
  h2:
    fontFamily: Noto Serif TC
    fontSize: 24px
    fontWeight: 700
    lineHeight: 1.3
  body-serif:
    fontFamily: Noto Serif TC
    fontSize: 15.5px
    fontWeight: 400
    lineHeight: 1.9
  body-md:
    fontFamily: Noto Sans TC
    fontSize: 15px
    fontWeight: 400
    lineHeight: 1.75
  label-caps:
    fontFamily: Noto Sans TC
    fontSize: 12.5px
    fontWeight: 700
    lineHeight: 1
    letterSpacing: 0.14em
rounded:
  none: 0px
spacing:
  xs: 4px
  sm: 8px
  md: 16px
  lg: 32px
  xl: 64px
components:
  section-head:
    borderTop: "3px solid {colors.primary}"
    titleFont: "{typography.h2}"
  tag-bull:
    textColor: "#990F3D"
    backgroundColor: "#F2DCE2"
  tag-bear:
    textColor: "#1E6A45"
    backgroundColor: "#DFEBE0"
  footer-band:
    backgroundColor: "#262A33"
    textColor: "#C9C3BC"
---

## Overview

Briefast 是一份「每個交易日開盤前出刊的 AI 股市晨報」。視覺定位是嚴肅財經刊物頭版，不是券商 App、不是行銷頁。設計語言遷移自 Financial Times 的版式系統（Origami：暖紙面、襯線 display＋無襯線 UI 雙字型制、規線分欄、無卡片無陰影），未使用 FT 商標或其註冊商用色。讀者是台灣散戶與上班族，場景是開盤前三到五分鐘快速掃讀。

## Colors

淺色（預設主題）：

| Token | 值 | 用途 |
|---|---|---|
| paper | `#F9F1E7` | 頁面紙底（暖米色） |
| ink | `#33302E` | 主文字、粗規線 |
| ink-soft | `#66605A` | 次要文字、meta |
| rule | `#E0D8CA` | 髮絲分隔線（僅大區塊層級） |
| rule-heavy | `#33302E` | 3px 粗規線（報頭、章節首） |
| up | `#990F3D` | 看漲紅（酒紅，兼品牌重點色） |
| up-tint | `#F2DCE2` | 看漲標籤底 |
| down | `#1E6A45` | 看跌綠 |
| down-tint | `#DFEBE0` | 看跌標籤底 |
| link | `#0D7680` | 連結青 |
| band / band-ink | `#262A33` / `#C9C3BC` | 頁尾深帶與其文字 |

深色主題：paper `#201D1A`、ink `#EDE6DC`、ink-soft `#A89F94`、rule `#38342D`、rule-heavy `#EDE6DC`、up `#E5697F`、up-tint `#3B252A`、down `#6FBF8F`、down-tint `#22302A`、link `#5FB8C0`、band `#16130F`。

檔案紙色（僅用於歷史報告檢視頁，即日期早於最新報告的整頁檢視）：淺色 paper 改 `#EAE0CA`、link 同步加深為 `#0B6870`（同色相，維持對比 ≥4.5:1）；深色 paper 改 `#2A251E`（其餘 token 不變）。舊紙色是中性紙色的深淺變化，不是新色相；最新報告與其他頁面一律不得套用。

規則：

- **台股慣例紅漲綠跌，永不反轉**。看漲／看好一律紅系，看跌／看壞一律綠系；中性條目無色標籤（描邊灰）。
- 紅綠是全站唯二的彩色，稀缺性讓「今天誰多誰空」三秒可掃。不引入第三彩色系（連結青只用於連結）。
- 所有文字組合對比 ≥4.5:1（實算最低：ink-soft on paper 5.54）。

## Typography

雙字型制：**Noto Serif TC**（報頭 900、頭條與章節標題 700、股名、詳情正文）×**Noto Sans TC**（UI、標籤、理由、列表、頁尾）。數字一律 `font-variant-numeric: tabular-nums`。正文 ≥14px、標籤 ≥12px。中文引號用「」。生產環境自行部署字型（syralit.toml `[[theme.font_faces]]`）或 Google Fonts；fallback 序列 serif：`"Noto Serif TC", "Songti TC", "PMingLiU", serif`，sans：`"Noto Sans TC", "PingFang TC", "Microsoft JhengHei", sans-serif`。

## Layout

報紙規線系統，max-width 1220px 置中。分區只用線不用卡片，但**規線只到大區塊層級**：3px 粗規線開章節、1px 髮絲線分直欄與大型條目（詳情、歷史列）；小條目（觀察列表、多空條目、產業要聞）之間一律用留白分隔，不畫線——線太密會割眼。報頭置中題字＋3px+1px 雙規線＋dateline。頭版為「總覽主欄＋今日觀察右欄」雙欄；多空判斷四直欄（紅綠欄頭底線）；產業四直欄；詳情兩欄（左股名標籤、右內文）。斷點：1080px 四欄→兩欄、767px 全單欄。

## Elevation & Depth

全平面。零陰影、零光暈、零漸層。層次全部由字級、字重、規線粗細與紅綠底 tint 承擔。

## Shapes

全站零圓角（`border-radius: 0`），含按鈕、標籤、分頁。銳角是刊物氣質的一部分。

## Components

- **section-head**：3px 粗規線＋襯線 24px 標題＋可選 13px 灰註。
- **多空欄（call-col）**：欄頭＝三角形符號（CSS border 畫，非 icon 字型非 emoji）＋標題＋3px 紅／綠底線；欄內條目以留白分隔（不畫線），股名襯線 18px 著紅／綠。
- **多空標籤（tag）**：12.5px 粗體＋tint 底色；中性版無底色、1px rule 描邊、灰字。
- **歷史列表**：日期粗體＋headline 灰字，髮絲線分隔；分頁當前頁 ink 底 paper 字。
- **頁尾深帶**：band 底、免責聲明 13.5px，永遠存在。

## Syralit 對映（實作參考）

`syralit.toml [theme]`：`accent = up`、`background_color = paper`、`secondary_background_color = paper`（同紙面，不做卡面色）、`text_color = ink`、`border_color = rule`、`link_color = link`、`radius = "0px"`、`button_radius = "0px"`、`font = sans 序列`、`heading_font = serif 序列`、`red_color = up`、`green_color = down`、mode = "system" 並提供上表深色 token。版面細節（規線、雙欄、欄頭底線）用自訂樣式實作。

## Do's and Don'ts

- ✅ 新增元件時先問「報紙會怎麼排」：用規線與字級，不加卡片與陰影。
- ✅ 資料數字用 tabular-nums；日期是報告身分，永遠顯眼。
- ❌ 不用 emoji 當 icon；方向符號用 CSS 三角形或 ▲▼ 排版字元。
- ❌ 不用圓角、陰影、漸層、紫色系。
- ❌ 不把紅綠用在多空語意以外的裝飾上；不反轉紅漲綠跌。
- ❌ 不加行銷 CTA、註冊引導——這是刊物，不是 landing page。
