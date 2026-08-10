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

Non-secret runtime configuration SHALL live in a single syralit.toml file tracked in version control, and the repo SHALL NOT carry a separate example copy of it. The tracked syralit.toml SHALL NOT contain a secrets section. Secrets and deployment variables SHALL live in a .env file at the repo root that is excluded from version control, and the repo SHALL track a .env.example listing every key with its purpose and no real values. The admin password SHALL be supplied through the BRIEFAST_ADMIN_PASSWORD key of that .env, resolved as an environment variable at runtime. Docker Compose SHALL take BRIEFAST_CONFIG, BRIEFAST_PORT, and BRIEFAST_ADMIN_PASSWORD from the same .env through its native variable substitution, requiring no extra flags. Deployments needing environment-specific settings SHALL point BRIEFAST_CONFIG at their own file rather than editing the tracked one.

#### Scenario: Clone and run without setup steps

- **WHEN** an operator clones the repo, copies .env.example to .env, fills in the admin password, and starts the stack with Compose
- **THEN** the tracked syralit.toml is mounted and the password reaches the container with no additional flags or copying of config templates

#### Scenario: Secrets never committed

- **WHEN** git status is inspected after configuring an admin password
- **THEN** the password appears in no tracked file, .env is ignored, and the tracked syralit.toml carries no secrets section at all

#### Scenario: Environment-specific override

- **WHEN** a deployment sets BRIEFAST_CONFIG in its .env to its own configuration file
- **THEN** Compose mounts that file instead of the tracked syralit.toml, with the same parsing behavior as before

<!-- @trace
source: track-syralit-config
updated: 2026-08-10
code:
  - docker-compose.yml
  - README.md
  - syralit.toml
  - syralit.toml.example
  - AGENTS.md
  - .env.example
-->