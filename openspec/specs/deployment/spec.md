# deployment Specification

## Purpose

TBD - created by archiving change 'daily-report-site'. Update Purpose after archive.

## Requirements

### Requirement: Docker compose deployment

The repo SHALL contain a Dockerfile and docker-compose.yml such that docker compose up -d --build builds the site and serves it on the configured port. The data directory SHALL be mounted as a volume so reports, keys, and the update log survive container rebuilds and restarts.

#### Scenario: Data survives restart

- **WHEN** a report and a key exist and the container is rebuilt and restarted
- **THEN** the report history and the key are still present and served


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
### Requirement: Configuration injection

Runtime secrets (admin password, optional overrides) SHALL be provided via environment variables or a syralit.toml file mounted at deploy time. The repo SHALL provide syralit.toml.example documenting every setting, and the real syralit.toml SHALL be excluded from version control via .gitignore.

#### Scenario: Example config provided

- **WHEN** an operator prepares a new deployment
- **THEN** copying syralit.toml.example and filling in the admin password yields a working configuration

#### Scenario: Secrets never committed

- **WHEN** git status is inspected after creating a local syralit.toml
- **THEN** syralit.toml is ignored and only syralit.toml.example is tracked

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