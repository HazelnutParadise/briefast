## ADDED Requirements

### Requirement: Industry event headlines display

The industry news summary section SHALL render each event of a section as its own unit: the event headline as a heading styled distinctly from body text, followed by that event's body markdown. Multiple events within one section SHALL be visually separable from each other. Headline styling SHALL be defined in application code and SHALL NOT be alterable by report payload content.

#### Scenario: Two events in one section

- **WHEN** the 科技 section carries two events
- **THEN** both event headlines render as headings above their own bodies and the two events read as separate units

## MODIFIED Requirements

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
