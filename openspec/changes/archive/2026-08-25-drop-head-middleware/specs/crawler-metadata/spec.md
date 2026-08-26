## MODIFIED Requirements

### Requirement: Public pages declare Traditional Chinese

Every HTML page the application serves SHALL declare the document language as `zh-Hant-TW`. This includes the admin pages, whose interface is written in Traditional Chinese.

The declaration SHALL come from the application shell configuration, so that it applies to every page the framework renders and does not depend on any per-page metadata being produced.

#### Scenario: Home page language attribute

- **WHEN** a client requests the home page
- **THEN** the response HTML root element carries `lang="zh-Hant-TW"` and no longer carries `lang="en"`

#### Scenario: History page language attribute

- **WHEN** a client requests the history list page or a dated report view
- **THEN** the response HTML root element carries `lang="zh-Hant-TW"`

#### Scenario: Admin page language attribute

- **WHEN** a client requests an admin page
- **THEN** the response HTML root element carries `lang="zh-Hant-TW"`

#### Scenario: Language applies without per-page metadata

- **WHEN** a client requests a page that produces no per-page metadata, such as an admin page
- **THEN** the response HTML root element still carries `lang="zh-Hant-TW"`

### Requirement: Each public page carries its own title and description

Public page HTML responses SHALL carry a `title` element and a `meta name="description"` element whose content identifies that specific page. The title SHALL remain the page's own title after the browser has loaded the application, and SHALL NOT be replaced by a site-wide name.

#### Scenario: Home page reflects the latest report

- **WHEN** a client requests the home page and a report exists
- **THEN** the title contains the latest report headline and the site name
- **AND** the description is derived from that report's overview text with Markdown markup removed

#### Scenario: Dated report view reflects that report

- **WHEN** a client requests the history page with a date query parameter matching a stored report
- **THEN** the title contains that report's headline and the site name

#### Scenario: History list page

- **WHEN** a client requests the history list page without a date query parameter
- **THEN** the title identifies the history listing and the site name

#### Scenario: Description length is bounded

- **WHEN** a report overview is longer than the description limit
- **THEN** the description is truncated to at most 150 characters and ends with an ellipsis

##### Example: description derivation

| Overview source | Emitted description |
|---|---|
| `## 盤前總覽\n\n台股今日開低走高。` | `台股今日開低走高。` |
| A 400-character overview | First 150 characters followed by `…` |
| Empty overview | Site default description |

#### Scenario: Browser title is not overwritten after connecting

- **WHEN** the browser finishes loading the application on a public page
- **THEN** the document title is still the title that page was served with

## REMOVED Requirements

### Requirement: Metadata rewriting never disturbs other responses

**Reason**: There is no longer a rewriting layer. Per-page metadata is produced by the framework before the shell is rendered, so no response is buffered, no body is modified, and `Content-Length`, streaming and connection upgrades are never touched by application code.

**Migration**: The guarantees this requirement protected are now structural rather than tested. Admin and API responses carry no per-page metadata because those handlers are not given a metadata source; realtime transports are unaffected because nothing intercepts responses. The behaviours that remain observable are covered by the admin scenarios under the language and metadata requirements.
