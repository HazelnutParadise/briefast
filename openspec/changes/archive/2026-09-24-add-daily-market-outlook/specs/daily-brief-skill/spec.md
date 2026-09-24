## ADDED Requirements

### Requirement: Evidence-based daily market outlook

The daily workflow SHALL include a `market_outlook` object in every new report and SHALL form its direction for the report date from the verified report evidence, including relevant news, prior-trading-day chip data when available, index-significant stocks, and relevant conference information when available. It SHALL identify countervailing evidence and data gaps in the explanation, choose `uncertain` when a reliable direction cannot be supported, and SHALL NOT fabricate point targets, probabilities, or evidence. `up` and `down` SHALL refer to the report-date TAIEX close relative to the preceding trading-day close; `range` SHALL mean choppy trading without a clear directional edge.

#### Scenario: Conflicting signals

- **WHEN** report evidence contains material opposing signals with no defensible directional edge
- **THEN** the workflow emits `direction: "uncertain"` and explains the conflict

#### Scenario: Available evidence

- **WHEN** the workflow composes a new report after collecting news, chips, and relevant conference information
- **THEN** the JSON includes a supported direction and a nonblank `summary_md` referring only to verified inputs
