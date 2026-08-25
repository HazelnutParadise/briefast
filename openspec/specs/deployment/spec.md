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

Non-secret runtime configuration SHALL live in a single syralit.toml file tracked in version control and built into the container image, and the repo SHALL NOT carry a separate example copy of it. The tracked syralit.toml SHALL NOT contain a secrets section. The image SHALL run correctly with no configuration file mounted, binding the address declared in the baked-in file rather than falling back to a loopback default. Deploy-time injection SHALL be limited to secrets and the host port, which SHALL live in a .env file at the repo root that is excluded from version control, with a tracked .env.example listing every key with its purpose and no real values. The admin password SHALL be supplied through the BRIEFAST_ADMIN_PASSWORD key of that .env, resolved as an environment variable at runtime. Docker Compose SHALL take BRIEFAST_PORT, BRIEFAST_ADMIN_PASSWORD, and BRIEFAST_DATA from the same .env through its native variable substitution, requiring no extra flags. The only mount SHALL be the data directory, whose source SHALL default to a named volume and SHALL become a host path when BRIEFAST_DATA names one. The database path inside the container SHALL remain fixed and SHALL NOT be exposed as a deploy-time variable.

The address the server listens on SHALL be derived from the resolved configuration, meaning the tracked syralit.toml with the framework defaults applied to any field it leaves unset. It SHALL NOT be read from a source that is only populated while a page function is executing.

#### Scenario: Image runs standalone

- **WHEN** the image is started with no configuration file mounted
- **THEN** the service binds the address from the baked-in syralit.toml and answers requests from outside the container

#### Scenario: Listen address is resolved from configuration

- **WHEN** the application computes the address it will listen on, outside of any page function
- **THEN** the address is the host and port declared in the tracked syralit.toml
- **AND** neither part is an empty string nor zero

##### Example: address resolution

| syralit.toml | Resolved listen address |
|---|---|
| `host = "0.0.0.0"`, `port = 8600` | `0.0.0.0:8600` |
| neither key present | the framework defaults, never `:0` |

#### Scenario: Only the database is persisted

- **WHEN** the Compose configuration is expanded with no BRIEFAST_DATA set
- **THEN** the only mount is the named database volume at the data directory, with no bind mount for the configuration file

#### Scenario: Data directory placed on the host

- **WHEN** BRIEFAST_DATA in .env names a host path
- **THEN** the data directory is bind mounted from that path and the container database path is unchanged

#### Scenario: Secrets never committed

- **WHEN** git status is inspected after configuring an admin password
- **THEN** the password appears in no tracked file, .env is ignored, and the tracked syralit.toml carries no secrets section at all


<!-- @trace
source: fix-listen-address
updated: 2026-08-25
code:
  - main.go
  - go.sum
  - go.mod
tests:
  - main_test.go
-->

---
### Requirement: Data directory works on any host mount

The container SHALL start successfully when the data directory is bind mounted from a host path, regardless of whether that path already exists or which user owns it, and SHALL require no manual ownership change on the host. The image SHALL achieve this by entering an entrypoint as root that creates the data directory and gives it to the application user, then dropping privileges so the application process itself never runs as root. When the container is started as a non-root user, the entrypoint SHALL skip the ownership step and execute the application directly rather than failing. The data directory SHALL be derived from the existing database path setting, with no additional configuration variable.

#### Scenario: Fresh host path

- **WHEN** the data mount points at a host path that does not yet exist
- **THEN** the container starts, the database opens, and the service answers requests without any host-side chown

#### Scenario: Host path owned by another user

- **WHEN** the data mount points at an existing host path owned by a different user than the application user
- **THEN** the entrypoint takes ownership of it and the database opens normally

#### Scenario: Application does not run as root

- **WHEN** the container is running normally
- **THEN** the application process runs as the unprivileged application user

<!-- @trace
source: bind-mount-safe-data-dir
updated: 2026-08-15
code:
  - README.md
  - Dockerfile
  - docker-entrypoint.sh
-->