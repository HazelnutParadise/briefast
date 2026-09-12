## MODIFIED Requirements

### Requirement: Fixed report layout

The report page SHALL render exactly five fixed content sections in this order: pre-market overview, today's watch, stock calls, industry news summary, stock news detail. On the homepage only, an upcoming-conferences section SHALL be inserted between the stock calls section and the advertising slot whenever at least one stored conference has held_on on or after the current Asia/Taipei date; when none exists the section SHALL be omitted entirely. Historical report views SHALL NOT render the conferences section. Section headers, ordering, and styling SHALL be defined in application code and SHALL NOT be alterable by report payload or conference payload content.

#### Scenario: Section order independent of payload

- **WHEN** a report is rendered
- **THEN** the five fixed sections appear in the fixed order regardless of field order in the stored JSON

#### Scenario: Conferences section placement on the homepage

- **WHEN** the homepage renders and an upcoming conference exists
- **THEN** the conferences section appears after the stock calls section and before the advertising slot, and the five fixed sections keep their order

#### Scenario: Homepage without upcoming conferences

- **WHEN** the homepage renders and no conference has held_on on or after today
- **THEN** no conferences section and no placeholder text is rendered

#### Scenario: Historical view carries no conferences section

- **WHEN** a dated report view renders while upcoming conferences exist
- **THEN** the page contains no conferences section

## ADDED Requirements

### Requirement: Upcoming conferences section

The upcoming-conferences section SHALL be headed 近期法說會 with a note stating that predictions are the model's pre-conference view. It SHALL render one card per upcoming conference ordered by held_on ascending then symbol, each card showing the company name and symbol, the conference date with weekday and time when held_at is non-empty, a prediction tag, the brief headline, and a link labelled 會前報告 → to /conference/?symbol=<symbol>&date=<held_on>. The prediction tag SHALL read 會後看漲 in the bullish red accent for bull, 會後看跌 in the bearish green accent for bear, and 無法判斷 in the neutral outlined style for none (Taiwan market convention, never reversed). Cards SHALL use the four-column rule layout of the stock calls section with no rounded corners and no shadows, collapsing to two columns at 1080px and one column at 767px.

When at least one settled conference with a bull or bear prediction exists, the section head SHALL show an accuracy line reading 近 N 場命中 X 場, where N is the number of settled bull-or-bear conferences counted from the most recent 20 by held_on and X is how many of those have outcome hit; settled conferences with prediction none SHALL NOT count. When no such conference exists the accuracy line SHALL be omitted with no placeholder. Below the cards, the section SHALL list up to five most recently settled conferences (by held_on descending), each showing company name and symbol, held_on, the prediction tag, the actual close-to-close change as a signed percentage colored red when positive and green when negative, and 命中 or 落空 (or nothing for prediction none). The settled list SHALL be omitted when no settled conference exists.

#### Scenario: Cards with prediction tags

- **WHEN** upcoming conferences exist for 2330 (bull, 2026-10-16 14:00) and 1314 (none, 2026-10-14)
- **THEN** the section shows the 1314 card first with a neutral 無法判斷 tag and the 2330 card second with a red 會後看漲 tag, each linking to its brief page

#### Scenario: Accuracy line computed from settled outcomes

- **WHEN** 7 settled conferences exist: 5 with bull or bear predictions of which 3 are hit, and 2 with prediction none
- **THEN** the section head reads 近 5 場命中 3 場

##### Example: accuracy counting

| Settled conferences (most recent first) | Accuracy line |
|------------------------------------------|---------------|
| bull hit, bear miss, none, bull hit | 近 3 場命中 2 場 |
| 25 bull-or-bear settled, 14 hit among the newest 20 | 近 20 場命中 14 場 |
| only prediction none settled | (omitted) |

#### Scenario: Recently settled list

- **WHEN** a conference for 2330 settled with pre_close 1,080 and post_close 1,095 under prediction bull
- **THEN** the settled list shows 2330 with a red +1.39% change and 命中

#### Scenario: No settled conferences

- **WHEN** upcoming conferences exist but none has been settled
- **THEN** the section renders the cards, no accuracy line, and no settled list


### Requirement: Conference brief page

The site SHALL serve /conference/ taking symbol and date query parameters, rendering the stored brief for that (symbol, held_on) pair on the standard paper background with the masthead, navigation back to the homepage, and the fixed disclaimer footer. The page SHALL show: company name and symbol; conference date with weekday, time, and venue; the prediction tag; the brief headline; the summary markdown under a 判斷依據 heading; a 基本面 block; the watch points block labelled 觀察重點; the chip block when chips is present (same rendering as stock news entries); source links; and a line stating generated_at and announced_on.

The 基本面 block SHALL render, when the fundamentals.revenue object is present, a table with rows for 當月營收, 上月營收, 去年同月營收, 今年累計營收, and 去年累計營收 in thousands of NTD grouped by thousands, plus computed percentage changes month-over-month, year-over-year, and year-to-date-over-prior-year, each signed and colored red when positive and green when negative. When the fundamentals.quarter object is present it SHALL render the label, revenue, operating income, net income, and EPS. When fundamentals is absent the block SHALL be omitted with no placeholder.

When the brief has been settled, the page SHALL show a result band under the masthead stating pre_close with its date, post_close with its date, the signed percentage change, and 命中 or 落空 (or 未列入統計 for prediction none). When the pair is not stored, the page SHALL render a not-found state with links to the homepage, and SHALL NOT crash. A missing or malformed symbol or date parameter SHALL render the same not-found state.

#### Scenario: Brief page renders all blocks

- **WHEN** a visitor opens /conference/?symbol=2330&date=2026-10-16 and the stored brief carries fundamentals with revenue and quarter, chips, and three sources
- **THEN** the page shows the prediction tag, 判斷依據, a 基本面 table with signed percentage changes, 觀察重點, the 籌碼面 block, three source links, and the disclaimer footer

##### Example: revenue percentage rendering

| Row | Current | Comparison | Rendered |
|-----|---------|------------|----------|
| 上月比較 | 13,744,103 | 13,382,706 | +2.70% in red |
| 去年同月 | 13,744,103 | 13,535,929 | +1.54% in red |
| 去年累計 | 85,211,435 | 90,000,000 | -5.32% in green |

#### Scenario: Settled brief shows result band

- **WHEN** the brief for (2330, 2026-10-16) has pre_close 1,080 on 2026-10-15 and post_close 1,095 on 2026-10-17 under prediction bull
- **THEN** the page shows a result band with both closes and dates, +1.39%, and 命中

#### Scenario: Unknown conference

- **WHEN** a visitor opens /conference/?symbol=9999&date=2026-01-01 with no stored brief
- **THEN** the page renders a not-found state linking to the homepage and no error


### Requirement: Live update on conference changes

When a conference brief is ingested or settled, connected homepage sessions SHALL re-render so the conferences section reflects the new state without a manual reload, using the same shared version signal as report ingestion.

#### Scenario: Open homepage shows new card

- **WHEN** a visitor has the homepage open and a new upcoming conference brief is ingested
- **THEN** the conferences section appears or gains the new card without a reload
