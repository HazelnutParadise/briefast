## ADDED Requirements

### Requirement: Docker compose deployment

The repo SHALL contain a Dockerfile and docker-compose.yml such that docker compose up -d --build builds the site and serves it on the configured port. The data directory SHALL be mounted as a volume so reports, keys, and the update log survive container rebuilds and restarts.

#### Scenario: Data survives restart

- **WHEN** a report and a key exist and the container is rebuilt and restarted
- **THEN** the report history and the key are still present and served

### Requirement: Configuration injection

Runtime secrets (admin password, optional overrides) SHALL be provided via environment variables or a syralit.toml file mounted at deploy time. The repo SHALL provide syralit.toml.example documenting every setting, and the real syralit.toml SHALL be excluded from version control via .gitignore.

#### Scenario: Example config provided

- **WHEN** an operator prepares a new deployment
- **THEN** copying syralit.toml.example and filling in the admin password yields a working configuration

#### Scenario: Secrets never committed

- **WHEN** git status is inspected after creating a local syralit.toml
- **THEN** syralit.toml is ignored and only syralit.toml.example is tracked
