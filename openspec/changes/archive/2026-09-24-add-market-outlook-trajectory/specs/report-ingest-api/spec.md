## MODIFIED Requirements

### Requirement: Daily market outlook payload

The report API SHALL accept an optional `market_outlook` object with `direction`, `summary_md`, and optional `trajectory_md`. When present, `direction` MUST be one of `up`, `down`, `range`, or `uncertain`, and `summary_md` MUST contain non-whitespace text. A non-null provided `trajectory_md` MUST contain non-whitespace text; null SHALL be treated as absent for compatibility. Invalid outlooks SHALL produce HTTP 400 with an error naming the offending `market_outlook` field and SHALL NOT persist the report. Successful reads SHALL return the same outlook, including its trajectory when present. Reports without the object or without the optional trajectory SHALL remain valid and readable.

#### Scenario: Valid outlook round trip

- **WHEN** a valid report contains `market_outlook` with `direction` set to `up`, nonblank `summary_md`, and nonblank `trajectory_md`
- **THEN** POST returns 200 and a dated GET returns the same outlook including the trajectory

#### Scenario: Invalid outlook rejected

- **WHEN** a report contains an unsupported direction, a blank summary, and a whitespace-only trajectory
- **THEN** POST returns 400 with errors naming all three offending `market_outlook` fields and no report is stored

#### Scenario: Legacy outlook accepted

- **WHEN** a valid report contains an outlook with direction and summary but omits `trajectory_md`
- **THEN** POST returns 200 and the report remains readable

#### Scenario: Legacy report accepted

- **WHEN** a valid report omits `market_outlook`
- **THEN** POST returns 200 and the report remains readable
