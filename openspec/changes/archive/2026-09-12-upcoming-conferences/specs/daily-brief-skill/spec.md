## ADDED Requirements

### Requirement: Conference announcement detection

The skill SHALL direct the agent, after the five collection batches complete, to scan the TWSE material-information dataset opendata/t187ap04_L and the TPEx dataset openapi/v1/mopsfin_t187ap04_O already fetched in batch 2 for entries whose 符合條款 field is 第12款, and to parse from each entry's 說明 text the conference date (the value after 1.召開法人說明會之日期：, ROC calendar converted to YYYY-MM-DD), the time (after 2.召開法人說明會之時間：, normalised to HH:MM or empty), the venue (after 3.召開法人說明會之地點：), and the summary message (after 4.法人說明會擇要訊息：). The agent SHALL report how many 第12款 entries were found per market.

The skill SHALL define a company-held conference as an entry whose subject and summary message contain none of the invitation markers 受邀, 應…之邀, 邀請, 參加, 論壇, Summit, Conference, 投資論壇, and whose date field is a single date rather than a range. Entries failing this test SHALL be classified as invited-forum appearances and SHALL NOT be registered. Only company-held conferences whose date is on or after the current Taipei date SHALL be registered. The agent SHALL query GET /api/conferences first and SHALL NOT re-post a brief for a (symbol, held_on) pair already listed in upcoming unless the parsed date, time, or venue differs from the stored values.

#### Scenario: Company-held conference registered

- **WHEN** a 第12款 entry for 1314 reads 說明本公司2026年第二季財務暨營運概況 with date 115/09/16 and no invitation marker
- **THEN** the agent classifies it as company-held and proceeds to compose and post a brief for (1314, 2026-09-16)

#### Scenario: Invited forum skipped

- **WHEN** a 第12款 entry reads 本公司受邀參加瑞銀證券舉辦之投資論壇
- **THEN** the agent classifies it as an invited appearance, registers nothing, and counts it in the run report

##### Example: classification table

| Subject or summary | Date field | Classification |
|--------------------|------------|----------------|
| 說明本公司2026年第二季財務暨營運概況 | 115/09/16 | company-held, register |
| 本公司受邀參加群益證券舉辦之法人說明會 | 115/09/11 | invited, skip |
| 公告本公司將於115年09月11日舉行法人說明會 | 115/09/11 | company-held, register |
| 2026/9/11受邀參加CLSA; 2026/9/16受邀參加KeyBanc | 115/09/11 ~ 115/10/13 | invited and range, skip |

#### Scenario: Already registered conference not re-posted

- **WHEN** GET /api/conferences lists (2330, 2026-10-16) in upcoming with the same date, time, and venue as the parsed entry
- **THEN** the agent posts nothing for that pair


### Requirement: Conference fundamentals reference data

For every conference to be registered, the skill SHALL direct the agent to fetch, during collection, the monthly revenue datasets (TWSE opendata/t187ap05_L, TPEx openapi/v1/mopsfin_t187ap05_O) and the latest quarterly income statement datasets (TWSE opendata/t187ap06_L_ci and its industry variants basi, fh, ins, bd, mim; TPEx openapi/v1/mopsfin_t187ap06_O_ci and the same variants), save each response to a file in the working directory, and read only the rows for the conference symbols with file-query tools, never loading a full-market response into the conversation context. TWSE requests SHALL join the existing TWSE throttling sequence. The agent SHALL read the revenue row's 資料年月, 當月營收, 上月營收, 去年當月營收, 當月累計營收, 去年累計營收 and the income statement row's 年度, 季別, 營業收入, 營業利益（損失）, 本期淨利（淨損）, 基本每股盈餘（元）, and fill the brief's fundamentals object from them, converting ROC year-month to YYYY-MM and labelling the quarter as cumulative (for example 2026Q2 累計) because the dataset reports year-to-date figures.

The skill SHALL state that the income statement dataset carries only the latest cumulative quarter, so a comparison against the prior quarter or the prior year's quarter SHALL be written only when a fetched news article or announcement in the run states those figures, cited in sources; otherwise the quarter comparison sentence SHALL be omitted. The monthly revenue comparisons (month-over-month, year-over-year, year-to-date) SHALL always be stated with direction (變好 or 變壞) in summary_md when the revenue row exists. A symbol absent from a dataset SHALL leave the corresponding fundamentals object omitted with no placeholder, and the omission SHALL be listed in the run report.

#### Scenario: Fundamentals filled from saved files

- **WHEN** the agent registers (1101, 2026-09-20) and the saved revenue file carries 資料年月 11507 with 當月營收 13744103, 去年當月營收 13535929
- **THEN** the brief's fundamentals.revenue has month 2026-07, current 13744103, last_year_month 13535929, and summary_md states the year-over-year direction as 變好

#### Scenario: Quarter comparison only from cited sources

- **WHEN** no fetched article states the prior-year quarter EPS for the conference company
- **THEN** summary_md carries the latest cumulative quarter figures with their label and no prior-year quarter comparison sentence

#### Scenario: Symbol missing from dataset

- **WHEN** the conference symbol has no row in the income statement files
- **THEN** fundamentals.quarter is omitted, the revenue object is still filled if present, and the gap is listed in the run report


### Requirement: Pre-conference brief composition

The skill SHALL embed the full conference brief JSON schema with a filled example and SHALL direct the agent to compose one brief per registered conference and POST it to /api/conference with the Bearer key. The brief's prediction SHALL be bull, bear, or none, judged from the fundamentals direction, the collected news window for that symbol, and chip data as corroboration only; when the agent cannot support a direction from fetched material, prediction SHALL be none and summary_md SHALL open with why no direction is given. summary_md SHALL open with the judgement, then state each revenue comparison with direction, the latest quarter figures with their label, and the news basis with the same citation discipline as stock_news entries (no memory-sourced numbers, every figure dated and sourced). watch_md SHALL list one to three concrete items the conference is expected to clarify. chips SHALL follow the chip block composition rules by market. The agent SHALL confirm a 200 response per posted brief and record every POST outcome in the run report.

The skill SHALL state that a conference brief is independent of the same-day report's calls: a symbol is permitted to carry a short-term call in the report and a different conference prediction, and the brief SHALL name the horizon difference when both exist.

#### Scenario: Brief posted with direction

- **WHEN** revenue comparisons are all positive and the news window carries a dated order announcement for the company
- **THEN** the agent posts a brief with prediction bull, summary_md opening with the judgement and citing the revenue month and the article, and receives 200

#### Scenario: No direction when material is thin

- **WHEN** the company has no news in the window and its revenue row is missing
- **THEN** the agent posts a brief with prediction none whose summary_md opens with the reason no direction is given


### Requirement: Conference settlement after the event

The skill SHALL direct the agent, on every run, to take the pending_settlement array from GET /api/conferences and, for each entry whose held_on is at least one trading day before the previous Taipei trading day, fetch that stock's daily trading history: for twse the endpoint rwd/zh/afterTrading/STOCK_DAY on www.twse.com.tw with date set to the first day of the month containing held_on and stockNo set to the symbol (joining the TWSE throttling sequence), and for tpex the endpoint www/zh-tw/afterTrading/tradingStock on www.tpex.org.tw with code set to the symbol and date set to the first day of that month; when the day after held_on falls in the next month the agent SHALL fetch that month too. The agent SHALL take pre_close as the closing price of the last trading day strictly before held_on and post_close as the closing price of the first trading day strictly after held_on, both with their dates, and POST them to /api/conference/settle. When the post-conference trading day is not yet present in the response, the entry SHALL be left pending and reported. The agent SHALL record each settlement response and outcome in the run report and SHALL NOT mention settlement mechanics in report content.

#### Scenario: Settlement posted two trading days after the conference

- **WHEN** a conference was held on 2026-10-16 (Friday), today is 2026-10-20 (Tuesday), and the history carries closes for 2026-10-15 and 2026-10-19
- **THEN** the agent posts pre_close from 2026-10-15 and post_close from 2026-10-19 and records the returned outcome

#### Scenario: Post-conference close not yet available

- **WHEN** a conference was held yesterday and the history ends on the conference date
- **THEN** the agent posts no settlement, leaves the entry pending, and lists it in the run report

##### Example: close selection across a month boundary

| held_on | Trading days around it | pre_close_date | post_close_date | Months fetched |
|---------|------------------------|----------------|-----------------|----------------|
| 2026-10-16 | 10-15, 10-16, 10-19 | 2026-10-15 | 2026-10-19 | 2026-10 |
| 2026-09-30 | 09-29, 09-30, 10-01 | 2026-09-29 | 2026-10-01 | 2026-09 and 2026-10 |
