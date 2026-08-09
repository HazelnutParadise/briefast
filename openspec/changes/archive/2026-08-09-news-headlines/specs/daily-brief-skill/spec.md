## ADDED Requirements

### Requirement: One-sentence headlines

The skill SHALL require every industries event and every stock_news entry to carry a non-empty one-sentence headline. Each headline MUST state what happened and why it matters, SHALL NOT be a bare topic label such as a two-word category name, and SHALL NOT restate the first sentence of the body it introduces. The skill SHALL require events within one industry section to be grouped by theme so that one event covers a related cluster of developments rather than one event per article.

#### Scenario: Industry event headline states the development

- **WHEN** the agent writes a 科技 event covering rising memory contract prices
- **THEN** the event headline names the price move and its consequence rather than a label such as 記憶體供需

#### Scenario: Stock headline distinct from body

- **WHEN** the agent writes a stock_news entry whose summary opens by describing a capacity expansion
- **THEN** the entry headline conveys the development and its significance in wording that is not a copy of that opening sentence

## MODIFIED Requirements

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

### Requirement: Watch points in dedicated fields

The skill SHALL require every stock_news entry and every industries entry to carry a non-empty watch_md field listing one or two forward-looking points the reader can verify later. Each stock-level watch point MUST tie back to the news or fundamentals cited in that entry; each industry-level watch point MUST tie back to the developments reported in that section's events. The skill SHALL forbid duplicating watch points inside event or stock summary bodies and SHALL forbid templated filler watch points that name nothing verifiable.

#### Scenario: Stock watch points populated

- **WHEN** the agent writes a stock_news entry for a company announcing capacity expansion
- **THEN** the entry's watch_md names verifiable follow-ups, such as the timing of new capacity coming online, and summary_md contains no watch-point segment

#### Scenario: Industry watch points populated

- **WHEN** the agent writes the 科技 section covering a memory price increase
- **THEN** the section's watch_md names verifiable follow-ups, such as the next contract-price announcement to watch
