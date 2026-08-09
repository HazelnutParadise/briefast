## ADDED Requirements

### Requirement: Fixed industry taxonomy

The skill SHALL instruct the agent to emit the industries array using exactly four fixed section names — 科技, 金融, 傳產, 房地產 — defined as a closed taxonomy in which every Taiwan-market-relevant industry news item is assigned to exactly one section. Within each section, the skill SHALL require the day's content to be grouped under bold dynamic sub-theme lead-ins inside summary_md. The skill SHALL allow a section to be omitted only after the agent has actively checked that day's sources and found no news for it; omission MUST be a verified decision, not an oversight. The skill SHALL NOT allow section names outside the four fixed names.

#### Scenario: News assigned to fixed sections

- **WHEN** the agent composes the industries array from the day's collected news
- **THEN** every industry entry uses one of the four fixed names, and each entry's summary_md groups items under bold sub-theme lead-ins

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

### Requirement: Analyst perspective guidance

The skill SHALL instruct the agent to write industry sections from an industry analyst's perspective — explaining supply-demand shifts, pricing, and cross-company implications rather than restating headlines — and to write stock_news entries from an equity analyst's perspective — interpreting each news item against the company's situation, including fundamentals context when fresh data is available.

#### Scenario: Industry section interprets rather than restates

- **WHEN** the agent writes an industry section covering a memory price increase
- **THEN** the summary explains the supply-demand driver and which parts of the chain benefit or suffer, instead of only repeating the headline

### Requirement: Fundamentals freshness discipline

The skill SHALL require that every fundamentals figure (revenue, EPS, margins, or similar) cited in the report comes from data fetched during the current run from a public source, and that each figure is labeled with its data period. The skill SHALL forbid writing fundamentals figures from model memory. When fresh fundamentals data cannot be fetched, the skill SHALL direct the agent to omit fundamentals context rather than fabricate or approximate it.

#### Scenario: Unfetchable fundamentals omitted

- **WHEN** the agent wants to add revenue context for a stock but cannot fetch current-period data during the run
- **THEN** the stock entry is written from the news alone, with no fundamentals figures

#### Scenario: Cited figure carries its period

- **WHEN** the agent cites a monthly revenue figure fetched from TWSE OpenAPI
- **THEN** the figure appears with its data period, such as 7 月營收

### Requirement: Watch points in dedicated fields

The skill SHALL require every stock_news entry and every industries entry to carry a non-empty watch_md field listing one or two forward-looking points the reader can verify later. Each stock-level watch point MUST tie back to the news or fundamentals cited in that entry; each industry-level watch point MUST tie back to that section's reported developments. The skill SHALL forbid duplicating watch points inside summary_md and SHALL forbid templated filler watch points that name nothing verifiable.

#### Scenario: Stock watch points populated

- **WHEN** the agent writes a stock_news entry for a company announcing capacity expansion
- **THEN** the entry's watch_md names verifiable follow-ups, such as the timing of new capacity coming online, and summary_md contains no watch-point segment

#### Scenario: Industry watch points populated

- **WHEN** the agent writes the 科技 section covering a memory price increase
- **THEN** the section's watch_md names verifiable follow-ups, such as the next contract-price announcement to watch

### Requirement: Source collection completeness gate

The skill SHALL require the agent to record the success or failure of every configured news source during collection and to retry each failed source once. When the primary source cnyes still fails after the retry, the skill SHALL direct the agent to stop before composing the report, publish nothing, and report which source failed and why. When only secondary sources fail after the retry, the skill SHALL allow publication and SHALL require the agent to list every missing source in its run report.

#### Scenario: Primary source failure blocks publication

- **WHEN** cnyes remains unreachable after one retry
- **THEN** the agent stops without composing or POSTing a report and reports the failed source

#### Scenario: Secondary source failure disclosed

- **WHEN** two secondary sources remain unreachable after one retry and cnyes succeeded
- **THEN** the agent publishes the report and lists both missing sources in its run report

### Requirement: Same-day overwrite guard

The skill SHALL require the agent, before POSTing a report for a date that already has a published report, to compare the new report's industries count and stock_news count against the published version. When either count is lower, the skill SHALL direct the agent to stop, explain why the new report carries less content, and wait for instruction instead of overwriting.

#### Scenario: Shrinking rerun halted

- **WHEN** a rerun for a date produces three industries while the published report for that date has four
- **THEN** the agent stops before POSTing, explains the reduction, and waits for instruction

#### Scenario: Non-shrinking rerun proceeds

- **WHEN** a rerun produces counts equal to or greater than the published report for that date
- **THEN** the agent POSTs the report without an extra confirmation step

### Requirement: Taiwan-market coverage boundary

The skill SHALL limit calls and stock_news entries to Taiwan-listed stocks. The skill SHALL direct international market developments and foreign-listed company news into overview_md only, and SHALL NOT allow foreign-listed symbols in calls or stock_news.

#### Scenario: Foreign stock news routed to overview

- **WHEN** the day's news includes a major move in a US-listed company relevant to Taiwan supply chains
- **THEN** the development is summarized in overview_md, and no calls or stock_news entry is created for the foreign symbol
