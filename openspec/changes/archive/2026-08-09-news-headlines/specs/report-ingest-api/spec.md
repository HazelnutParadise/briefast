## MODIFIED Requirements

### Requirement: Report schema validation

The API SHALL validate the request body against the report schema before persisting. On any violation the API SHALL respond 400 with a JSON body containing ok=false and an errors array describing every violation, and SHALL NOT persist any part of the report. The rejected attempt SHALL be recorded in the update log with a rejected-schema action. Validation SHALL enforce at minimum: date matches YYYY-MM-DD; headline, overview_md, and watch_md are non-empty; every call entry value is one of short_bull, short_bear, long_bull, long_bear, none; every source has a url; every industries entry and every stock_news entry carries a non-empty watch_md; every industries entry carries at least one event; every event carries a non-empty headline and a non-empty summary_md; every stock_news entry carries a non-empty headline. Each violation message SHALL identify the offending entry by index so an agent can locate it without guessing.

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
