## ADDED Requirements

### Requirement: Daily workflow skill exists in the repo

The repo SHALL contain a Cowork skill at skills/daily-brief/SKILL.md (project-root skills directory) defining the daily pre-market workflow as ordered steps: collect industry and stock market news, judge stock calls from the news, compose the report JSON, and POST it to the site API.

#### Scenario: Skill discoverable

- **WHEN** an agent loads skills/daily-brief/SKILL.md in this repo
- **THEN** the daily-brief skill is available with the four workflow steps documented

### Requirement: Skill output contract

The skill SHALL embed the full report JSON schema with a filled example and SHALL instruct the agent to produce a payload conforming to it. The skill SHALL require that every stock call carries a one-line reason and at least one news source, and that stocks with major news but no clear direction are emitted as stock_news entries with call none.

#### Scenario: Unclear direction handled

- **WHEN** the agent finds major news about a stock but cannot judge direction
- **THEN** the skill directs it to emit a stock_news entry with call none and no calls-list entry

### Requirement: Endpoint configuration via environment

The skill SHALL read the site URL from BRIEFAST_URL and the API key from BRIEFAST_API_KEY, and SHALL NOT contain a hardcoded URL or key. The skill SHALL instruct the agent to stop and report when either variable is missing.

#### Scenario: Missing credentials

- **WHEN** BRIEFAST_API_KEY is not set when the skill runs
- **THEN** the workflow stops before posting and reports the missing configuration

### Requirement: Failure reporting

The skill SHALL instruct the agent to treat a non-200 API response as a failure, include the response body in its report, and retry once on 5xx errors before reporting failure.

#### Scenario: Validation rejection surfaced

- **WHEN** the API responds 400 with an errors array
- **THEN** the agent fixes the payload per the errors when possible, otherwise reports the errors verbatim
