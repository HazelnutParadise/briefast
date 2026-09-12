# report-viewing Specification

## Purpose

TBD - created by archiving change 'daily-report-site'. Update Purpose after archive.

## Requirements

### Requirement: Homepage shows the latest report

The homepage SHALL render the report with the newest date found in the report store, with the masthead showing the report date and generated-at time. When no report exists, the homepage SHALL show an empty state ("尚無報告") and SHALL NOT crash.

#### Scenario: Latest report rendered

- **WHEN** reports exist for 2026-08-06 and 2026-08-07 and a visitor opens the homepage
- **THEN** the 2026-08-07 report is rendered and the masthead shows that date

#### Scenario: Empty state

- **WHEN** no report file exists and a visitor opens the homepage
- **THEN** an empty-state message is shown and the page renders without error


<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
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


<!-- @trace
source: upcoming-conferences
updated: 2026-09-12
code:
  - internal/store/migration_v2.sql
  - internal/site/styles.go
  - internal/seo/meta.go
  - main.go
  - DESIGN.md
  - AGENTS.md
  - internal/seo/endpoints.go
  - internal/site/conference.go
  - internal/api/conference.go
  - internal/report/conference.go
  - internal/site/site.go
  - internal/store/conference.go
  - internal/store/store.go
  - skills/daily-brief/SKILL.md
  - docs/plans/2026-09-12-main-upcoming-conferences-plan.md
tests:
  - internal/site/conference_test.go
  - internal/store/conference_test.go
  - main_test.go
  - internal/report/conference_test.go
  - internal/site/site_test.go
  - internal/seo/seo_test.go
  - internal/api/conference_test.go
  - internal/report/skill_example_test.go
  - internal/store/store_test.go
-->

---
### Requirement: Stock calls display

The stock calls section SHALL group entries into four cards: short-term bullish, short-term bearish, long-term bullish, long-term bearish. Bullish cards SHALL use red accents and bearish cards SHALL use green accents (Taiwan market convention). Each entry SHALL show the stock name, symbol, and its one-line reason.

#### Scenario: Calls grouped into four cards

- **WHEN** a report contains entries in all four call lists
- **THEN** each entry appears in its matching card with name, symbol, and reason


<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
### Requirement: Stock news detail display

The stock news detail section SHALL render one entry per stock_news item: stock name and symbol, the entry headline, a call tag matching its call value, the summary markdown, and its news source links. The headline SHALL render in the entry's identification area alongside the stock name rather than inside the body text. Entries with call value none SHALL show no tag.

#### Scenario: Tagged entry

- **WHEN** a stock_news item has call short_bull
- **THEN** its entry shows a short-term bullish tag

#### Scenario: Neutral entry without tag

- **WHEN** a stock_news item has call none
- **THEN** its entry shows the stock name and summary with no call tag

#### Scenario: Headline shown with the stock identity

- **WHEN** a stock_news item carries a headline
- **THEN** that headline appears in the entry's identification area with the stock name and symbol


<!-- @trace
source: news-headlines
updated: 2026-08-09
code:
  - internal/api/report.go
  - internal/site/site.go
  - skills/daily-brief/SKILL.md
  - internal/site/styles.go
  - internal/report/schema.go
  - internal/api/read.go
  - main.go
tests:
  - internal/api/report_test.go
  - internal/report/schema_test.go
  - internal/api/read_test.go
  - internal/site/site_test.go
  - main_test.go
-->

---
### Requirement: History browsing

The history page SHALL list report dates from newest to oldest, 10 per page, each row showing the date and the report headline. Selecting a date SHALL render that day's report using the same layout as the homepage.

When the selected date is older than the newest stored report date, the rendered report SHALL display a prominent historical-report notice below the masthead that names the report date, states the content is not the latest report, and links to both the homepage (latest report) and the history list. When the selected date equals the newest stored report date, the notice SHALL NOT appear.

Every page under the history route — the history list, a report view of any date, and the not-found state — SHALL render on the distinct archive paper background, applied by overriding the paper color token, with a dedicated value for the light theme and for the dark theme. The archive background marks the history section itself, independent of whether the viewed report is the newest. Every existing text-on-paper combination MUST keep a contrast ratio of at least 4.5:1 on the archive background; when a text token falls below that ratio on the archive paper, the override SHALL darken or lighten that token while keeping its hue, and SHALL NOT introduce any new hue. The homepage SHALL NOT use the archive background.

A historical report view SHALL offer masthead navigation to both the homepage and the history list. The not-found state for a missing date SHALL also link back to the homepage in addition to the history list. Determining the newest stored report date MUST NOT require loading the full report payload.

#### Scenario: History list ordering and paging

- **WHEN** 25 reports exist and a visitor opens the history page
- **THEN** the 10 newest dates are listed in descending order with pagination to reach the rest

##### Example: ordering

- **GIVEN** reports dated 2026-08-05, 2026-08-07, 2026-08-06
- **WHEN** the history page renders
- **THEN** rows appear in order: 2026-08-07, 2026-08-06, 2026-08-05

#### Scenario: View a historical report

- **WHEN** a visitor selects 2026-08-05 from the history list and a newer report exists
- **THEN** the 2026-08-05 report renders with the same five-section layout as the homepage, plus a historical-report notice naming 2026-08-05 with links to the homepage and the history list

#### Scenario: Whole history section uses archive paper background

- **WHEN** a visitor opens the history list, any dated report view, or a history URL whose date has no report
- **THEN** the page overrides the paper token to the archive paper values for both light and dark themes

#### Scenario: Newest report opened from history carries no notice

- **WHEN** a visitor selects the newest stored date from the history list
- **THEN** the report renders on the archive background but without the historical-report notice

##### Example: notice and background by page

| Page | Archive background | Historical notice |
| --- | --- | --- |
| Homepage (latest report) | no | no |
| History list | yes | no |
| History view of an older date | yes | yes |
| History view of the newest date | yes | no |
| History view of a missing date | yes | no |

#### Scenario: Missing date links home

- **WHEN** a visitor opens a history URL whose date has no stored report
- **THEN** the not-found state offers links to both the homepage and the history list


<!-- @trace
source: history-route-archive-background
updated: 2026-08-16
code:
  - DESIGN.md
  - internal/site/site.go
  - internal/site/styles.go
tests:
  - internal/site/site_test.go
-->

---
### Requirement: Disclaimer footer

Every public report view SHALL display a fixed disclaimer footer stating the content is AI-generated from public news, is for reference only, and is not investment advice. The footer text SHALL be defined in application code, not in report payloads.

#### Scenario: Footer always present

- **WHEN** any report page (homepage or historical) is rendered
- **THEN** the disclaimer footer is displayed

<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
### Requirement: Watch points display

The industry news summary section and the stock news detail section SHALL render each entry's watch_md as a visually distinct block labeled 觀察重點, placed after the entry's summary content, with styling defined in application code. When an entry has no watch_md value or an empty one — such as reports stored before the field existed — the block SHALL be omitted entirely, with no placeholder text and no rendering error.

#### Scenario: Watch points block rendered

- **WHEN** a report is rendered and a stock_news entry carries a non-empty watch_md
- **THEN** the entry shows a 觀察重點 block after its summary content

#### Scenario: Legacy report without the field

- **WHEN** a report stored before the watch_md field existed is rendered
- **THEN** industry and stock entries render without a 觀察重點 block and without placeholder content

<!-- @trace
source: daily-brief-analyst-guidance
updated: 2026-08-09
code:
  - internal/api/read.go
  - internal/report/schema.go
  - skills/daily-brief/SKILL.md
  - internal/api/report.go
  - main.go
  - internal/site/styles.go
  - internal/site/site.go
tests:
  - internal/api/report_test.go
  - main_test.go
  - internal/site/site_test.go
  - internal/api/read_test.go
  - internal/report/schema_test.go
-->

---
### Requirement: Industry event headlines display

The industry news summary section SHALL render each event of a section as its own unit: the event headline as a heading styled distinctly from body text, followed by that event's body markdown. Multiple events within one section SHALL be visually separable from each other. Headline styling SHALL be defined in application code and SHALL NOT be alterable by report payload content.

#### Scenario: Two events in one section

- **WHEN** the 科技 section carries two events
- **THEN** both event headlines render as headings above their own bodies and the two events read as separate units

<!-- @trace
source: news-headlines
updated: 2026-08-09
code:
  - internal/api/report.go
  - internal/site/site.go
  - skills/daily-brief/SKILL.md
  - internal/site/styles.go
  - internal/report/schema.go
  - internal/api/read.go
  - main.go
tests:
  - internal/api/report_test.go
  - internal/report/schema_test.go
  - internal/api/read_test.go
  - internal/site/site_test.go
  - main_test.go
-->

---
### Requirement: Chip data display

For every stock_news entry that carries a chips object, the stock news detail section SHALL render a chip block labeled 籌碼面 after the entry's watch points block (or after the summary content when no watch points block renders) and before the source links. The block SHALL contain three horizontal bar rows for foreign investors, investment trust, and dealers. Bar widths SHALL be computed server-side, scaled proportionally to the largest absolute net value among the three rows; a zero net value SHALL render its row with a zero-width bar. Net-buy bars SHALL use the bullish red accent and net-sell bars SHALL use the bearish green accent (Taiwan market convention), with no rounded corners and no shadows, consistent with the established layout language. Each row SHALL show its net value converted from shares to lots (divided by 1,000 and rounded to the nearest integer) with an explicit plus sign for positive values.

Below the bars, the block SHALL render one summary line stating the three-institutional total (in lots, same conversion), the margin balance change and short balance change (in lots, from the chips object's trading-unit values), and the chip data date. When margin_change or short_change is absent from the chips object, the corresponding fragment SHALL be omitted with no placeholder text.

When an entry has no chips object — such as reports stored before the field existed or entries whose chip data was unavailable — the chip block SHALL be omitted entirely, with no placeholder text and no rendering error. Chip block styling SHALL be defined in application code and SHALL NOT be alterable by report payload content.

#### Scenario: Chip block rendered with scaled bars

- **WHEN** a stock_news entry carries chips with foreign_net 54,758,664, trust_net -15,000, and dealer_net 2,063,215
- **THEN** the entry shows a 籌碼面 block where the foreign row renders the widest bar in the bullish red accent, the trust row renders a near-zero-width bar in the bearish green accent, and the dealer row renders a proportionally narrower bar in the bullish red accent

##### Example: bar scaling and lot conversion

| Row | Net (shares) | Bar direction and color | Relative width | Shown value |
|------|--------------|-------------------------|----------------|-------------|
| 外資 | 54,758,664 | net buy, bullish red | 100% | +54,759 張 |
| 投信 | -15,000 | net sell, bearish green | ~0.03% | -15 張 |
| 自營商 | 2,063,215 | net buy, bullish red | ~3.8% | +2,063 張 |

#### Scenario: Summary line with margin data

- **WHEN** an entry's chips carries total_net 56,806,879 shares, margin_change -18,270 lots, short_change 9,056 lots, and date 2026-08-20
- **THEN** the summary line states the total as +56,807 張, 融資 -18,270 張, 融券 +9,056 張, and the data date

#### Scenario: Missing margin fragment omitted

- **WHEN** an entry's chips carries institutional values but no margin_change and no short_change
- **THEN** the summary line shows the institutional total and the data date with no margin or short fragment and no placeholder text

#### Scenario: Legacy entry without chips

- **WHEN** a report stored before the chips field existed is rendered
- **THEN** stock news entries render without a 籌碼面 block and without placeholder content

<!-- @trace
source: add-stock-chip-data
updated: 2026-08-21
code:
  - internal/site/styles.go
  - internal/site/site.go
  - skills/daily-brief/SKILL.md
  - internal/report/schema.go
tests:
  - internal/site/site_test.go
  - internal/api/report_test.go
  - internal/report/schema_test.go
-->

---
### Requirement: Advertising slot

Report pages SHALL carry one advertising slot placed between the stock calls section and the industry news section. The slot SHALL be visually delimited from editorial content by a rule and a label reading 廣告, because the advertisement itself is rendered by a third party whose appearance the application does not control.

The slot SHALL consist of a mount point rendered as part of the page markup and a loader script delivered through a node that executes scripts in the main document. The loader script SHALL NOT re-execute when the page re-renders in response to a background update, so that a refresh of the report does not count an additional impression.

Pages without report content SHALL NOT carry the slot.

#### Scenario: Slot appears between calls and industries

- **WHEN** a client requests a page showing a report
- **THEN** the rendered markup contains the advertising mount point after the stock calls section and before the industry news section

#### Scenario: Slot is labelled and delimited

- **WHEN** the advertising slot is rendered
- **THEN** it carries a label reading 廣告 and a rule separating it from the surrounding editorial content

#### Scenario: Loader runs in the main document

- **WHEN** a report page is rendered in a browser
- **THEN** the loader script executes and resolves the mount point that the page markup provided

#### Scenario: Background update does not re-run the loader

- **WHEN** a new report is published and the open page re-renders through the live update channel
- **THEN** the advertising node is reused and its loader script does not run again

#### Scenario: Pages without a report carry no slot

- **WHEN** a client requests the history listing, or a dated view for which no report exists
- **THEN** the rendered markup contains no advertising mount point

##### Example: slot presence by page

| Page | Advertising slot |
|---|---|
| home page showing the latest report | present |
| dated view of a stored report | present |
| history listing | absent |
| dated view with no stored report | absent |
| waiting-for-report home page | absent |

<!-- @trace
source: article-ad-slot
updated: 2026-08-25
code:
  - main.go
  - DESIGN.md
  - internal/site/styles.go
  - AGENTS.md
  - internal/seo/meta.go
  - internal/seo/middleware.go
  - internal/site/site.go
tests:
  - internal/seo/seo_test.go
  - internal/site/site_test.go
-->

---
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


<!-- @trace
source: upcoming-conferences
updated: 2026-09-12
code:
  - internal/store/migration_v2.sql
  - internal/site/styles.go
  - internal/seo/meta.go
  - main.go
  - DESIGN.md
  - AGENTS.md
  - internal/seo/endpoints.go
  - internal/site/conference.go
  - internal/api/conference.go
  - internal/report/conference.go
  - internal/site/site.go
  - internal/store/conference.go
  - internal/store/store.go
  - skills/daily-brief/SKILL.md
  - docs/plans/2026-09-12-main-upcoming-conferences-plan.md
tests:
  - internal/site/conference_test.go
  - internal/store/conference_test.go
  - main_test.go
  - internal/report/conference_test.go
  - internal/site/site_test.go
  - internal/seo/seo_test.go
  - internal/api/conference_test.go
  - internal/report/skill_example_test.go
  - internal/store/store_test.go
-->

---
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


<!-- @trace
source: upcoming-conferences
updated: 2026-09-12
code:
  - internal/store/migration_v2.sql
  - internal/site/styles.go
  - internal/seo/meta.go
  - main.go
  - DESIGN.md
  - AGENTS.md
  - internal/seo/endpoints.go
  - internal/site/conference.go
  - internal/api/conference.go
  - internal/report/conference.go
  - internal/site/site.go
  - internal/store/conference.go
  - internal/store/store.go
  - skills/daily-brief/SKILL.md
  - docs/plans/2026-09-12-main-upcoming-conferences-plan.md
tests:
  - internal/site/conference_test.go
  - internal/store/conference_test.go
  - main_test.go
  - internal/report/conference_test.go
  - internal/site/site_test.go
  - internal/seo/seo_test.go
  - internal/api/conference_test.go
  - internal/report/skill_example_test.go
  - internal/store/store_test.go
-->

---
### Requirement: Live update on conference changes

When a conference brief is ingested or settled, connected homepage sessions SHALL re-render so the conferences section reflects the new state without a manual reload, using the same shared version signal as report ingestion.

#### Scenario: Open homepage shows new card

- **WHEN** a visitor has the homepage open and a new upcoming conference brief is ingested
- **THEN** the conferences section appears or gains the new card without a reload

<!-- @trace
source: upcoming-conferences
updated: 2026-09-12
code:
  - internal/store/migration_v2.sql
  - internal/site/styles.go
  - internal/seo/meta.go
  - main.go
  - DESIGN.md
  - AGENTS.md
  - internal/seo/endpoints.go
  - internal/site/conference.go
  - internal/api/conference.go
  - internal/report/conference.go
  - internal/site/site.go
  - internal/store/conference.go
  - internal/store/store.go
  - skills/daily-brief/SKILL.md
  - docs/plans/2026-09-12-main-upcoming-conferences-plan.md
tests:
  - internal/site/conference_test.go
  - internal/store/conference_test.go
  - main_test.go
  - internal/report/conference_test.go
  - internal/site/site_test.go
  - internal/seo/seo_test.go
  - internal/api/conference_test.go
  - internal/report/skill_example_test.go
  - internal/store/store_test.go
-->