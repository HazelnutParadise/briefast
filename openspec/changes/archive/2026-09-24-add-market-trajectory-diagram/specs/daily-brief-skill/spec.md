## MODIFIED Requirements

### Requirement: Evidence-based daily market outlook

The daily workflow SHALL include a `market_outlook` object with `direction`, `trajectory_md`, and `summary_md` in every new report. It SHALL form the report-date direction and a baseline opening-to-intraday-to-close trajectory from verified report evidence, including relevant news, prior-trading-day chip data when available, index-significant stocks, and relevant conference information when available. The trajectory SHALL state at least one observable condition that would alter the baseline. When evidence supports all three phases, the workflow SHALL include `trajectory_chart` with `open`, `midday`, and `close` qualitative levels relative to the preceding trading-day close, consistent with `trajectory_md` and the closing `direction`. If evidence cannot support a reliable intraday path, it SHALL explicitly say so, name the decisive condition to watch, and omit `trajectory_chart` rather than invent an intraday turn. The explanation SHALL identify countervailing evidence and data gaps, choose `uncertain` when a reliable direction cannot be supported, and SHALL NOT fabricate point targets, probabilities, precise turn times, live market observations, or evidence. `up` and `down` SHALL refer to the report-date TAIEX close relative to the preceding trading-day close; `range` SHALL mean choppy trading without a clear directional edge.

#### Scenario: Conflicting signals

- **WHEN** report evidence contains material opposing signals with no defensible directional edge
- **THEN** the workflow emits `direction: "uncertain"`, omits `trajectory_chart`, explains the conflict, and writes a trajectory that says the intraday path cannot be reliably judged and identifies the key condition to monitor

#### Scenario: Available evidence

- **WHEN** the workflow composes a new report after collecting news, chips, and relevant conference information sufficient to support opening, intraday, and closing phases
- **THEN** the JSON includes a supported direction, a nonblank baseline trajectory covering all three phases plus a condition that would change it, a consistent three-phase `trajectory_chart`, and a nonblank `summary_md` referring only to verified inputs

#### Scenario: Direction support without reliable path

- **WHEN** evidence supports a closing direction but not a defensible three-phase path
- **THEN** the workflow retains the supported direction, states in `trajectory_md` that the intraday path cannot be reliably judged and identifies the decisive condition, and omits `trajectory_chart`
