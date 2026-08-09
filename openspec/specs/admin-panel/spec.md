# admin-panel Specification

## Purpose

TBD - created by archiving change 'daily-report-site'. Update Purpose after archive.

## Requirements

### Requirement: Admin login gate

The /admin page SHALL show only a login form until the operator authenticates with the admin password, which SHALL be read from configuration (syralit.toml secrets or the BRIEFAST_ADMIN_PASSWORD environment variable). When no admin password is configured, the admin page SHALL refuse to serve admin functions and SHALL display a configuration error instead.

#### Scenario: Unauthenticated visitor

- **WHEN** a visitor opens /admin without logging in
- **THEN** only the login form is shown and no key or log data is rendered

#### Scenario: Missing password configuration

- **WHEN** no admin password is configured and /admin is opened
- **THEN** a configuration error is shown and login is not possible


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
### Requirement: API key creation

An authenticated admin SHALL be able to create an API key by entering a name. The system SHALL generate a random token, persist the key with its name and creation time to the key store, and show it in the key list immediately.

#### Scenario: Create a named key

- **WHEN** the admin creates a key named cowork-daily
- **THEN** a random token is generated and the key appears in the list with its name and creation time


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
### Requirement: API keys remain fully viewable

The key list SHALL display the complete plaintext token of every non-revoked key on every visit, not only at creation time. Tokens SHALL therefore be stored in plaintext in the key store.

#### Scenario: Token visible on a later visit

- **WHEN** the admin logs in again days after creating a key
- **THEN** the full token of that key is still displayed in the list


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
### Requirement: API key revocation

An authenticated admin SHALL be able to revoke a key. Revocation SHALL persist a revocation timestamp and take effect immediately for API authentication. Revoked keys SHALL be visually distinguished from active keys in the list.

#### Scenario: Revoke a key

- **WHEN** the admin revokes a key
- **THEN** the key is marked revoked in the list and subsequent API requests using it receive 401


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
### Requirement: Report update log view

The admin panel SHALL display the report update log from newest to oldest, each entry showing the time, the key name used, the report date, and the action.

#### Scenario: Log entries listed

- **WHEN** two reports were ingested and the admin opens the update log
- **THEN** both entries are listed newest first with time, key name, and report date

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