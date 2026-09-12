## ADDED Requirements

### Requirement: Conference brief ingestion endpoint

The API SHALL expose POST /api/conference accepting one JSON object describing a company-held investor conference and its pre-conference brief. The endpoint SHALL require the same Bearer key authentication as POST /api/report, re-checking key status on every request; a missing, unknown, or revoked key SHALL receive 401, SHALL NOT be persisted, and SHALL be recorded in the update log with a rejected-auth action. The request body SHALL be limited to 2 MiB and SHALL contain exactly one JSON object with no unknown fields.

The object SHALL carry: symbol (Taiwan stock code), name, market (one of twse, tpex), held_on (YYYY-MM-DD), held_at (HH:MM or empty), venue, announced_on (YYYY-MM-DD), prediction (one of bull, bear, none), headline, summary_md, watch_md, an optional fundamentals object, an optional chips object identical in shape and validation to the stock_news chips object, sources (array of title and url), and generated_at (RFC 3339). Validation SHALL enforce: symbol non-empty; name non-empty; market is twse or tpex; held_on and announced_on are valid dates; held_at is empty or HH:MM; prediction is one of the three values; headline, summary_md, and watch_md non-empty; every source has a non-empty url; chips.date valid when chips is present; generated_at parses as RFC 3339. Every violation SHALL be listed in a 400 response with ok=false and an errors array, nothing SHALL be persisted, and the attempt SHALL be recorded in the update log with a rejected-schema action.

The fundamentals object, when present, SHALL carry an optional revenue object and an optional quarter object. The revenue object SHALL carry month (YYYY-MM), current, prev_month, last_year_month, ytd, last_year_ytd (all integers in thousands of NTD). The quarter object SHALL carry label (free text naming the period, for example 2026Q2 累計), revenue, operating_income, net_income (integers in thousands of NTD), and eps (decimal string). Absent fields SHALL be omitted, never zero-filled.

On a valid authenticated request the API SHALL upsert the row keyed by symbol and held_on (last write wins, settlement fields preserved), insert one update-log row with the key name snapshotted, the held_on date in the report_date column, and a conference_ok action, trigger the live update notification, and respond 200 with ok=true, symbol, and held_on.

#### Scenario: Valid brief persisted and logged

- **WHEN** a valid brief for symbol 2330 held on 2026-10-16 is posted with an active key
- **THEN** the row for (2330, 2026-10-16) is stored, one update-log row with action conference_ok and report_date 2026-10-16 is inserted, connected viewers are notified, and the response is 200 with symbol and held_on

#### Scenario: Same conference re-posted keeps settlement

- **WHEN** a brief for (2330, 2026-10-16) that was already settled is posted again with a new summary
- **THEN** the stored brief content is replaced, the stored pre_close, post_close, and outcome fields are unchanged, and a second conference_ok log row is inserted

#### Scenario: Invalid brief rejected atomically

- **WHEN** a brief has prediction "up" and an empty headline
- **THEN** the API responds 400 listing both violations and no row is written

##### Example: validation matrix

| Field | Value | Result |
|-------|-------|--------|
| prediction | bull | accepted |
| prediction | up | 400, prediction must be bull, bear, or none |
| market | otc | 400, market must be twse or tpex |
| held_at | 14:00 | accepted |
| held_at | 2pm | 400, held_at must be HH:MM or empty |
| fundamentals | absent | accepted |
| chips.date | 2026-13-01 | 400, chips.date must be a valid date |

#### Scenario: Unauthenticated post rejected and logged

- **WHEN** a request arrives without a valid bearer key
- **THEN** the API responds 401, nothing is persisted, and a rejected-auth entry with action conference_rejected_auth is recorded


### Requirement: Conference settlement endpoint

The API SHALL expose POST /api/conference/settle accepting a JSON object with symbol, held_on, pre_close_date, pre_close, post_close_date, and post_close, where the two dates are YYYY-MM-DD and the two closes are positive decimals. The endpoint SHALL require the same Bearer key authentication and logging behaviour as ingestion, using settle_rejected_auth and settle_rejected_schema actions. A settlement for a (symbol, held_on) pair with no stored brief SHALL receive 404. Validation SHALL enforce pre_close_date strictly before held_on and post_close_date strictly after held_on.

On success the API SHALL store the four values and a settled_at timestamp, and SHALL compute outcome from the stored prediction: hit when prediction is bull and post_close is greater than pre_close, hit when prediction is bear and post_close is less than pre_close, miss for bull or bear otherwise (an unchanged close counts as miss), and empty outcome when prediction is none. The API SHALL insert one update-log row with action settle_ok and respond 200 with ok=true, symbol, held_on, and outcome. Re-posting a settlement SHALL overwrite the previous settlement and recompute outcome.

#### Scenario: Bullish prediction settled as hit

- **WHEN** the stored prediction is bull, pre_close is 1,080 on 2026-10-15, and post_close is 1,095 on 2026-10-17
- **THEN** outcome is hit, settled_at is set, and the response carries outcome hit

##### Example: outcome matrix

| Prediction | pre_close | post_close | Outcome |
|------------|-----------|------------|---------|
| bull | 1080 | 1095 | hit |
| bull | 1080 | 1080 | miss |
| bull | 1080 | 1060 | miss |
| bear | 500 | 480 | hit |
| bear | 500 | 510 | miss |
| none | 500 | 520 | (empty) |

#### Scenario: Settlement dates outside the conference rejected

- **WHEN** a settlement carries pre_close_date equal to held_on
- **THEN** the API responds 400 naming pre_close_date and nothing is stored

#### Scenario: Settlement for unknown conference

- **WHEN** a settlement names a (symbol, held_on) pair with no stored brief
- **THEN** the API responds 404 with ok=false and no settlement is stored


### Requirement: Conference listing endpoint

The API SHALL expose GET /api/conferences requiring the same Bearer key authentication and logging as the report read endpoint (list_rejected_auth action on rejection). The response SHALL be a JSON object with two arrays: upcoming, holding every stored conference whose held_on is on or after the current date in Asia/Taipei ordered by held_on ascending then symbol; and pending_settlement, holding every stored conference whose held_on is before the current Taipei date and whose settled_at is empty, ordered by held_on ascending. Each element SHALL carry symbol, name, market, held_on, held_at, prediction, and, for pending_settlement, nothing more; upcoming elements SHALL additionally carry venue and announced_on so the daily workflow can detect a rescheduled or relocated conference without decoding briefs. The full brief text SHALL NOT be included in the listing. Successful listings SHALL NOT be written to the update log.

#### Scenario: Listing separates upcoming from pending

- **WHEN** today in Taipei is 2026-09-14 and stored conferences are held on 2026-09-11 (unsettled), 2026-09-12 (settled), 2026-09-14, and 2026-09-16
- **THEN** upcoming contains 2026-09-14 then 2026-09-16, and pending_settlement contains only 2026-09-11

#### Scenario: Empty store

- **WHEN** no conference is stored
- **THEN** the response is 200 with two empty arrays


### Requirement: Conference storage

Conference rows SHALL live in a dedicated conferences table created by a new schema migration version 2, leaving migration 1 untouched. The primary key SHALL be (symbol, held_on). The table SHALL store the brief payload as the ingested JSON blob plus indexed columns for held_on, prediction, and settled_at so that the home page can list upcoming conferences and count settled outcomes without decoding payloads. Rows SHALL carry created_at and updated_at fixed-width timestamps consistent with the reports table. Opening a database that already carries migration 1 SHALL apply migration 2 in one transaction and record version 2 in schema_migrations.

#### Scenario: Migration applied on existing database

- **WHEN** the application opens a database at schema version 1
- **THEN** the conferences table exists afterwards and schema_migrations records version 2

#### Scenario: Fresh database

- **WHEN** the application opens a new database
- **THEN** both migrations apply and schema_migrations records version 2
