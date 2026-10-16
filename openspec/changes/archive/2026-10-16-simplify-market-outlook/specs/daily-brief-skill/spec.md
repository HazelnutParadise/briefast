## MODIFIED Requirements

### Requirement: Evidence-based daily market outlook

The daily workflow SHALL include a `market_outlook` object with `direction` and nonblank `summary_md` in every new report. It SHALL omit `trajectory_md` and `trajectory_chart` from new reports and SHALL NOT forecast intraday continuation, reversals, three-phase paths, or turn times. It SHALL assess the opening tendency and closing direction independently from verified pre-report evidence. It SHALL express a supported opening tendency in `summary_md` relative to the preceding TAIEX close, without treating night futures as the opening price. When opening evidence is insufficient, it SHALL briefly state that the opening tendency is uncertain without discarding a supported closing judgment. The `direction` field SHALL describe the report-date TAIEX close relative to the preceding trading-day close: `up`, `down`, `range` (choppy trading without a clear directional edge), or `uncertain`. When closing evidence is insufficient it SHALL choose `uncertain`, while retaining any supported opening judgment. The summary SHALL explain the basis, material countervailing signals, and decisive observable conditions without repetitive statements that the intraday path cannot be judged.

Before composing the outlook, the workflow SHALL check the latest completed U.S. trading session available before generation, including dated Nasdaq Composite and SOX performance, and record sources and session dates. It SHALL obtain the completed TAIFEX TX night session attributed to the report trading date and compare it with the preceding trading date's regular close for the same contract. It SHALL use valid TX data as opening evidence and explain material agreement or divergence with U.S. markets and Taiwan-specific news, prior-trading-day chip data, index-significant stocks, and relevant conference information when available. It SHALL NOT derive the closing direction from TX or U.S. movement alone. It SHALL NOT fabricate point targets, probabilities, live observations, or evidence. Missing or unusable TX data SHALL be retried once; the workflow SHALL continue with other verified inputs and record the failure only in the operator report. Existing stored reports and API support for legacy trajectory fields SHALL remain unchanged.

#### Scenario: Opening and closing supported

- **WHEN** verified evidence supports a higher opening and a bullish close
- **THEN** the workflow emits `direction: "up"`, explains both independent judgments in `summary_md`, and omits both trajectory fields

#### Scenario: Opening supported without closing support

- **WHEN** TX supports a lower opening but other evidence cannot support a closing direction
- **THEN** the workflow emits `direction: "uncertain"`, retains the lower-opening judgment with its basis and closing uncertainty in `summary_md`, and omits both trajectory fields

#### Scenario: Closing supported without opening support

- **WHEN** evidence supports a bearish close but cannot support an opening tendency
- **THEN** the workflow emits `direction: "down"`, explains the closing basis and opening uncertainty in `summary_md`, and omits both trajectory fields

#### Scenario: Conflicting signals

- **WHEN** opposing signals support neither an opening tendency nor a closing direction
- **THEN** the workflow emits `direction: "uncertain"`, explains the conflict and decisive observation in `summary_md`, and omits both trajectory fields without an intraday-path disclaimer

#### Scenario: U.S. semiconductor signal diverges

- **WHEN** the latest completed U.S. session shows SOX and Nasdaq moving in opposite directions, or SOX conflicts with verified Taiwan-specific signals
- **THEN** the workflow checks the session dates, explains the conflicting evidence and its weighting in `summary_md`, and does not infer intraday turns from the U.S. indices

#### Scenario: U.S. market data are missing or stale

- **WHEN** a U.S. index has no verifiable close and date for the latest completed session
- **THEN** the workflow treats that index as unavailable, does not describe an older close as last night's performance, and bases the outlook on the remaining verified evidence

#### Scenario: TX night evidence supports the opening only

- **WHEN** a verified TX night result has a nonzero change and is available before report generation
- **THEN** the workflow weighs it in the opening tendency, compares it with U.S. and Taiwan evidence, and does not derive the midday or closing level from TX alone

#### Scenario: TX night evidence conflicts with other evidence

- **WHEN** the verified TX night direction conflicts materially with SOX, Nasdaq, or Taiwan-specific evidence
- **THEN** the workflow explains the divergence and its weighting in `summary_md`, uses conditional opening language in `summary_md`, and omits both trajectory fields

#### Scenario: TX night source is missing or stale

- **WHEN** TAIFEX data for the report-attributed night session are absent, unavailable before composition, or fail date and contract validation after one retry
- **THEN** the workflow omits TX as evidence, completes the report using other verified inputs, and records the TX failure only in the operator report
