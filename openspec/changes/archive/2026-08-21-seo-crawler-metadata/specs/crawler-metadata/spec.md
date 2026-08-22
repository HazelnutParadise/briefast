## ADDED Requirements

### Requirement: Public pages declare Traditional Chinese

Public page HTML responses SHALL declare the document language as `zh-Hant-TW`.

#### Scenario: Home page language attribute

- **WHEN** a client requests the home page
- **THEN** the response HTML root element carries `lang="zh-Hant-TW"` and no longer carries `lang="en"`

#### Scenario: History page language attribute

- **WHEN** a client requests the history list page or a dated report view
- **THEN** the response HTML root element carries `lang="zh-Hant-TW"`

### Requirement: Each public page carries its own title and description

Public page HTML responses SHALL carry a `title` element and a `meta name="description"` element whose content identifies that specific page.

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

### Requirement: Public pages carry canonical and social sharing metadata

Public page HTML responses SHALL carry a canonical link and Open Graph and Twitter card metadata built from the resolved site base URL.

#### Scenario: Canonical points at the requested page

- **WHEN** a client requests any public page
- **THEN** the response carries a `link rel="canonical"` whose href is the absolute URL of that page

#### Scenario: Open Graph metadata is complete

- **WHEN** a client requests any public page
- **THEN** the response carries `og:type`, `og:site_name`, `og:title`, `og:description`, `og:url` and `og:locale`
- **AND** `og:url` matches the canonical href
- **AND** `og:locale` is `zh_TW`

#### Scenario: Twitter card metadata is present

- **WHEN** a client requests any public page
- **THEN** the response carries `twitter:card` with the value `summary`

### Requirement: Site serves robots.txt

The system SHALL serve `/robots.txt` describing which paths crawlers may fetch and where the sitemap lives.

#### Scenario: robots.txt content

- **WHEN** a client requests `/robots.txt`
- **THEN** the response status is 200 with content type `text/plain`
- **AND** the body disallows `/admin/` and `/api/`
- **AND** the body declares the sitemap using an absolute URL

### Requirement: Site serves a sitemap covering every report

The system SHALL serve `/sitemap.xml` listing the home page, the history list page, and every stored report.

#### Scenario: Sitemap lists stored reports

- **WHEN** a client requests `/sitemap.xml` and reports exist
- **THEN** the response status is 200 with content type `application/xml`
- **AND** the body is a `urlset` containing one entry per stored report, each with a `lastmod` value

#### Scenario: Sitemap with no reports

- **WHEN** a client requests `/sitemap.xml` and no reports exist
- **THEN** the response is a valid `urlset` containing the home page and the history list page only

#### Scenario: Sitemap query failure

- **WHEN** the report listing query fails
- **THEN** the response status is 500 and no partial XML document is written

### Requirement: Site base URL resolution

The system SHALL resolve the absolute site base URL used by canonical, Open Graph and sitemap output from configuration first and from request headers otherwise.

#### Scenario: Configured base URL wins

- **WHEN** the `BRIEFAST_SITE_URL` environment variable is set
- **THEN** absolute URLs use that value regardless of request headers

#### Scenario: Forwarded protocol header is honoured

- **WHEN** `BRIEFAST_SITE_URL` is unset and the request carries an `X-Forwarded-Proto` header
- **THEN** absolute URLs use that protocol together with the request host

#### Scenario: Direct connection fallback

- **WHEN** `BRIEFAST_SITE_URL` is unset and no `X-Forwarded-Proto` header is present
- **THEN** absolute URLs use `https` for a TLS connection and `http` otherwise, together with the request host

### Requirement: Metadata rewriting never disturbs other responses

Head rewriting SHALL apply only to successful HTML responses on public pages and SHALL leave every other response byte-identical.

#### Scenario: Realtime channels pass through

- **WHEN** a request path begins with the framework asset prefix, including the WebSocket and server-sent event endpoints
- **THEN** the response is passed through without buffering and streaming and connection upgrade remain available

#### Scenario: Non-HTML responses pass through

- **WHEN** a downstream handler responds with a content type other than `text/html`, or with a status other than 200
- **THEN** the response body is written unchanged

#### Scenario: Unrecognised HTML passes through

- **WHEN** an HTML response contains no `title` element to replace
- **THEN** the response body is written unchanged

#### Scenario: Content length stays accurate

- **WHEN** an HTML response is rewritten
- **THEN** the `Content-Length` header equals the byte length of the rewritten body

#### Scenario: Admin and API routes are untouched

- **WHEN** a client requests an admin page or an API endpoint
- **THEN** the response is identical to the response produced without head rewriting

### Requirement: Metadata failures degrade to site defaults

Report lookup failures SHALL NOT turn into page errors.

#### Scenario: Report lookup fails

- **WHEN** the report lookup backing a page's metadata returns an error or finds no report
- **THEN** the page is served with the site default title and description
- **AND** the response status is unchanged
