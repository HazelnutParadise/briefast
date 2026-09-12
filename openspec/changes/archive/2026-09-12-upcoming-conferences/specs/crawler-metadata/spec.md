## MODIFIED Requirements

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

#### Scenario: Conference brief page reflects that brief

- **WHEN** a client requests the conference page with symbol and date query parameters matching a stored brief
- **THEN** the title contains the company name, the words 法說會前預測, and the site name
- **AND** the description is derived from the brief's summary text with Markdown markup removed and bounded like report descriptions

#### Scenario: Conference page for an unknown brief

- **WHEN** a client requests the conference page with parameters matching no stored brief
- **THEN** the title and description fall back to the site defaults and the canonical points at the home page

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


### Requirement: Site serves a sitemap covering every report

The system SHALL serve `/sitemap.xml` listing the home page, the history list page, every stored report, and every stored conference brief page.

#### Scenario: Sitemap lists stored reports

- **WHEN** a client requests `/sitemap.xml` and reports exist
- **THEN** the response status is 200 with content type `application/xml`
- **AND** the body is a `urlset` containing one entry per stored report, each with a `lastmod` value

#### Scenario: Sitemap lists conference brief pages

- **WHEN** a client requests `/sitemap.xml` and conference briefs exist
- **THEN** the body contains one entry per stored brief at the conference page URL with symbol and date query parameters, each with a `lastmod` taken from the brief's updated time

#### Scenario: Sitemap with no reports

- **WHEN** a client requests `/sitemap.xml` and no reports or briefs exist
- **THEN** the response is a valid `urlset` containing the home page and the history list page only

#### Scenario: Sitemap query failure

- **WHEN** the report listing query or the conference listing query fails
- **THEN** the response status is 500 and no partial XML document is written
