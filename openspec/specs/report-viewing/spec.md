# report-viewing Specification

## Purpose

TBD - created by archiving change 'daily-report-site'. Update Purpose after archive.

## Requirements

### Requirement: Homepage shows the latest report

The homepage SHALL render the report with the newest date found in the report store, with the masthead showing the report date and generated-at time. When no report exists, the homepage SHALL show an empty state ("尚無報告") and SHALL NOT crash.

#### Scenario: Latest report rendered

- **WHEN** reports exist for 2026-08-06 and 2026-08-07 and a visitor opens the homepage
- **THEN** the 2026-08-07 report is rendered and the masthead shows that date

#### Scenario: Empty state

- **WHEN** no report file exists and a visitor opens the homepage
- **THEN** an empty-state message is shown and the page renders without error


<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
### Requirement: Fixed report layout

The report page SHALL render exactly five content sections in this order: pre-market overview, today's watch, stock calls, industry news summary, stock news detail. Section headers, ordering, and styling SHALL be defined in application code and SHALL NOT be alterable by report payload content.

#### Scenario: Section order independent of payload

- **WHEN** a report is rendered
- **THEN** the five sections appear in the fixed order regardless of field order in the stored JSON


<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
### Requirement: Stock calls display

The stock calls section SHALL group entries into four cards: short-term bullish, short-term bearish, long-term bullish, long-term bearish. Bullish cards SHALL use red accents and bearish cards SHALL use green accents (Taiwan market convention). Each entry SHALL show the stock name, symbol, and its one-line reason.

#### Scenario: Calls grouped into four cards

- **WHEN** a report contains entries in all four call lists
- **THEN** each entry appears in its matching card with name, symbol, and reason


<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
### Requirement: Stock news detail display

The stock news detail section SHALL render one entry per stock_news item: stock name and symbol, the entry headline, a call tag matching its call value, the summary markdown, and its news source links. The headline SHALL render in the entry's identification area alongside the stock name rather than inside the body text. Entries with call value none SHALL show no tag.

#### Scenario: Tagged entry

- **WHEN** a stock_news item has call short_bull
- **THEN** its entry shows a short-term bullish tag

#### Scenario: Neutral entry without tag

- **WHEN** a stock_news item has call none
- **THEN** its entry shows the stock name and summary with no call tag

#### Scenario: Headline shown with the stock identity

- **WHEN** a stock_news item carries a headline
- **THEN** that headline appears in the entry's identification area with the stock name and symbol


<!-- @trace
source: news-headlines
updated: 2026-08-09
code:
  - internal/api/report.go
  - internal/site/site.go
  - skills/daily-brief/SKILL.md
  - internal/site/styles.go
  - internal/report/schema.go
  - internal/api/read.go
  - main.go
tests:
  - internal/api/report_test.go
  - internal/report/schema_test.go
  - internal/api/read_test.go
  - internal/site/site_test.go
  - main_test.go
-->

---
### Requirement: History browsing

The history page SHALL list report dates from newest to oldest, 10 per page, each row showing the date and the report headline. Selecting a date SHALL render that day's report using the same layout as the homepage.

#### Scenario: History list ordering and paging

- **WHEN** 25 reports exist and a visitor opens the history page
- **THEN** the 10 newest dates are listed in descending order with pagination to reach the rest

##### Example: ordering

- **GIVEN** reports dated 2026-08-05, 2026-08-07, 2026-08-06
- **WHEN** the history page renders
- **THEN** rows appear in order: 2026-08-07, 2026-08-06, 2026-08-05

#### Scenario: View a historical report

- **WHEN** a visitor selects 2026-08-05 from the history list
- **THEN** the 2026-08-05 report renders with the same five-section layout as the homepage


<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
### Requirement: Disclaimer footer

Every public report view SHALL display a fixed disclaimer footer stating the content is AI-generated from public news, is for reference only, and is not investment advice. The footer text SHALL be defined in application code, not in report payloads.

#### Scenario: Footer always present

- **WHEN** any report page (homepage or historical) is rendered
- **THEN** the disclaimer footer is displayed

<!-- @trace
source: daily-report-site
updated: 2026-08-09
code:
  - syralit.toml.example
  - internal/report/schema.go
  - internal/store/schema.sql
  - docs/design/design-demos/master.png
  - docs/design/design-demos/roulette.png
  - README.md
  - docs/design/design-demos/benchmark.png
  - go.sum
  - skills/daily-brief/SKILL.md
  - DESIGN.md
  - docs/design/fallback-spec.md
  - .dockerignore
  - main.go
  - internal/store/store.go
  - internal/admin/admin.go
  - internal/api/report.go
  - internal/site/site.go
  - .spectra/touched/daily-report-site.json
  - skills/daily-brief/scripts/seen.py
  - docker-compose.yml
  - internal/site/styles.go
  - CLAUDE.md
  - docs/design/design-demos/master.html
  - docs/design/design-demos/benchmark.html
  - go.mod
  - Dockerfile
  - AGENTS.md
  - docs/design/design-demos/roulette.html
  - .spectra/changes/daily-report-site.started
tests:
  - internal/admin/admin_test.go
  - internal/report/schema_test.go
  - internal/site/site_test.go
  - internal/store/store_test.go
  - main_test.go
  - internal/api/report_test.go
-->

---
### Requirement: Watch points display

The industry news summary section and the stock news detail section SHALL render each entry's watch_md as a visually distinct block labeled 觀察重點, placed after the entry's summary content, with styling defined in application code. When an entry has no watch_md value or an empty one — such as reports stored before the field existed — the block SHALL be omitted entirely, with no placeholder text and no rendering error.

#### Scenario: Watch points block rendered

- **WHEN** a report is rendered and a stock_news entry carries a non-empty watch_md
- **THEN** the entry shows a 觀察重點 block after its summary content

#### Scenario: Legacy report without the field

- **WHEN** a report stored before the watch_md field existed is rendered
- **THEN** industry and stock entries render without a 觀察重點 block and without placeholder content

<!-- @trace
source: daily-brief-analyst-guidance
updated: 2026-08-09
code:
  - internal/api/read.go
  - internal/report/schema.go
  - skills/daily-brief/SKILL.md
  - internal/api/report.go
  - main.go
  - internal/site/styles.go
  - internal/site/site.go
tests:
  - internal/api/report_test.go
  - main_test.go
  - internal/site/site_test.go
  - internal/api/read_test.go
  - internal/report/schema_test.go
-->

---
### Requirement: Industry event headlines display

The industry news summary section SHALL render each event of a section as its own unit: the event headline as a heading styled distinctly from body text, followed by that event's body markdown. Multiple events within one section SHALL be visually separable from each other. Headline styling SHALL be defined in application code and SHALL NOT be alterable by report payload content.

#### Scenario: Two events in one section

- **WHEN** the 科技 section carries two events
- **THEN** both event headlines render as headings above their own bodies and the two events read as separate units

<!-- @trace
source: news-headlines
updated: 2026-08-09
code:
  - internal/api/report.go
  - internal/site/site.go
  - skills/daily-brief/SKILL.md
  - internal/site/styles.go
  - internal/report/schema.go
  - internal/api/read.go
  - main.go
tests:
  - internal/api/report_test.go
  - internal/report/schema_test.go
  - internal/api/read_test.go
  - internal/site/site_test.go
  - main_test.go
-->