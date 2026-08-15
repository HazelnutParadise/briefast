## ADDED Requirements

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
