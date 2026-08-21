# daily-brief-skill Specification

## Purpose

TBD - created by archiving change 'daily-report-site'. Update Purpose after archive.

## Requirements

### Requirement: Daily workflow skill exists in the repo

The repo SHALL contain a Cowork skill at skills/daily-brief/SKILL.md (project-root skills directory) defining the daily pre-market workflow as ordered steps: collect industry and stock market news, judge stock calls from the news, compose the report JSON, and POST it to the site API.

#### Scenario: Skill discoverable

- **WHEN** an agent loads skills/daily-brief/SKILL.md in this repo
- **THEN** the daily-brief skill is available with the four workflow steps documented


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
### Requirement: Skill output contract

The skill SHALL embed the full report JSON schema with a filled example and SHALL instruct the agent to produce a payload conforming to it. The skill SHALL require that every stock call carries a one-line reason and at least one news source, and that stocks with major news but no clear direction are emitted as stock_news entries with call none.

#### Scenario: Unclear direction handled

- **WHEN** the agent finds major news about a stock but cannot judge direction
- **THEN** the skill directs it to emit a stock_news entry with call none and no calls-list entry


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
### Requirement: Endpoint configuration via environment

The skill SHALL direct the agent to read the site URL from the BRIEFAST_URL key and the API key from the BRIEFAST_API_KEY key of a .env file at the execution workspace root, in KEY=VALUE format, and SHALL NOT contain a hardcoded URL or key. The skill directory SHALL ship a .env.example template listing the two keys with placeholder values and no real secrets. The skill SHALL instruct the agent to stop and report when the .env file is missing or either key is missing or blank.

#### Scenario: Missing credentials

- **WHEN** the workspace root has no .env file, or BRIEFAST_API_KEY is missing or blank in it, when the skill runs
- **THEN** the workflow stops before collecting and reports the missing configuration

#### Scenario: Template present

- **WHEN** a user prepares a new execution workspace
- **THEN** the skill's .env.example shows the two required keys with placeholder values to copy into the workspace .env


<!-- @trace
source: daily-brief-source-overhaul
updated: 2026-08-10
code:
  - skills/daily-brief/SKILL.md
  - skills/daily-brief/.env.example
-->

---
### Requirement: Failure reporting

The skill SHALL instruct the agent to treat a non-200 API response as a failure, include the response body in its report, and retry once on 5xx errors before reporting failure.

#### Scenario: Validation rejection surfaced

- **WHEN** the API responds 400 with an errors array
- **THEN** the agent fixes the payload per the errors when possible, otherwise reports the errors verbatim

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
### Requirement: Fixed industry taxonomy

The skill SHALL instruct the agent to emit the industries array using exactly four fixed section names — 科技, 金融, 傳產, 房地產 — defined as a closed taxonomy in which every Taiwan-market-relevant industry news item is assigned to exactly one section. Within each section, the skill SHALL require the day's content to be split into events, each carrying its own headline and body. The skill SHALL allow a section to be omitted only after the agent has actively checked that day's sources and found no news for it; omission MUST be a verified decision, not an oversight. The skill SHALL NOT allow section names outside the four fixed names, and SHALL NOT allow a retained section to carry an empty event list.

#### Scenario: News assigned to fixed sections

- **WHEN** the agent composes the industries array from the day's collected news
- **THEN** every industry entry uses one of the four fixed names and carries at least one event, each event having its own headline and body

##### Example: closed-taxonomy assignment

| News topic | Section |
| --- | --- |
| 記憶體價格續漲、封測追單 | 科技 |
| 金管會要求國銀內控清查 | 金融 |
| 貨櫃運價指數反彈 | 傳產 |
| 北台灣預售屋成交降溫 | 房地產 |
| 生技廠取得美國藥證 | 傳產 |

#### Scenario: Empty section omitted after verification

- **WHEN** the agent has checked all sources and found no 房地產 news in the collection window
- **THEN** the 房地產 entry is omitted from industries without placeholder content


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
### Requirement: Analyst perspective guidance

The skill SHALL instruct the agent to write industry sections from an industry analyst's perspective — explaining supply-demand shifts, pricing, and cross-company implications rather than restating headlines — and to write stock_news entries from an equity analyst's perspective — interpreting each news item against the company's situation, including fundamentals context when fresh data is available.

#### Scenario: Industry section interprets rather than restates

- **WHEN** the agent writes an industry section covering a memory price increase
- **THEN** the summary explains the supply-demand driver and which parts of the chain benefit or suffer, instead of only repeating the headline


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
### Requirement: Fundamentals freshness discipline

The skill SHALL require that every fundamentals figure (revenue, EPS, margins, or similar) cited in the report comes from data fetched during the current run from a public source, and that each figure is labeled with its data period. The skill SHALL forbid writing fundamentals figures from model memory. When fresh fundamentals data cannot be fetched, the skill SHALL direct the agent to omit fundamentals context rather than fabricate or approximate it.

#### Scenario: Unfetchable fundamentals omitted

- **WHEN** the agent wants to add revenue context for a stock but cannot fetch current-period data during the run
- **THEN** the stock entry is written from the news alone, with no fundamentals figures

#### Scenario: Cited figure carries its period

- **WHEN** the agent cites a monthly revenue figure fetched from TWSE OpenAPI
- **THEN** the figure appears with its data period, such as 7 月營收


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
### Requirement: Watch points in dedicated fields

The skill SHALL require every stock_news entry and every industries entry to carry a non-empty watch_md field listing one or two forward-looking points the reader can verify later. Each stock-level watch point MUST tie back to the news or fundamentals cited in that entry; each industry-level watch point MUST tie back to the developments reported in that section's events. The skill SHALL forbid duplicating watch points inside event or stock summary bodies and SHALL forbid templated filler watch points that name nothing verifiable.

#### Scenario: Stock watch points populated

- **WHEN** the agent writes a stock_news entry for a company announcing capacity expansion
- **THEN** the entry's watch_md names verifiable follow-ups, such as the timing of new capacity coming online, and summary_md contains no watch-point segment

#### Scenario: Industry watch points populated

- **WHEN** the agent writes the 科技 section covering a memory price increase
- **THEN** the section's watch_md names verifiable follow-ups, such as the next contract-price announcement to watch


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
### Requirement: Source collection completeness gate

The skill SHALL require the agent to record the success or failure of every configured news source during collection and to retry each failed source once. A source that remains unreachable after the retry SHALL NOT stop publication; the skill SHALL direct the agent to publish from what it did collect and to list every missing source in its run report, including the primary source when it is among the failures. Collection mechanics SHALL stay internal: the skill SHALL forbid naming the source roster, the batch structure, or any source outage inside report content such as overview_md or any summary field, so the published report never discloses coverage gaps or how collection works. Per-article citations in the sources field are unaffected. The skill SHALL NOT treat a source that returned successfully with no articles in the window as a failure.

#### Scenario: Primary source failure stays internal

- **WHEN** cnyes remains unreachable after one retry
- **THEN** the agent composes the report from the remaining sources, names the outage only in its run report, and the published report contains no mention of the missing source or of thinner coverage

#### Scenario: Secondary source failure disclosed internally

- **WHEN** two secondary sources remain unreachable after one retry and cnyes succeeded
- **THEN** the agent publishes the report and lists both missing sources in its run report, with nothing about them in report content


<!-- @trace
source: batched-source-collection
updated: 2026-08-13
code:
  - skills/daily-brief/SKILL.md
-->

---
### Requirement: Same-day overwrite guard

The skill SHALL require the agent, before POSTing a report for a date that already has a published report, to compare the new report's industries count and stock_news count against the published version. When either count is lower, the skill SHALL direct the agent to stop, explain why the new report carries less content, and wait for instruction instead of overwriting.

#### Scenario: Shrinking rerun halted

- **WHEN** a rerun for a date produces three industries while the published report for that date has four
- **THEN** the agent stops before POSTing, explains the reduction, and waits for instruction

#### Scenario: Non-shrinking rerun proceeds

- **WHEN** a rerun produces counts equal to or greater than the published report for that date
- **THEN** the agent POSTs the report without an extra confirmation step


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
### Requirement: Taiwan-market coverage boundary

The skill SHALL limit calls and stock_news entries to Taiwan-listed stocks. The skill SHALL direct international market developments and foreign-listed company news into overview_md only, and SHALL NOT allow foreign-listed symbols in calls or stock_news.

#### Scenario: Foreign stock news routed to overview

- **WHEN** the day's news includes a major move in a US-listed company relevant to Taiwan supply chains
- **THEN** the development is summarized in overview_md, and no calls or stock_news entry is created for the foreign symbol

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
### Requirement: One-sentence headlines

The skill SHALL require every industries event and every stock_news entry to carry a non-empty one-sentence headline. Each headline MUST state what happened and why it matters, SHALL NOT be a bare topic label such as a two-word category name, and SHALL NOT restate the first sentence of the body it introduces. The skill SHALL require events within one industry section to be grouped by theme so that one event covers a related cluster of developments rather than one event per article.

#### Scenario: Industry event headline states the development

- **WHEN** the agent writes a 科技 event covering rising memory contract prices
- **THEN** the event headline names the price move and its consequence rather than a label such as 記憶體供需

#### Scenario: Stock headline distinct from body

- **WHEN** the agent writes a stock_news entry whose summary opens by describing a capacity expansion
- **THEN** the entry headline conveys the development and its significance in wording that is not a copy of that opening sentence

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
### Requirement: Pinned source endpoints

The skill SHALL pin every news source to a named list endpoint verified to be reachable and to yield full article text, and SHALL NOT direct the agent to discover sources from portal homepages. The pinned set SHALL be: cnyes category list API with the six categories tw_stock, headline, tech, tw_macro, cnyeshouse, and wd_stock, whose list response already carries full article content and stock tags so no detail-page fetch is made, with wd_stock feeding overview_md only; TWSE OpenAPI datasets t187ap04_L (material information) and t187ap05_L (monthly revenue) pulled whole and filtered by the collection window; CTEE rss_web livenews RSS for the six categories policy, stock, finance, industry, house, and tech, with article bodies fetched over plain HTTP and no headless browser; CNA finance RSS; TechNews site feed; LTN business RSS; and CNBC top stories RSS at www.cnbc.com/id/100003114/device/rss/rss.html, with CNBC article bodies fetched over plain HTTP and subject to the foreign wire scope requirement. The skill SHALL NOT list WSJ Markets as a source, because its article pages sit behind a paywall and feed summaries alone cannot support judgement, and SHALL NOT list UDN as a source.

#### Scenario: Collection uses pinned endpoints

- **WHEN** the agent collects news for the day
- **THEN** every list fetch targets one of the pinned endpoints and cnyes article content is taken from the list response without a detail-page fetch

#### Scenario: No headless browser required

- **WHEN** the agent collects CTEE articles
- **THEN** list discovery uses the CTEE RSS endpoints and article bodies are fetched over plain HTTP


<!-- @trace
source: remove-wsj-source
updated: 2026-08-15
code:
  - skills/daily-brief/SKILL.md
-->

---
### Requirement: TWSE request throttling

The skill SHALL require TWSE OpenAPI requests to run serially with an interval of at least five seconds between requests, and SHALL forbid parallel TWSE fetches. When TWSE responds 429 or refuses the connection, the skill SHALL direct the agent to wait sixty seconds and retry once; when the retry also fails, the agent SHALL treat TWSE as a failed secondary source under the source collection completeness gate.

#### Scenario: Rate limited then recovered

- **WHEN** a TWSE request returns 429 and the retry sixty seconds later succeeds
- **THEN** collection continues normally with no failure disclosure needed for TWSE

#### Scenario: Rate limited twice

- **WHEN** a TWSE request returns 429 and the retry sixty seconds later also fails
- **THEN** the agent proceeds without TWSE data and lists TWSE as a missing source in its run report


<!-- @trace
source: daily-brief-source-overhaul
updated: 2026-08-10
code:
  - skills/daily-brief/SKILL.md
  - skills/daily-brief/.env.example
-->

---
### Requirement: Two-pass ranked judgement

The skill SHALL split judgement into two passes run in order: an industry pass using the industry analyst perspective, then a stock pass using the equity analyst perspective. Before writing in each pass, the skill SHALL require the agent to rank the window's news by importance using that pass's criteria, compared in listed order. Industry-pass criteria: (1) supply-demand or price shifts that already happened, (2) policy or regulatory measures about to take binding effect, (3) breadth of impact across the chain, (4) durability of the shift, (5) certainty of the information. Stock-pass criteria: (1) quantifiable earnings impact already visible, (2) factual tier — company disclosure above company outlook above analyst views above supply-chain chatter above rumor, (3) how soon the impact shows up in results, (4) irreversible change in competitive position, (5) fund-flow signals as corroboration only. The ranking SHALL govern only which events lead a section, how much depth each item receives, and ordering within summaries; the skill SHALL forbid using the ranking to cap or exclude items from the report, and the no-quota rule SHALL remain in force.

#### Scenario: Ranking shapes emphasis not inclusion

- **WHEN** the industry pass ranks a confirmed contract-price increase above an unconfirmed capacity rumor
- **THEN** the price increase leads its section with more depth while the rumor still appears with appropriate hedging, and neither is dropped because of rank

#### Scenario: Same news ranked differently per pass

- **WHEN** a policy announcement ranks high in the industry pass but has no quantifiable earnings impact for a specific stock
- **THEN** the industry section leads with it while the stock pass gives it lower emphasis in the affected entries


<!-- @trace
source: daily-brief-source-overhaul
updated: 2026-08-10
code:
  - skills/daily-brief/SKILL.md
  - skills/daily-brief/.env.example
-->

---
### Requirement: Analyst summary writing craft

The skill SHALL include writing-craft rules for both judgement passes, all bound by two guards: information may come only from the collected news or data fetched during the current run, and when the required material is absent the sentence is omitted entirely rather than replaced with a placeholder such as a cannot-estimate phrase. Industry-pass craft: figures SHALL carry direction, magnitude, period, and a comparison baseline only as stated in the news; price or supply-demand shifts SHALL be attributed to the supply side or the demand side together with the durability implication, and when the news does not state the driver the summary SHALL say the driver is unclear rather than invent one — this is the only rule where missing material is stated instead of omitted; summaries SHALL describe impact distribution as who benefits, who absorbs cost, and who can pass it through; events covered in depth SHALL state whether the shift is cyclical or structural, consistent with the durability ranking criterion. Stock-pass craft: each entry SHALL open with what the news changes about the judgement, one level deeper than the entry headline, and entries with call none SHALL open with why direction cannot be judged; magnitude conversion against company financials SHALL appear only when the base figures were fetched during the current run; expectation-gap statements SHALL rest only on expectations stated in the news or on events already registered in the dedup script; promotional wording SHALL be translated into facts, with unquantified deals downgraded to stated intent; attribution SHALL name sources only as the news names them, keeping vague subjects vague with the certainty tier lowered accordingly. The skill SHALL mark its example table as illustrative, not as sentence templates to copy.

#### Scenario: Unstated driver acknowledged not invented

- **WHEN** the news reports a price rise without stating whether supply cuts or demand pull drives it
- **THEN** the event summary marks the driver as unclear instead of asserting a cause

#### Scenario: Magnitude omitted without placeholder

- **WHEN** a stock entry cites a contract win but no base revenue figure was fetched this run
- **THEN** the entry carries no magnitude comparison and no cannot-estimate placeholder sentence

#### Scenario: None entry opens with the uncertainty

- **WHEN** a stock entry carries call none
- **THEN** its summary opens by stating why the direction cannot be judged from the day's news

<!-- @trace
source: daily-brief-source-overhaul
updated: 2026-08-10
code:
  - skills/daily-brief/SKILL.md
  - skills/daily-brief/.env.example
-->

---
### Requirement: Batched collection with per-batch counts

The skill SHALL split collection into five named batches, each naming the endpoints it covers and each reporting its own counts. Batch one SHALL cover the six cnyes category endpoints and report, per category, how many items were returned and how many of those were unread. Batch two SHALL cover the two TWSE datasets and report, per dataset, how many records fall inside the collection window. Batch three SHALL cover the six CTEE category feeds and report unread counts per category. Batch four SHALL cover the CNA, TechNews, and LTN feeds and report unread counts per source. Batch five SHALL cover the CNBC top stories feed and report its returned and unread counts.

The five batches MAY run in parallel with one another, since they target independent hosts. Within batch two the two TWSE requests SHALL remain serial under the existing interval rule, because that limit applies per host rather than per batch. After all batches return, the skill SHALL require a consolidated table covering every source, and SHALL forbid starting judgement until that table is complete — an absent entry SHALL be treated as an unattempted source, not as an absence of news. A batch whose sources remain unreachable after the retry SHALL be recorded as failed in the consolidated table and handled under the source collection completeness gate: publication proceeds from what was collected, including when the failed batch is the primary source batch or batch five, and every missing source is listed in the run report.

#### Scenario: Counts reported per batch

- **WHEN** the agent runs the five batches concurrently
- **THEN** each batch reports its own per-endpoint counts and judgement waits until all five have reported

#### Scenario: Missing batch blocks judgement

- **WHEN** the consolidated table lacks an entry for the CTEE batch
- **THEN** judgement does not begin and the agent resolves the gap by fetching that batch or recording it as a failed source

#### Scenario: Primary batch failure does not stop the run

- **WHEN** the cnyes batch fails after its retry while the other batches returned successfully
- **THEN** the agent composes and publishes the report from the remaining batches, lists cnyes in its run report, and the published report content discloses nothing about the outage

#### Scenario: Foreign wire batch failure does not stop the run

- **WHEN** the CNBC feed remains unreachable after one retry
- **THEN** the agent publishes from the remaining batches and lists CNBC as a missing source in its run report


<!-- @trace
source: remove-wsj-source
updated: 2026-08-15
code:
  - skills/daily-brief/SKILL.md
-->

---
### Requirement: Foreign wire scope and cross-language deduplication

The skill SHALL restrict foreign wire content (CNBC top stories) to overview_md and to background context inside industry event summaries, and SHALL NOT allow foreign wire items to create or support calls or stock_news entries. Because the similarity check in seen.py compares article bodies within one language and cannot match a Chinese article against an English article covering the same event, the skill SHALL direct the agent to deduplicate foreign wire items by event against the Taiwanese sources collected in the same run: a foreign wire item whose event is already covered by a Taiwanese source SHALL be judged as not reported, recorded via seen.py record with decision skipped, and the Taiwanese version SHALL be the one used. Only foreign wire items covering events absent from the Taiwanese sources SHALL feed report content.

#### Scenario: Event already covered by Taiwanese media

- **WHEN** a CNBC item covers a Fed statement that a cnyes article in the same collection window already reports
- **THEN** the CNBC item is recorded as skipped and overview_md draws on the cnyes version

#### Scenario: Event only in foreign wires

- **WHEN** a CNBC item reports an export-control detail that no Taiwanese source in the window covers
- **THEN** the item informs overview_md or industry event background, and no calls or stock_news entry cites it


<!-- @trace
source: remove-wsj-source
updated: 2026-08-15
code:
  - skills/daily-brief/SKILL.md
-->

---
### Requirement: Previous-close reference data

The skill SHALL direct the agent to fetch previous trading day closing prices during collection, using two plain-HTTP requests: the TWSE listed-stock daily trading dataset exchangeReport/STOCK_DAY_ALL on openapi.twse.com.tw, and the TPEx mainboard daily close quotes dataset openapi/v1/tpex_mainboard_daily_close_quotes on www.tpex.org.tw. The TWSE request SHALL follow the existing TWSE request throttling requirement; the TPEx host is outside that throttle. The agent SHALL validate the trading date embedded in each response (ROC calendar format) against the expected previous trading day in Taipei time, and SHALL treat a response carrying an older date as a failed fetch rather than silently using stale prices. The agent SHALL report, per source, how many records were returned and the embedded data date.

The skill SHALL require the agent to save each response to a file in the working directory and to read prices only by querying those files for the symbols named in the collected news, using file-query tools such as jq or grep. The agent SHALL NOT load a full-market response into the conversation context. Date validation and count reporting SHALL likewise be satisfied by querying the saved files for the embedded date field and the record count, not by reading the responses whole.

Closing prices SHALL serve as judgement context — for example whether news appears already priced in, or as a magnitude baseline — and MAY be quoted in report content only with the price date stated. The skill SHALL forbid creating calls or stock_news entries from price movement alone without supporting news from the collection window. When both price fetches fail after one retry each, the agent SHALL proceed to judge and publish from news alone, list the missing price data in its run report, and disclose nothing about the gap in report content, consistent with the source collection completeness gate's internal-disclosure rule.

#### Scenario: Prices fetched and dated

- **WHEN** the agent fetches both closing-price datasets before judgement on a normal trading morning
- **THEN** each response's embedded date matches the previous Taipei trading day, the agent reports record counts and the data date per source, and judgement proceeds with prices available as context

#### Scenario: Stale date treated as failure

- **WHEN** the TWSE response carries an embedded date older than the previous Taipei trading day after the retry
- **THEN** the agent treats the TWSE price fetch as failed, judges without listed-stock prices, and lists the gap in its run report with nothing disclosed in report content

#### Scenario: No call from price movement alone

- **WHEN** a stock shows a sharp previous-day price move but no news about it exists in the collection window
- **THEN** no calls or stock_news entry is created for that stock, and the price move at most informs context inside entries already backed by news

#### Scenario: Prices queried from saved files by symbol

- **WHEN** the agent needs previous-close prices for the stocks named in the collected news
- **THEN** the agent queries the saved response files for those symbols only, and the full-market payload is never read into the conversation context

##### Example: price quoting with date

| Situation | Allowed in report content | Not allowed |
|-----------|--------------------------|-------------|
| TSMC news entry, price date 2026-08-14 | 昨收（8/14）1,080 元，利多公布前已連漲 | 昨收 1,080 元（無日期標示） |


<!-- @trace
source: previous-close-read-discipline
updated: 2026-08-17
code:
  - skills/daily-brief/SKILL.md
-->

---
### Requirement: Trading-day gate before collection

The skill SHALL determine, after reading `.env` and before any news collection begins, whether the current date in the Asia/Taipei timezone is a Taiwan stock market trading day. On a non-trading day the run SHALL end immediately: no collection, no judgement, no report composition, and no POST. The operator-facing run report SHALL state the closure reason (weekend, or the matching calendar entry name).

The determination MUST be based on verifiable rules, never on model memory or general knowledge: Saturday and Sunday SHALL always be treated as non-trading days; any other date SHALL be checked against the TWSE OpenAPI holiday schedule endpoint `https://openapi.twse.com.tw/v1/holidaySchedule/holidaySchedule`, whose entries carry dates in ROC format (e.g. `1150101` for 2026-01-01).

Calendar entries MUST be classified before use: an entry whose name or description states the market is closed (「市場無交易」, 「放假」, 「補假」) marks a non-trading day; informational entries that mark trading days (「開始交易日」, 「最後交易日」) MUST NOT be treated as closures. A date with no matching closure entry on a weekday SHALL be treated as a trading day.

If the holiday schedule fetch fails on a weekday, the skill SHALL retry once; if it still fails, the skill SHALL proceed with the full run as if the date were a trading day and SHALL list the endpoint failure in the operator-facing run report. The weekend rule requires no endpoint and SHALL apply regardless of endpoint availability. The closure decision and gate mechanics SHALL NOT appear in report content, consistent with the existing rule that collection mechanics stay internal.

#### Scenario: Weekend run ends without publishing

- **WHEN** the skill runs on a Saturday or Sunday in the Asia/Taipei timezone
- **THEN** the run ends before collection, nothing is POSTed, and the run report states the weekend closure — without calling the holiday schedule endpoint

#### Scenario: Holiday matched in calendar ends the run

- **WHEN** the current Taipei date matches a holiday schedule entry whose name or description states the market is closed
- **THEN** the run ends before collection, nothing is POSTed, and the run report cites the entry name as the closure reason

#### Scenario: Informational calendar entry does not close the market

- **WHEN** the current Taipei date matches only an informational entry such as a first-or-last-trading-day notice
- **THEN** the date is treated as a trading day and the run proceeds normally

##### Example: entry classification

| Entry name | Classification |
| --- | --- |
| 市場無交易，僅辦理結算交割作業 | non-trading day |
| 農曆除夕及春節 (description contains 放假) | non-trading day |
| 國曆新年開始交易日 | trading day |
| 農曆春節前最後交易日 | trading day |

#### Scenario: Calendar endpoint failure on a weekday

- **WHEN** the holiday schedule endpoint fails twice (initial attempt plus one retry) on a weekday
- **THEN** the run proceeds as a trading day and the endpoint failure is listed in the operator-facing run report, not in report content

<!-- @trace
source: add-trading-day-gate
updated: 2026-08-16
code:
  - skills/daily-brief/SKILL.md
-->

---
### Requirement: Chip reference data

The skill SHALL direct the agent to fetch previous trading day chip data (institutional investors' net buy/sell and margin trading balances) during collection, using four plain-HTTP requests: the TWSE daily institutional trading dataset rwd/zh/fund/T86 and the TWSE margin trading dataset rwd/zh/marginTrading/MI_MARGN on www.twse.com.tw, both called with the target date as a YYYYMMDD parameter and the all-stocks select type; the TPEx institutional daily trading dataset openapi/v1/tpex_3insti_daily_trading and the TPEx mainboard margin balance dataset openapi/v1/tpex_mainboard_margin_balance on www.tpex.org.tw. Requests to www.twse.com.tw SHALL follow the existing TWSE request throttling requirement as part of the same serial sequence; the TPEx host is outside that throttle.

The agent SHALL validate the trading date embedded in each response (ROC calendar format) against the expected previous trading day in Taipei time, and SHALL treat a response carrying an older date as a failed fetch rather than silently using stale chip data. The agent SHALL report, per source, how many records were returned and the embedded data date.

The skill SHALL require the agent to save each response to a file in the working directory and to read chip values only by querying those files for the symbols named in the report's stock_news entries, using file-query tools such as jq or grep. The agent SHALL NOT load a full-market chip response into the conversation context. Date validation and count reporting SHALL likewise be satisfied by querying the saved files.

When a chip source fails after one retry (TWSE retries under the throttling requirement's sixty-second rule), the agent SHALL proceed to judge and publish without that source's chip data, list the gap in its run report, and disclose nothing about the gap in report content, consistent with the source collection completeness gate's internal-disclosure rule.

#### Scenario: Chip data fetched and dated

- **WHEN** the agent fetches the four chip datasets before judgement on a normal trading morning
- **THEN** each response's embedded date matches the previous Taipei trading day and the agent reports record counts and the data date per source

#### Scenario: Stale date treated as failure

- **WHEN** the TPEx institutional response carries an embedded date older than the previous Taipei trading day after the retry
- **THEN** the agent treats that fetch as failed, composes TPEx-listed entries without institutional chip values, and lists the gap in its run report with nothing disclosed in report content

#### Scenario: TWSE chip requests join the throttle sequence

- **WHEN** the agent fetches the two www.twse.com.tw chip datasets in the same run as other TWSE requests
- **THEN** all TWSE requests run serially with at least five seconds between them and no TWSE fetches run in parallel

#### Scenario: Chip values queried from saved files by symbol

- **WHEN** the agent needs chip values for the stocks in stock_news
- **THEN** the agent queries the saved response files for those symbols only, and no full-market chip payload is read into the conversation context


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
### Requirement: Chip block composition

For every stock_news entry whose market's chip fetches succeeded, the skill SHALL direct the agent to fill the entry's chips block from the saved files of that stock's market (TWSE files for listed stocks, TPEx files for OTC stocks) as follows: date is the validated chip data date in YYYY-MM-DD; foreign_net is the sum of the foreign-investors-excluding-foreign-dealers net column and the foreign-dealers net column, in shares; trust_net and dealer_net are the investment trust and dealer net columns, in shares; total_net is taken directly from the source's three-institutional total column, in shares, not recomputed; margin_change and short_change are the current-day balance minus the previous-day balance from the margin dataset, in trading units (lots).

When a symbol is absent from a chip file (such as a newly listed stock or an instrument without margin trading), the agent SHALL omit the chips block or the affected fields rather than filling zeros or estimates. When a market's chip fetches failed, the agent SHALL omit the chips block for that market's entries entirely.

#### Scenario: Listed stock composed from TWSE files

- **WHEN** a listed stock appears in stock_news and both TWSE chip fetches succeeded with validated dates
- **THEN** its chips block carries the validated date, the summed foreign net, trust net, dealer net, the source's institutional total, and the margin and short balance changes computed as today minus previous day

##### Example: foreign net summation

- **GIVEN** T86 columns for a symbol: foreign-excluding-foreign-dealers net = 54,758,664 shares, foreign-dealers net = 0 shares
- **WHEN** the agent composes the chips block
- **THEN** foreign_net is 54758664

#### Scenario: Symbol absent from a chip file

- **WHEN** a stock_news symbol does not appear in the margin dataset
- **THEN** the entry's chips block omits margin_change and short_change instead of carrying zeros

#### Scenario: Failed market omits chips and still publishes

- **WHEN** both TWSE chip fetches fail after retry while TPEx fetches succeed
- **THEN** listed-stock entries carry no chips block, OTC entries carry theirs, the report is published, and the run report lists the TWSE chip gap


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
### Requirement: Chip citation discipline in judgement

Chip data SHALL serve as corroboration for judgements already grounded in news from the collection window, consistent with the stock-pass ranking criterion that treats fund-flow signals as corroboration only. The skill SHALL forbid creating calls or stock_news entries from chip movement alone without supporting news. When report content cites a chip figure, the citation MUST state the chip data date. Entries whose chips block is omitted SHALL be judged from news and price context alone with no placeholder language about missing chip data.

#### Scenario: No call from chip movement alone

- **WHEN** a stock shows heavy institutional net buying but no news about it exists in the collection window
- **THEN** no calls or stock_news entry is created for that stock

#### Scenario: Chip citation carries its date

- **WHEN** a stock_news summary cites institutional net buying to corroborate a bullish call
- **THEN** the cited figure states the chip data date

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