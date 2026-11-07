## ADDED Requirements

### Requirement: Verified TX night extraction endpoint

The API SHALL expose POST /api/tx-night accepting multipart form fields `report_date` and `previous_date` (YYYY-MM-DD) and files `report_csv` and `previous_csv` holding TAIFEX TX daily CSV downloads in CP950 encoding. It SHALL require Bearer key authentication, recording 401 and 400 responses in the update log. It SHALL select the unique highest-volume nonzero single-month TX `盤後` row attributed to the report date and compare its closing price with the unique `一般` row for the same contract on the preceding trading date. It SHALL reject stale dates, missing required columns, malformed or ambiguous rows, missing matching contracts, and nonpositive or nonfinite prices with a 400 response whose errors describe the failure and which carries no partial result. On success it SHALL respond 200 with ok=true, `report_date`, `previous_regular_date`, `contract`, both closing prices, `change_points`, `change_pct`, `night_volume`, and `direction` set to `up`, `down`, or `flat`.

#### Scenario: Completed night session

- **WHEN** the report date is 2026-09-24, the preceding trading date is 2026-09-23, the highest-volume report-date TX night contract is 202610 with close 47909 and volume 28947, and the preceding regular session for 202610 closed at 48336
- **THEN** the API responds 200 with contract 202610, `change_points: -427`, and `direction: "down"`, using no report-date regular-session price

#### Scenario: Weekend and contract rollover

- **WHEN** the report date follows a weekend or contract rollover and the supplied preceding trading date is several calendar days earlier
- **THEN** the API pairs the report-date night close only with the matching contract on the supplied preceding trading date, without subtracting one calendar day or substituting another contract

#### Scenario: Unusable source rows

- **WHEN** the report-date night row is absent, two eligible contracts tie for highest volume, the preceding same-contract regular row is missing, or a required date, price, volume, or CSV header is invalid
- **THEN** the API responds 400 with an errors array describing the failure and no result fields

#### Scenario: Flat overnight move

- **WHEN** the two valid same-contract closing prices are equal
- **THEN** the API responds with zero point and percentage change and `direction: "flat"`, without forcing a bullish or bearish signal
