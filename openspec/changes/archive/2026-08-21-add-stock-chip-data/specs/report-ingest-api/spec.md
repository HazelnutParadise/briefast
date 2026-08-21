## MODIFIED Requirements

### Requirement: Report schema validation

The API SHALL validate the request body against the report schema before persisting. On any violation the API SHALL respond 400 with a JSON body containing ok=false and an errors array describing every violation, and SHALL NOT persist any part of the report. The rejected attempt SHALL be recorded in the update log with a rejected-schema action. Validation SHALL enforce at minimum: date matches YYYY-MM-DD; headline, overview_md, and watch_md are non-empty; every call entry value is one of short_bull, short_bear, long_bull, long_bear, none; every source has a url; every industries entry and every stock_news entry carries a non-empty watch_md; every industries entry carries at least one event; every event carries a non-empty headline and a non-empty summary_md; every stock_news entry carries a non-empty headline. Each violation message SHALL identify the offending entry by index so an agent can locate it without guessing.

A stock_news entry SHALL accept an optional chips object carrying the previous trading day's chip data: date (chip data date), foreign_net, trust_net, dealer_net, and total_net (net buy/sell in shares, integers, negative for net sell), and optional margin_change and short_change (margin and short balance change in trading units, integers, negative for decrease). An absent or null chips object SHALL be accepted. When chips is present, validation SHALL enforce that its date matches YYYY-MM-DD, with the violation message identifying the stock_news index. Reports stored without chips SHALL remain readable and valid.

#### Scenario: Invalid payload rejected atomically

- **WHEN** a payload has an empty headline and a malformed date
- **THEN** the API responds 400 with both violations listed in errors and no report row is persisted

#### Scenario: Missing watch points rejected

- **WHEN** a payload has a stock_news entry without watch_md and an industries entry whose watch_md is only whitespace
- **THEN** the API responds 400 with both violations listed in errors and no report row is persisted

#### Scenario: Missing headlines rejected

- **WHEN** a payload has an industries entry whose second event has a blank headline and a stock_news entry with no headline
- **THEN** the API responds 400 with both violations listed in errors, each naming the offending index, and no report row is persisted

#### Scenario: Industry section without events rejected

- **WHEN** a payload has an industries entry whose events array is empty
- **THEN** the API responds 400 with that violation listed in errors and no report row is persisted

#### Scenario: Chips accepted and round-tripped

- **WHEN** a valid payload has a stock_news entry whose chips object carries a valid date and integer net values
- **THEN** the API responds 200 and a subsequent authenticated read of that date returns the entry with the identical chips object

#### Scenario: Entry without chips accepted

- **WHEN** a valid payload has stock_news entries with no chips object
- **THEN** the API responds 200 and the report is persisted

#### Scenario: Chips with malformed date rejected

- **WHEN** a payload's third stock_news entry carries a chips object whose date is not a valid YYYY-MM-DD value
- **THEN** the API responds 400 with a violation naming stock_news index 2 and no report row is persisted
