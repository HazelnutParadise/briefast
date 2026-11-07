## MODIFIED Requirements

### Requirement: Authenticated report read endpoint

The API SHALL expose GET /api/report/{date} returning `{"ok":true,"guide":…,"report":…}` on success, where `report` is the stored report for that date in the same JSON structure that was ingested. The endpoint SHALL require the same Bearer key authentication as ingestion and SHALL re-check key status on every request so revocation takes effect immediately. Requests with a missing, invalid, or revoked key SHALL receive 401 and SHALL be recorded in the update log with a read-rejected-auth action. A date that is not a valid YYYY-MM-DD value SHALL receive 400. A valid date with no stored report SHALL receive 404. All non-200 responses SHALL carry ok=false and an errors array. Successful reads SHALL NOT be written to the update log.

#### Scenario: Authorized read returns the report

- **WHEN** a request with a valid Bearer key asks for a date that has a stored report
- **THEN** the API responds 200 with ok=true, the guide pointer, and that report's full JSON under `report` including every industries and stock_news entry

#### Scenario: Read report can be resubmitted

- **WHEN** an agent takes the `report` object from a read response and POSTs it unchanged to /api/report
- **THEN** the ingest accepts it without an unknown-field error

#### Scenario: Revoked key rejected and logged

- **WHEN** a request uses a key that was revoked after it was issued
- **THEN** the API responds 401 and records a read-rejected-auth entry in the update log

#### Scenario: Unknown date

- **WHEN** a request with a valid Bearer key asks for a well-formed date that has no stored report
- **THEN** the API responds 404 with ok=false and no update log entry
