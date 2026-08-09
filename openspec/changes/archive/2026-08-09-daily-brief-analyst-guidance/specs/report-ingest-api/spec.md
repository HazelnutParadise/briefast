## ADDED Requirements

### Requirement: Authenticated report read endpoint

The API SHALL expose GET /api/report/{date} returning the stored report for that date as the same JSON structure that was ingested. The endpoint SHALL require the same Bearer key authentication as ingestion and SHALL re-check key status on every request so revocation takes effect immediately. Requests with a missing, invalid, or revoked key SHALL receive 401 and SHALL be recorded in the update log with a read-rejected-auth action. A date that is not a valid YYYY-MM-DD value SHALL receive 400. A valid date with no stored report SHALL receive 404. All non-200 responses SHALL carry ok=false and an errors array. Successful reads SHALL NOT be written to the update log.

#### Scenario: Authorized read returns the report

- **WHEN** a request with a valid Bearer key asks for a date that has a stored report
- **THEN** the API responds 200 with that report's full JSON including every industries and stock_news entry

#### Scenario: Revoked key rejected and logged

- **WHEN** a request uses a key that was revoked after it was issued
- **THEN** the API responds 401 and records a read-rejected-auth entry in the update log

#### Scenario: Unknown date

- **WHEN** a request with a valid Bearer key asks for a well-formed date that has no stored report
- **THEN** the API responds 404 with ok=false and no update log entry

## MODIFIED Requirements

### Requirement: Report schema validation

The API SHALL validate the request body against the report schema before persisting. On any violation the API SHALL respond 400 with a JSON body containing ok=false and an errors array describing every violation, and SHALL NOT persist any part of the report. The rejected attempt SHALL be recorded in the update log with a rejected-schema action. Validation SHALL enforce at minimum: date matches YYYY-MM-DD; headline, overview_md, and watch_md are non-empty; every call entry value is one of short_bull, short_bear, long_bull, long_bear, none; every source has a url; every industries entry and every stock_news entry carries a non-empty watch_md.

#### Scenario: Invalid payload rejected atomically

- **WHEN** a payload has an empty headline and a malformed date
- **THEN** the API responds 400 with both violations listed in errors and no report row is persisted

#### Scenario: Missing watch points rejected

- **WHEN** a payload has a stock_news entry without watch_md and an industries entry whose watch_md is only whitespace
- **THEN** the API responds 400 with both violations listed in errors and no report row is persisted
