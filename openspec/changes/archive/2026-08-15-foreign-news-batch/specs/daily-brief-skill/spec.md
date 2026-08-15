## ADDED Requirements

### Requirement: Foreign wire scope and cross-language deduplication

The skill SHALL restrict foreign wire content (WSJ Markets and CNBC top stories) to overview_md and to background context inside industry event summaries, and SHALL NOT allow foreign wire items to create or support calls or stock_news entries. Because the similarity check in seen.py compares article bodies within one language and cannot match a Chinese article against an English article covering the same event, the skill SHALL direct the agent to deduplicate foreign wire items by event against the Taiwanese sources collected in the same run: a foreign wire item whose event is already covered by a Taiwanese source SHALL be judged as not reported, recorded via seen.py record with decision skipped, and the Taiwanese version SHALL be the one used. Only foreign wire items covering events absent from the Taiwanese sources SHALL feed report content.

#### Scenario: Event already covered by Taiwanese media

- **WHEN** a CNBC item covers a Fed statement that a cnyes article in the same collection window already reports
- **THEN** the CNBC item is recorded as skipped and overview_md draws on the cnyes version

#### Scenario: Event only in foreign wires

- **WHEN** a WSJ Markets item reports an export-control detail that no Taiwanese source in the window covers
- **THEN** the item informs overview_md or industry event background, and no calls or stock_news entry cites it

## MODIFIED Requirements

### Requirement: Pinned source endpoints

The skill SHALL pin every news source to a named list endpoint verified to be reachable and to yield full article text, and SHALL NOT direct the agent to discover sources from portal homepages. The pinned set SHALL be: cnyes category list API with the six categories tw_stock, headline, tech, tw_macro, cnyeshouse, and wd_stock, whose list response already carries full article content and stock tags so no detail-page fetch is made, with wd_stock feeding overview_md only; TWSE OpenAPI datasets t187ap04_L (material information) and t187ap05_L (monthly revenue) pulled whole and filtered by the collection window; CTEE rss_web livenews RSS for the six categories policy, stock, finance, industry, house, and tech, with article bodies fetched over plain HTTP and no headless browser; CNA finance RSS; TechNews site feed; LTN business RSS; WSJ Markets RSS at feeds.a.dj.com/rss/RSSMarketsMain.xml, used from feed item titles and descriptions only because WSJ article pages sit behind a paywall; and CNBC top stories RSS at www.cnbc.com/id/100003114/device/rss/rss.html, with CNBC article bodies fetched over plain HTTP; both foreign wire feeds are subject to the foreign wire scope requirement. The skill SHALL NOT list UDN as a source.

#### Scenario: Collection uses pinned endpoints

- **WHEN** the agent collects news for the day
- **THEN** every list fetch targets one of the pinned endpoints and cnyes article content is taken from the list response without a detail-page fetch

#### Scenario: No headless browser required

- **WHEN** the agent collects CTEE articles
- **THEN** list discovery uses the CTEE RSS endpoints and article bodies are fetched over plain HTTP

### Requirement: Batched collection with per-batch counts

The skill SHALL split collection into five named batches, each naming the endpoints it covers and each reporting its own counts. Batch one SHALL cover the six cnyes category endpoints and report, per category, how many items were returned and how many of those were unread. Batch two SHALL cover the two TWSE datasets and report, per dataset, how many records fall inside the collection window. Batch three SHALL cover the six CTEE category feeds and report unread counts per category. Batch four SHALL cover the CNA, TechNews, and LTN feeds and report unread counts per source. Batch five SHALL cover the WSJ Markets and CNBC top stories feeds and report unread counts per source.

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

- **WHEN** both feeds in batch five remain unreachable after one retry
- **THEN** the agent publishes from the remaining batches and lists the two feeds as missing sources in its run report
