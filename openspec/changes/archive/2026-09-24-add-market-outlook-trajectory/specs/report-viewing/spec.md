## MODIFIED Requirements

### Requirement: Market outlook display

The homepage and dated report view SHALL display an available daily market outlook within the pre-market overview, with a fixed label, localized direction, and rendered explanatory text. When `trajectory_md` exists, both views SHALL label and render it before the explanation. Bullish direction SHALL use Taiwan-market red, bearish direction SHALL use green, and `range` and `uncertain` SHALL use neutral styling. When a stored report has no outlook, neither view SHALL show an empty outlook block; when it has no trajectory, neither view SHALL show an empty trajectory label. The five fixed sections SHALL keep their existing order.

#### Scenario: Outlook trajectory on both report views

- **WHEN** a report has an `up` outlook and a nonblank trajectory and a visitor opens the homepage or its dated history page
- **THEN** the overview displays `今日大盤走勢預測`, `偏多`, `預期走法`, the trajectory, and the explanation in that order with red direction styling

#### Scenario: Old outlook without trajectory

- **WHEN** a stored report has direction and summary but no `trajectory_md`
- **THEN** both views render the direction and explanation without a trajectory label or placeholder

#### Scenario: Old report without outlook

- **WHEN** a stored report has no `market_outlook`
- **THEN** the overview renders without an outlook block or placeholder
