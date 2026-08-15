## ADDED Requirements

### Requirement: Previous-close reference data

The skill SHALL direct the agent to fetch previous trading day closing prices during collection, using two plain-HTTP requests: the TWSE listed-stock daily trading dataset exchangeReport/STOCK_DAY_ALL on openapi.twse.com.tw, and the TPEx mainboard daily close quotes dataset openapi/v1/tpex_mainboard_daily_close_quotes on www.tpex.org.tw. The TWSE request SHALL follow the existing TWSE request throttling requirement; the TPEx host is outside that throttle. The agent SHALL validate the trading date embedded in each response (ROC calendar format) against the expected previous trading day in Taipei time, and SHALL treat a response carrying an older date as a failed fetch rather than silently using stale prices. The agent SHALL report, per source, how many records were returned and the embedded data date.

Closing prices SHALL serve as judgement context — for example whether news appears already priced in, or as a magnitude baseline — and MAY be quoted in report content only with the price date stated. The skill SHALL forbid creating calls or stock_news entries from price movement alone without supporting news from the collection window. When both price fetches fail after one retry each, the agent SHALL proceed to judge and publish from news alone, list the missing price data in its run report, and disclose nothing about the gap in report content, consistent with the source collection completeness gate's internal-disclosure rule.

#### Scenario: Prices fetched and dated

- **WHEN** the agent fetches both closing-price datasets before judgement on a normal trading morning
- **THEN** each response's embedded date matches the previous Taipei trading day, the agent reports record counts and the data date per source, and judgement proceeds with prices available as context

#### Scenario: Stale date treated as failure

- **WHEN** the TWSE response carries an embedded date older than the previous Taipei trading day after the retry
- **THEN** the agent treats the TWSE price fetch as failed, judges without listed-stock prices, and lists the gap in its run report with nothing disclosed in report content

#### Scenario: No call from price movement alone

- **WHEN** a stock shows a sharp previous-day price move but no news about it exists in the collection window
- **THEN** no calls or stock_news entry is created for that stock, and the price move at most informs context inside entries already backed by news

##### Example: price quoting with date

| Situation | Allowed in report content | Not allowed |
|-----------|--------------------------|-------------|
| TSMC news entry, price date 2026-08-14 | 昨收（8/14）1,080 元，利多公布前已連漲 | 昨收 1,080 元（無日期標示） |
| Stock fell 5% on 8/14, no news in window | (no entry at all) | short_bear call citing only the 5% drop |
