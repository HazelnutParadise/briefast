## MODIFIED Requirements

### Requirement: Public pages declare Traditional Chinese

Every HTML page the application serves SHALL declare the document language as `zh-Hant-TW`. This includes the admin pages, whose interface is written in Traditional Chinese.

The declaration SHALL come from the application shell itself rather than from the metadata rewriting layer, so that it survives a response the rewriting layer declines to modify.

#### Scenario: Home page language attribute

- **WHEN** a client requests the home page
- **THEN** the response HTML root element carries `lang="zh-Hant-TW"` and no longer carries `lang="en"`

#### Scenario: History page language attribute

- **WHEN** a client requests the history list page or a dated report view
- **THEN** the response HTML root element carries `lang="zh-Hant-TW"`

#### Scenario: Admin page language attribute

- **WHEN** a client requests an admin page
- **THEN** the response HTML root element carries `lang="zh-Hant-TW"`

#### Scenario: Language survives a declined rewrite

- **WHEN** an HTML response contains no `title` element, so the metadata rewriting layer passes the body through unchanged
- **THEN** the response HTML root element still carries `lang="zh-Hant-TW"`

#### Scenario: Rewriting layer leaves the language attribute alone

- **WHEN** the metadata rewriting layer rewrites a page's head
- **THEN** it replaces only the `title` element and the metadata tags, and performs no substitution on the `lang` attribute
