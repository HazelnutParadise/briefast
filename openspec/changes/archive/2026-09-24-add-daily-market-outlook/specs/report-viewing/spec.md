## ADDED Requirements

### Requirement: Market outlook display

The homepage and dated report view SHALL display an available daily market outlook within the pre-market overview, with a fixed label, localized direction, and rendered explanatory text. Bullish direction SHALL use Taiwan-market red, bearish direction SHALL use green, and `range` and `uncertain` SHALL use neutral styling. When a stored report has no outlook, neither view SHALL show an empty outlook block. The five fixed sections SHALL keep their existing order.

#### Scenario: Outlook on both report views

- **WHEN** a report has an `up` outlook and a visitor opens the homepage or its dated history page
- **THEN** the overview displays `今日大盤走勢預測`, `偏多`, and the explanation with red direction styling

#### Scenario: Old report without outlook

- **WHEN** a stored report has no `market_outlook`
- **THEN** the overview renders without an outlook block or placeholder
