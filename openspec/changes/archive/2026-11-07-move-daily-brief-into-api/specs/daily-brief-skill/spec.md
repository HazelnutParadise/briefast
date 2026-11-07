## MODIFIED Requirements

### Requirement: Daily workflow skill exists in the repo

The daily pre-market workflow SHALL be served by GET /api/guide as ordered steps: collect industry and stock market news, judge stock calls from the news, compose the report JSON, and POST it to the site API. The repo SHALL NOT ship a separate daily-brief skill package. The README SHALL provide the one-line Cowork task instruction that tells the agent to read the workspace .env and then follow GET /api/guide.

#### Scenario: Skill discoverable

- **WHEN** an agent configured only with the site URL and an API key calls GET /api/guide
- **THEN** it receives the daily-brief workflow with the four workflow steps documented

#### Scenario: No skill package left

- **WHEN** the repo is checked after this change
- **THEN** skills/daily-brief/ does not exist and no file refers to its SKILL.md or scripts

---
### Requirement: Endpoint configuration via environment

The guide SHALL direct the agent to read the site URL from the BRIEFAST_URL key and the API key from the BRIEFAST_API_KEY key of a .env file at the execution workspace root, in KEY=VALUE format, and SHALL NOT contain a hardcoded URL or key. The README SHALL show a workspace .env template listing the two keys with placeholder values and no real secrets. The guide and the README task instruction SHALL tell the agent to stop and report when the .env file is missing or either key is missing or blank.

#### Scenario: Missing credentials

- **WHEN** the workspace root has no .env file, or BRIEFAST_API_KEY is missing or blank in it, when the workflow runs
- **THEN** the workflow stops before collecting and reports the missing configuration

#### Scenario: Template present

- **WHEN** a user prepares a new execution workspace
- **THEN** the README shows the two required keys with placeholder values to copy into the workspace .env

---
### Requirement: Foreign wire scope and cross-language deduplication

The guide SHALL restrict foreign wire content (CNBC top stories) to overview_md and to background context inside industry event summaries, and SHALL NOT allow foreign wire items to create or support calls or stock_news entries. Because POST /api/seen/similar compares article bodies within one language and cannot match a Chinese article against an English article covering the same event, the guide SHALL direct the agent to deduplicate foreign wire items by event against the Taiwanese sources collected in the same run: a foreign wire item whose event is already covered by a Taiwanese source SHALL be judged as not reported, recorded via POST /api/seen/record with decision skipped, and the Taiwanese version SHALL be the one used. Only foreign wire items covering events absent from the Taiwanese sources SHALL feed report content.

#### Scenario: Event already covered by Taiwanese media

- **WHEN** a CNBC item covers a Fed statement that a cnyes article in the same collection window already reports
- **THEN** the CNBC item is recorded as skipped and overview_md draws on the cnyes version

#### Scenario: Event only in foreign wires

- **WHEN** a CNBC item reports an export-control detail that no Taiwanese source in the window covers
- **THEN** the item informs overview_md or industry event background, and no calls or stock_news entry cites it

## REMOVED Requirements

### Requirement: Verified TX night extraction

**Reason**: The local `tx_night.py` parser is removed together with the skill package. The server now performs the same extraction.

**Migration**: Upload the two TX CSV files to POST /api/tx-night, specified in the `tx-night-api` capability with the same selection rules and scenarios.
