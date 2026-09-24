## ADDED Requirements

### Requirement: Daily market outlook payload

The report API SHALL accept an optional `market_outlook` object with `direction` and `summary_md`. When present, `direction` MUST be one of `up`, `down`, `range`, or `uncertain`, and `summary_md` MUST contain non-whitespace text. Invalid outlooks SHALL produce HTTP 400 with an error naming `market_outlook` and SHALL NOT persist the report. Successful reads SHALL return the same outlook. Reports without the object SHALL remain valid and readable.

#### Scenario: Valid outlook round trip

- **WHEN** a valid report contains `market_outlook` with `direction` set to `up` and nonblank `summary_md`
- **THEN** POST returns 200 and a dated GET returns the same outlook

#### Scenario: Invalid outlook rejected

- **WHEN** a report contains an unsupported direction and a blank summary
- **THEN** POST returns 400 with both errors naming `market_outlook` and no report is stored

#### Scenario: Legacy report accepted

- **WHEN** a valid report omits `market_outlook`
- **THEN** POST returns 200 and the report remains readable
