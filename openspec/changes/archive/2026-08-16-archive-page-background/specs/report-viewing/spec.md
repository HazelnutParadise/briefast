## MODIFIED Requirements

### Requirement: History browsing

The history page SHALL list report dates from newest to oldest, 10 per page, each row showing the date and the report headline. Selecting a date SHALL render that day's report using the same layout as the homepage.

When the selected date is older than the newest stored report date, the rendered report SHALL display a prominent historical-report notice below the masthead that names the report date, states the content is not the latest report, and links to both the homepage (latest report) and the history list. When the selected date equals the newest stored report date, the notice SHALL NOT appear.

A view of a date older than the newest stored report SHALL additionally render the entire page on a distinct archive paper background, applied by overriding the paper color token, with a dedicated value for the light theme and for the dark theme. Every existing text-on-paper combination MUST keep a contrast ratio of at least 4.5:1 on the archive background; when a text token falls below that ratio on the archive paper, the override SHALL darken or lighten that token while keeping its hue, and SHALL NOT introduce any new hue. The archive background SHALL NOT apply when the selected date equals the newest stored report date, and SHALL NOT introduce any color hue outside the neutral paper range.

A historical report view SHALL offer masthead navigation to both the homepage and the history list. The not-found state for a missing date SHALL also link back to the homepage in addition to the history list. Determining the newest stored report date MUST NOT require loading the full report payload.

#### Scenario: History list ordering and paging

- **WHEN** 25 reports exist and a visitor opens the history page
- **THEN** the 10 newest dates are listed in descending order with pagination to reach the rest

##### Example: ordering

- **GIVEN** reports dated 2026-08-05, 2026-08-07, 2026-08-06
- **WHEN** the history page renders
- **THEN** rows appear in order: 2026-08-07, 2026-08-06, 2026-08-05

#### Scenario: View a historical report

- **WHEN** a visitor selects 2026-08-05 from the history list and a newer report exists
- **THEN** the 2026-08-05 report renders with the same five-section layout as the homepage, plus a historical-report notice naming 2026-08-05 with links to the homepage and the history list

#### Scenario: Historical view uses archive paper background

- **WHEN** a visitor opens a date older than the newest stored report
- **THEN** the page overrides the paper token to the archive paper values for both light and dark themes, so the whole page is visibly distinct from the latest report at any scroll position

#### Scenario: Newest report opened from history carries no notice

- **WHEN** a visitor selects the newest stored date from the history list
- **THEN** the report renders without the historical-report notice and without the archive background override

##### Example: notice presence by date

- **GIVEN** reports dated 2026-08-05, 2026-08-06, 2026-08-07
- **WHEN** a visitor opens each date from the history list
- **THEN** 2026-08-05 and 2026-08-06 show the historical-report notice and the archive background; 2026-08-07 shows neither

#### Scenario: Missing date links home

- **WHEN** a visitor opens a history URL whose date has no stored report
- **THEN** the not-found state offers links to both the homepage and the history list
