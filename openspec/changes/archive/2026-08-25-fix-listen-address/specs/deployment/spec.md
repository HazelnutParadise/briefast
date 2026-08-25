## MODIFIED Requirements

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
