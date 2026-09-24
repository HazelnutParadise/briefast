## MODIFIED Requirements

### Requirement: Daily market outlook payload

The report API SHALL accept an optional `market_outlook` object with `direction`, `summary_md`, optional `trajectory_md`, and optional `trajectory_chart`. When present, `direction` MUST be one of `up`, `down`, `range`, or `uncertain`, and `summary_md` MUST contain non-whitespace text. A non-null provided `trajectory_md` MUST contain non-whitespace text; null SHALL be treated as absent. A non-null `trajectory_chart` MUST have `open`, `midday`, and `close` values drawn only from `above`, `near`, and `below`, and MUST have a nonblank `trajectory_md`. Its `close` MUST be `above` for `up`, `below` for `down`, and `near` for `range`; `uncertain` MUST NOT carry a chart. Invalid outlooks SHALL produce HTTP 400 with an error naming the offending `market_outlook` field and SHALL NOT persist the report. Successful reads SHALL return the same outlook, including its trajectory and chart when present. Reports without the object or without either optional trajectory field SHALL remain valid and readable.

#### Scenario: Valid outlook round trip

- **WHEN** a valid report contains `market_outlook` with `direction: "up"`, nonblank `summary_md`, nonblank `trajectory_md`, and `trajectory_chart` with `open: "above"`, `midday: "near"`, and `close: "above"`
- **THEN** POST returns 200 and a dated GET returns the same outlook including the trajectory and chart

#### Scenario: Invalid outlook rejected

- **WHEN** a report contains an unsupported direction, a blank summary, and a whitespace-only trajectory
- **THEN** POST returns 400 with errors naming all three offending `market_outlook` fields and no report is stored

#### Scenario: Invalid chart rejected

- **WHEN** a report contains a chart with a missing phase, an unsupported phase value, a close level contrary to the direction, or no nonblank `trajectory_md`
- **THEN** POST returns 400 naming `market_outlook.trajectory_chart` and no report is stored

#### Scenario: Uncertain outlook has no chart

- **WHEN** a report with `direction: "uncertain"` contains `trajectory_chart`
- **THEN** POST returns 400 naming `market_outlook.trajectory_chart` and no report is stored

#### Scenario: Legacy outlook accepted

- **WHEN** a valid report contains an outlook with direction and summary but omits `trajectory_md` and `trajectory_chart`
- **THEN** POST returns 200 and the report remains readable

#### Scenario: Legacy report accepted

- **WHEN** a valid report omits `market_outlook`
- **THEN** POST returns 200 and the report remains readable
