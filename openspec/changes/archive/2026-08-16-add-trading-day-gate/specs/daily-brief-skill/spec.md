## ADDED Requirements

### Requirement: Trading-day gate before collection

The skill SHALL determine, after reading `.env` and before any news collection begins, whether the current date in the Asia/Taipei timezone is a Taiwan stock market trading day. On a non-trading day the run SHALL end immediately: no collection, no judgement, no report composition, and no POST. The operator-facing run report SHALL state the closure reason (weekend, or the matching calendar entry name).

The determination MUST be based on verifiable rules, never on model memory or general knowledge: Saturday and Sunday SHALL always be treated as non-trading days; any other date SHALL be checked against the TWSE OpenAPI holiday schedule endpoint `https://openapi.twse.com.tw/v1/holidaySchedule/holidaySchedule`, whose entries carry dates in ROC format (e.g. `1150101` for 2026-01-01).

Calendar entries MUST be classified before use: an entry whose name or description states the market is closed (「市場無交易」, 「放假」, 「補假」) marks a non-trading day; informational entries that mark trading days (「開始交易日」, 「最後交易日」) MUST NOT be treated as closures. A date with no matching closure entry on a weekday SHALL be treated as a trading day.

If the holiday schedule fetch fails on a weekday, the skill SHALL retry once; if it still fails, the skill SHALL proceed with the full run as if the date were a trading day and SHALL list the endpoint failure in the operator-facing run report. The weekend rule requires no endpoint and SHALL apply regardless of endpoint availability. The closure decision and gate mechanics SHALL NOT appear in report content, consistent with the existing rule that collection mechanics stay internal.

#### Scenario: Weekend run ends without publishing

- **WHEN** the skill runs on a Saturday or Sunday in the Asia/Taipei timezone
- **THEN** the run ends before collection, nothing is POSTed, and the run report states the weekend closure — without calling the holiday schedule endpoint

#### Scenario: Holiday matched in calendar ends the run

- **WHEN** the current Taipei date matches a holiday schedule entry whose name or description states the market is closed
- **THEN** the run ends before collection, nothing is POSTed, and the run report cites the entry name as the closure reason

#### Scenario: Informational calendar entry does not close the market

- **WHEN** the current Taipei date matches only an informational entry such as a first-or-last-trading-day notice
- **THEN** the date is treated as a trading day and the run proceeds normally

##### Example: entry classification

| Entry name | Classification |
| --- | --- |
| 市場無交易，僅辦理結算交割作業 | non-trading day |
| 農曆除夕及春節 (description contains 放假) | non-trading day |
| 國曆新年開始交易日 | trading day |
| 農曆春節前最後交易日 | trading day |

#### Scenario: Calendar endpoint failure on a weekday

- **WHEN** the holiday schedule endpoint fails twice (initial attempt plus one retry) on a weekday
- **THEN** the run proceeds as a trading day and the endpoint failure is listed in the operator-facing run report, not in report content
