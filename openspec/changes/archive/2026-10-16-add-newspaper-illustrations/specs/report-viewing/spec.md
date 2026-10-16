## ADDED Requirements

### Requirement: Newspaper illustrations and progressive reading

Both homepage and dated report views SHALL provide an illustrated industry guide using only the current report's industry event counts. The guide SHALL link to matching industry sections, label counts as news coverage rather than market performance, omit empty sectors, and omit itself when no events exist. Four reusable illustrations SHALL depict technology, finance, traditional industry, and real estate using etched ink, halftone, and muted slate-blue, ochre, and gray-brown washes. Unknown industries SHALL retain their escaped names and counts without an unrelated illustration. The full overview and each industry event body SHALL remain in the page behind native keyboard-operable disclosures, initially collapsed, with visible event headlines and unchanged watch items, sources, five-section order, outlook, and disclaimer. The illustration asset SHALL be embedded in the binary and served through a read-only cacheable URL with a valid image MIME type and conditional GET support. Desktop, mobile, dark mode, and print SHALL remain legible; print SHALL reveal disclosure content. Red and green SHALL retain bullish and bearish meaning, and illustrations SHALL NOT encode predictions, fabricated relationships, or fabricated values. No JSON or database change SHALL be required.

#### Scenario: Coverage proportions and navigation

- **WHEN** industries contain 2 technology events, 1 finance event, 1 traditional-industry event, and 0 real-estate events
- **THEN** the guide shows three links with counts 2, 1, and 1 and shares 50%, 25%, and 25%, omits real estate, and links each count to its own industry section

#### Scenario: Empty and unfamiliar sectors

- **WHEN** no industry has events
- **THEN** the guide renders no empty chart or illustration
- **WHEN** an unfamiliar industry has events
- **THEN** its escaped label and real count remain readable without a mismatched illustration

#### Scenario: Expand the complete report

- **WHEN** a reader activates the overview or an event headline disclosure using pointer or keyboard
- **THEN** its full original Markdown body becomes visible without navigation, additional data requests, or loss of the market outlook

#### Scenario: Embedded illustration delivery

- **WHEN** a browser requests the embedded image, repeats the request with its ETag, or attempts a POST
- **THEN** GET returns a valid WebP with cache headers, matching conditional GET returns 304, and POST returns 405

## MODIFIED Requirements

### Requirement: Market outlook display

The homepage and dated report view SHALL display an available daily market outlook after the pre-market overview headline and a disclosure containing its complete body, within the overview's main column, with a fixed label, localized direction, and rendered explanatory text. When `trajectory_chart` exists, both views SHALL render a three-stage qualitative path diagram after the direction and before the `trajectory_md` text. The diagram SHALL directly label opening, intraday, and closing levels relative to the preceding close, visibly identify itself as a scenario rather than index points, and expose an equivalent text description without hover or interaction. When `trajectory_md` exists, both views SHALL label and render it before the explanation. Bullish direction SHALL use Taiwan-market red, bearish direction SHALL use green, and `range` and `uncertain` SHALL use neutral styling. When a stored report has no outlook, neither view SHALL show an empty outlook block; when it has no trajectory or chart, neither view SHALL show an empty label or placeholder. The five fixed sections SHALL keep their existing order.

#### Scenario: Outlook chart and trajectory on both report views

- **WHEN** a report has an `up` outlook, `trajectory_chart` with `open: "above"`, `midday: "near"`, `close: "above"`, and a nonblank trajectory, and a visitor opens the homepage or dated history page
- **THEN** the overview displays its original headline and overview disclosure first, followed by `今日大盤走勢預測`, `偏多`, a three-point qualitative diagram with phase and level labels and a non-point caveat, `預期走法`, the trajectory, and the explanation in that order

#### Scenario: Old outlook without chart

- **WHEN** a stored report has direction, summary, and trajectory text but no `trajectory_chart`
- **THEN** both views render the existing direction, trajectory text, and explanation without a chart label or placeholder

#### Scenario: Old outlook without trajectory

- **WHEN** a stored report has direction and summary but no `trajectory_md` or `trajectory_chart`
- **THEN** both views render the direction and explanation without a trajectory or chart placeholder

#### Scenario: Old report without outlook

- **WHEN** a stored report has no `market_outlook`
- **THEN** the overview renders without an outlook block or placeholder
