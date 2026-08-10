## MODIFIED Requirements

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
