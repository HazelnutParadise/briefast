## MODIFIED Requirements

### Requirement: Market outlook display

The homepage and dated report view SHALL display an available daily market outlook within the pre-market overview, with a fixed label, localized direction, and rendered explanatory text. When `trajectory_chart` exists, both views SHALL render a three-stage qualitative path diagram after the direction and before the `trajectory_md` text. The diagram SHALL directly label opening, intraday, and closing levels relative to the preceding close, visibly identify itself as a scenario rather than index points, and expose an equivalent text description without hover or interaction. When `trajectory_md` exists, both views SHALL label and render it before the explanation. Bullish direction SHALL use Taiwan-market red, bearish direction SHALL use green, and `range` and `uncertain` SHALL use neutral styling. When a stored report has no outlook, neither view SHALL show an empty outlook block; when it has no trajectory or chart, neither view SHALL show an empty label or placeholder. The five fixed sections SHALL keep their existing order.

#### Scenario: Outlook chart and trajectory on both report views

- **WHEN** a report has an `up` outlook, `trajectory_chart` with `open: "above"`, `midday: "near"`, `close: "above"`, and a nonblank trajectory, and a visitor opens the homepage or dated history page
- **THEN** the overview displays `今日大盤走勢預測`, `偏多`, a three-point qualitative diagram with phase and level labels and a non-point caveat, `預期走法`, the trajectory, and the explanation in that order

#### Scenario: Old outlook without chart

- **WHEN** a stored report has direction, summary, and trajectory text but no `trajectory_chart`
- **THEN** both views render the existing direction, trajectory text, and explanation without a chart label or placeholder

#### Scenario: Old outlook without trajectory

- **WHEN** a stored report has direction and summary but no `trajectory_md` or `trajectory_chart`
- **THEN** both views render the direction and explanation without a trajectory or chart placeholder

#### Scenario: Old report without outlook

- **WHEN** a stored report has no `market_outlook`
- **THEN** the overview renders without an outlook block or placeholder
