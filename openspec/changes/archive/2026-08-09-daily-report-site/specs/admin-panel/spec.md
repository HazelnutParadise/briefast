## ADDED Requirements

### Requirement: Admin login gate

The /admin page SHALL show only a login form until the operator authenticates with the admin password, which SHALL be read from configuration (syralit.toml secrets or the BRIEFAST_ADMIN_PASSWORD environment variable). When no admin password is configured, the admin page SHALL refuse to serve admin functions and SHALL display a configuration error instead.

#### Scenario: Unauthenticated visitor

- **WHEN** a visitor opens /admin without logging in
- **THEN** only the login form is shown and no key or log data is rendered

#### Scenario: Missing password configuration

- **WHEN** no admin password is configured and /admin is opened
- **THEN** a configuration error is shown and login is not possible

### Requirement: API key creation

An authenticated admin SHALL be able to create an API key by entering a name. The system SHALL generate a random token, persist the key with its name and creation time to the key store, and show it in the key list immediately.

#### Scenario: Create a named key

- **WHEN** the admin creates a key named cowork-daily
- **THEN** a random token is generated and the key appears in the list with its name and creation time

### Requirement: API keys remain fully viewable

The key list SHALL display the complete plaintext token of every non-revoked key on every visit, not only at creation time. Tokens SHALL therefore be stored in plaintext in the key store.

#### Scenario: Token visible on a later visit

- **WHEN** the admin logs in again days after creating a key
- **THEN** the full token of that key is still displayed in the list

### Requirement: API key revocation

An authenticated admin SHALL be able to revoke a key. Revocation SHALL persist a revocation timestamp and take effect immediately for API authentication. Revoked keys SHALL be visually distinguished from active keys in the list.

#### Scenario: Revoke a key

- **WHEN** the admin revokes a key
- **THEN** the key is marked revoked in the list and subsequent API requests using it receive 401

### Requirement: Report update log view

The admin panel SHALL display the report update log from newest to oldest, each entry showing the time, the key name used, the report date, and the action.

#### Scenario: Log entries listed

- **WHEN** two reports were ingested and the admin opens the update log
- **THEN** both entries are listed newest first with time, key name, and report date
