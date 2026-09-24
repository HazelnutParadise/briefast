## ADDED Requirements

### Requirement: Verified TX night extraction

The daily skill SHALL provide a deterministic local parser for two saved TAIFEX TX daily CSV files. It SHALL accept a report trading date and the verified preceding Taiwan trading date, select the unique highest-volume nonzero single-month TX `盤後` row attributed to the report date, and compare its closing price with the unique `一般` row for the same contract on the preceding trading date. It SHALL reject stale dates, malformed or ambiguous rows, missing matching contracts, and nonpositive or nonfinite prices. On success it SHALL emit one JSON result containing both dates, the contract, both closing prices, point and percentage changes, night volume, and `up`, `down`, or `flat` direction. On failure it SHALL exit nonzero without emitting a partial result.

#### Scenario: Completed night session

- **WHEN** the report date is 2026-09-24, the preceding trading date is 2026-09-23, the highest-volume report-date TX night contract is 202610 with close 47909 and volume 28947, and the preceding regular session for 202610 closed at 48336
- **THEN** the parser emits contract 202610, `change_points: -427`, and `direction: "down"`, using no report-date regular-session price

#### Scenario: Weekend and contract rollover

- **WHEN** the report date follows a weekend or contract rollover and the supplied preceding trading date is several calendar days earlier
- **THEN** the parser pairs the report-date night close only with the matching contract on the supplied preceding trading date, without subtracting one calendar day or substituting another contract

#### Scenario: Unusable source rows

- **WHEN** the report-date night row is absent, two eligible contracts tie for highest volume, the preceding same-contract regular row is missing, or a required date, price, volume, or CSV header is invalid
- **THEN** the parser exits nonzero, describes the failure on stderr, and emits no JSON on stdout

#### Scenario: Flat overnight move

- **WHEN** the two valid same-contract closing prices are equal
- **THEN** the parser emits zero point and percentage change with `direction: "flat"`, without forcing a bullish or bearish signal

## MODIFIED Requirements

### Requirement: Evidence-based daily market outlook

The daily workflow SHALL include a `market_outlook` object with `direction`, `trajectory_md`, and `summary_md` in every new report. Before composing it, the workflow SHALL check the latest completed U.S. trading session available before report generation, including dated Nasdaq Composite and PHLX Semiconductor Sector Index (SOX) performance, and record each source and session date. It SHALL also obtain the completed TAIFEX TX night session attributed to the report trading date and compare it with the preceding trading date's regular close for the same contract, using only verified pre-report data. It SHALL use a valid TX result as evidence for the opening tendency, and SHALL explain material agreement or divergence with SOX, the broader U.S. market, and Taiwan-specific evidence. It SHALL NOT treat TX or U.S. movement alone as a deterministic closing forecast or a three-phase TAIEX path. It SHALL form the report-date direction and a baseline opening-to-intraday-to-close trajectory from verified report evidence, including relevant news, prior-trading-day chip data when available, index-significant stocks, and relevant conference information when available. The trajectory SHALL state at least one observable condition that would alter the baseline. When evidence supports all three phases, the workflow SHALL include `trajectory_chart` with `open`, `midday`, and `close` qualitative levels relative to the preceding trading-day close, consistent with `trajectory_md` and the closing `direction`. If evidence cannot support a reliable intraday path, it SHALL explicitly say so, name the decisive condition to watch, and omit `trajectory_chart` rather than invent an intraday turn. The explanation SHALL identify countervailing evidence and data gaps, choose `uncertain` when a reliable direction cannot be supported, and SHALL NOT fabricate point targets, probabilities, precise turn times, live market observations, or evidence. `up` and `down` SHALL refer to the report-date TAIEX close relative to the preceding trading-day close; `range` SHALL mean choppy trading without a clear directional edge. When verified TX night data are missing or unusable, the workflow SHALL retry once, continue with other verified evidence, and record the TX failure only in the operator report.

#### Scenario: Conflicting signals

- **WHEN** report evidence contains material opposing signals with no defensible directional edge
- **THEN** the workflow emits `direction: "uncertain"`, omits `trajectory_chart`, explains the conflict, and writes a trajectory that says the intraday path cannot be reliably judged and identifies the key condition to monitor

#### Scenario: Available evidence

- **WHEN** the workflow composes a new report after collecting news, chips, and relevant conference information sufficient to support opening, intraday, and closing phases
- **THEN** the JSON includes a supported direction, a nonblank baseline trajectory covering all three phases plus a condition that would change it, a consistent three-phase `trajectory_chart`, and a nonblank `summary_md` referring only to verified inputs

#### Scenario: Direction support without reliable path

- **WHEN** evidence supports a closing direction but not a defensible three-phase path
- **THEN** the workflow retains the supported direction, states in `trajectory_md` that the intraday path cannot be reliably judged and identifies the decisive condition, and omits `trajectory_chart`

#### Scenario: U.S. semiconductor signal diverges

- **WHEN** the latest completed U.S. session shows SOX and Nasdaq moving in opposite directions, or SOX conflicts with verified Taiwan-specific signals
- **THEN** the workflow checks the session dates, explains the conflicting evidence and its weighting in `summary_md`, and does not infer a three-phase TAIEX path from the U.S. indices alone

#### Scenario: U.S. market data are missing or stale

- **WHEN** a U.S. index has no verifiable close and date for the latest completed session
- **THEN** the workflow treats that index as unavailable, does not describe an older close as last night's performance, and bases the outlook on the remaining verified evidence

#### Scenario: TX night evidence supports the opening only

- **WHEN** a verified TX night result has a nonzero change and is available before report generation
- **THEN** the workflow weighs it in the opening tendency, compares it with U.S. and Taiwan evidence, and does not derive the midday or closing level from TX alone

#### Scenario: TX night evidence conflicts with other evidence

- **WHEN** the verified TX night direction conflicts materially with SOX, Nasdaq, or Taiwan-specific evidence
- **THEN** the workflow explains the divergence and its weighting in `summary_md`, uses conditional opening language in `trajectory_md`, and omits `trajectory_chart` unless all three phases have independent support

#### Scenario: TX night source is missing or stale

- **WHEN** TAIFEX data for the report-attributed night session are absent, unavailable before composition, or fail date and contract validation after one retry
- **THEN** the workflow omits TX as evidence, completes the report using other verified inputs, and records the TX failure only in the operator report
