## MODIFIED Requirements

### Requirement: Endpoint configuration via environment

The guide SHALL direct the agent to read the site URL from BRIEFAST_URL and the API key from BRIEFAST_API_KEY. For each key, the agent SHALL first use the process environment variable of that name; when it is unset or blank, the agent SHALL use the same key from a .env file at the execution workspace root in KEY=VALUE format. The guide SHALL NOT contain a hardcoded URL or key. The README SHALL show a workspace .env template listing the two keys with placeholder values and no real secrets, and SHALL state that a cloud schedule only needs the two environment variables. The guide and the README task instruction SHALL tell the agent to stop and report which key is missing when neither source provides a non-blank value for it.

#### Scenario: Missing credentials

- **WHEN** BRIEFAST_API_KEY is unset or blank in the environment and the workspace root has no .env file or no non-blank BRIEFAST_API_KEY in it, when the workflow runs
- **THEN** the workflow stops before collecting and reports the missing configuration

#### Scenario: Template present

- **WHEN** a user prepares a new execution workspace
- **THEN** the README shows the two required keys with placeholder values to copy into the workspace .env

#### Scenario: Cloud schedule with environment variables only

- **WHEN** a cloud schedule sets BRIEFAST_URL and BRIEFAST_API_KEY as environment variables and the workspace has no .env file
- **THEN** the workflow uses the environment values and proceeds without stopping

#### Scenario: Environment value takes precedence

- **WHEN** BRIEFAST_URL is set in both the environment and the workspace .env with different values
- **THEN** the agent uses the environment value
