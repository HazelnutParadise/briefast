## ADDED Requirements

### Requirement: Watch points display

The industry news summary section and the stock news detail section SHALL render each entry's watch_md as a visually distinct block labeled 觀察重點, placed after the entry's summary content, with styling defined in application code. When an entry has no watch_md value or an empty one — such as reports stored before the field existed — the block SHALL be omitted entirely, with no placeholder text and no rendering error.

#### Scenario: Watch points block rendered

- **WHEN** a report is rendered and a stock_news entry carries a non-empty watch_md
- **THEN** the entry shows a 觀察重點 block after its summary content

#### Scenario: Legacy report without the field

- **WHEN** a report stored before the watch_md field existed is rendered
- **THEN** industry and stock entries render without a 觀察重點 block and without placeholder content
