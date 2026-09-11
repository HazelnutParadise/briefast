## MODIFIED Requirements

### Requirement: Pinned source endpoints

The skill SHALL pin every news source to a named list endpoint verified to be reachable and to yield full article text, and SHALL NOT direct the agent to discover sources from portal homepages. The pinned set SHALL be: cnyes category list API with the six categories tw_stock, headline, tech, tw_macro, cnyeshouse, and wd_stock, whose list response already carries full article content and stock tags so no detail-page fetch is made, with wd_stock feeding overview_md only; TWSE OpenAPI datasets t187ap04_L (material information) and t187ap05_L (monthly revenue) pulled whole and filtered by the collection window; CTEE rss_web livenews RSS for the six categories policy, stock, finance, industry, house, and tech, with article bodies fetched over plain HTTP and no headless browser; CNA finance RSS; TechNews site feed; LTN business RSS; and CNBC top stories RSS at www.cnbc.com/id/100003114/device/rss/rss.html, with CNBC article bodies fetched over plain HTTP and subject to the foreign wire scope requirement. The skill SHALL NOT list WSJ Markets as a source, because its article pages sit behind a paywall and feed summaries alone cannot support judgement, and SHALL NOT list UDN as a source.

Because CNBC sits behind Akamai and rejects tool-style User-Agent values (curl, python-requests, Claude-User and similar) with 403 Access Denied on both the feed and article pages, the skill SHALL require every CNBC request — the RSS list fetch and each article-page fetch — to carry a desktop browser User-Agent header, SHALL print one working User-Agent string and a copyable curl invocation in the batch five instructions, and SHALL direct the agent, on any CNBC 403, to first confirm the header was sent and resend with it before treating the response as a source failure. Only a 403 or connection failure that persists with the browser User-Agent SHALL count as a CNBC failure under the source collection completeness gate. The CNBC feed description carries only a one-sentence summary, so the skill SHALL NOT offer feed summaries as a substitute for fetching the article page.

#### Scenario: Collection uses pinned endpoints

- **WHEN** the agent collects news for the day
- **THEN** every list fetch targets one of the pinned endpoints and cnyes article content is taken from the list response without a detail-page fetch

#### Scenario: No headless browser required

- **WHEN** the agent collects CTEE articles
- **THEN** list discovery uses the CTEE RSS endpoints and article bodies are fetched over plain HTTP

#### Scenario: CNBC requests carry a browser User-Agent

- **WHEN** the agent fetches the CNBC top stories RSS or any CNBC article page
- **THEN** the request carries the desktop browser User-Agent string given in the skill, the feed returns 200 with its items, and the article page returns 200 with parsable body text

##### Example: Same endpoint, different User-Agent

| User-Agent sent | RSS response | Article page response |
| --------------- | ------------ | --------------------- |
| curl default (curl/8.x) | 403 Access Denied, server AkamaiGHost | 403 |
| Claude-User/1.0 | 403 | 403 |
| Desktop Chrome string from the skill | 200, 30 items | 200, body text extractable |

#### Scenario: CNBC 403 is checked against the header before being recorded as failure

- **WHEN** a CNBC request returns 403
- **THEN** the agent verifies the browser User-Agent header was sent, resends with it if it was missing, and records CNBC as a failed source only when the browser-User-Agent request also fails after the one retry allowed by the completeness gate
