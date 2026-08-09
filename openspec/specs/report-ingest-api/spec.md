# report-ingest-api Specification

## Purpose

TBD - created by archiving change 'daily-report-site'. Update Purpose after archive.

## Requirements

### Requirement: Bearer key authentication

POST /api/report SHALL require an Authorization header of the form Bearer <key>. Requests whose key is missing, unknown, or revoked SHALL be rejected with HTTP 401 and SHALL NOT be persisted. Key validity SHALL be checked against the current key store on every request so that revocation takes effect immediately. Each rejected request SHALL be recorded in the update log with a rejected-auth action.

#### Scenario: Missing or unknown key

- **WHEN** a request arrives without a valid bearer key
- **THEN** the API responds 401, no report is persisted, and a rejected-auth entry is recorded in the update log

#### Scenario: Revoked key rejected

- **WHEN** a key is revoked in the admin panel and a request signed with it arrives afterwards
- **THEN** the API responds 401


<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
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


<!-- @trace
source: news-headlines
updated: 2026-08-09
code:
  - internal/api/report.go
  - internal/site/site.go
  - skills/daily-brief/SKILL.md
  - internal/site/styles.go
  - internal/report/schema.go
  - internal/api/read.go
  - main.go
tests:
  - internal/api/report_test.go
  - internal/report/schema_test.go
  - internal/api/read_test.go
  - internal/site/site_test.go
  - main_test.go
-->

---
### Requirement: Calls and stock news consistency

The API SHALL reject a report in which any symbol listed in the four calls lists has no matching stock_news entry. Stock_news entries with call none and no calls-list membership SHALL be accepted.

#### Scenario: Call without detail rejected

- **WHEN** calls.short_bull contains symbol 2330 but stock_news has no entry with symbol 2330
- **THEN** the API responds 400 naming the missing symbol

##### Example: consistency matrix

| calls lists | stock_news symbols | Result |
| ----------- | ------------------ | ------ |
| short_bull: 2330 | 2330 (call short_bull) | accepted |
| short_bull: 2330 | 2317 (call none) | 400, missing 2330 |
| (all empty) | 2317 (call none) | accepted |


<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
### Requirement: Successful ingestion

On a valid authenticated request the API SHALL upsert the report into the SQLite database keyed by its date, insert one update-log row recording the time, the key name (snapshotted at write time), the report date, and the action, trigger a live update notification, and respond 200 with ok=true and the report date. A subsequent valid request for the same date SHALL overwrite the stored report (last write wins).

#### Scenario: Valid report persisted

- **WHEN** a valid report for 2026-08-07 is posted with an active key
- **THEN** the stored report for 2026-08-07 is retrievable from the database, one update-log row is inserted, and the response is 200 ok

#### Scenario: Same-date overwrite

- **WHEN** a second valid report for 2026-08-07 is posted
- **THEN** the stored report contains the second payload and a second update-log row is inserted


<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
### Requirement: Live update push

When a report is successfully ingested, all currently connected browser sessions SHALL be refreshed to show the new report without a manual reload.

#### Scenario: Connected viewer sees new report

- **WHEN** a visitor has the homepage open and a new report is ingested
- **THEN** the visitor's page updates to the new report without pressing reload

<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
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

<!-- @trace
source: daily-brief-analyst-guidance
updated: 2026-08-09
code:
  - internal/api/read.go
  - internal/report/schema.go
  - skills/daily-brief/SKILL.md
  - internal/api/report.go
  - main.go
  - internal/site/styles.go
  - internal/site/site.go
tests:
  - internal/api/report_test.go
  - main_test.go
  - internal/site/site_test.go
  - internal/api/read_test.go
  - internal/report/schema_test.go
-->