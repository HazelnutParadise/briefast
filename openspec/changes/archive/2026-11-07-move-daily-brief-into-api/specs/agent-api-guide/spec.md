## ADDED Requirements

### Requirement: Workflow guide endpoint

The API SHALL expose GET /api/guide returning the complete daily-brief workflow guide as `text/markdown; charset=utf-8`. The guide content SHALL be embedded in the server binary and SHALL be deployed together with the validation rules it describes. The endpoint SHALL require Bearer key authentication and SHALL re-check key status on every request. Requests with a missing, invalid, or revoked key SHALL receive a 401 JSON response with ok=false and an errors array, and SHALL be recorded in the update log with a guide-rejected-auth action. Successful reads SHALL NOT be written to the update log.

#### Scenario: Authorized agent reads the guide

- **WHEN** a request with a valid Bearer key calls GET /api/guide
- **THEN** the API responds 200 with the Markdown guide covering the trading-day gate, collection, judgement, conference handling, report composition, and POST steps

#### Scenario: Revoked key rejected and logged

- **WHEN** a request with a revoked key calls GET /api/guide
- **THEN** the API responds 401 with ok=false and records a guide-rejected-auth entry in the update log

---
### Requirement: Guide pointer on every API response

Every JSON response under /api/, except the OpenAPI document, SHALL include a top-level `guide` string telling the agent to call GET /api/guide with the same API key and read the workflow guide before calling other endpoints. Every response under /api/, including the Markdown guide itself, SHALL carry the header `Link: </api/guide>; rel="help"`. The pointer SHALL appear on success and error responses alike, including 400, 401, 404, 405, and 5xx. GET /api/openapi.json SHALL carry the pointer text in `info.description` instead, so the document stays valid OpenAPI.

#### Scenario: Error response points to the guide

- **WHEN** an agent POSTs a malformed report to /api/report
- **THEN** the 400 response carries ok=false, the errors array, and the guide pointer

#### Scenario: Success response points to the guide

- **WHEN** an agent with a valid key calls GET /api/conferences
- **THEN** the 200 response carries the listing fields and the guide pointer

---
### Requirement: API overview and OpenAPI document

The API SHALL expose GET /api returning a JSON overview that lists every public /api endpoint with its method, a one-line purpose, whether it requires a Bearer key, and a link to its operation in the OpenAPI document. The API SHALL expose GET /api/openapi.json describing every endpoint's parameters, request and response schemas, Bearer authentication, and error responses. Both endpoints SHALL be readable without a key and SHALL contain no secrets.

#### Scenario: Agent without a working key discovers the API

- **WHEN** a request without an Authorization header calls GET /api
- **THEN** the API responds 200 with the endpoint list, the guide pointer, and a link to /api/openapi.json

#### Scenario: Overview stays complete

- **WHEN** the server's /api routes are compared with the overview and the OpenAPI document
- **THEN** every route appears in both, and neither lists a route the server does not serve

---
### Requirement: Unknown API paths answer in JSON

A request to any path under /api/ that no endpoint serves SHALL receive a 404 JSON response with ok=false, an errors array naming the unknown path, and the guide pointer. It SHALL NOT fall through to the public site's HTML.

#### Scenario: Mistyped endpoint

- **WHEN** an agent calls GET /api/reports
- **THEN** the API responds 404 in JSON with the guide pointer instead of the homepage HTML
