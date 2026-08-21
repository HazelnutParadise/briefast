## ADDED Requirements

### Requirement: Chip reference data

The skill SHALL direct the agent to fetch previous trading day chip data (institutional investors' net buy/sell and margin trading balances) during collection, using four plain-HTTP requests: the TWSE daily institutional trading dataset rwd/zh/fund/T86 and the TWSE margin trading dataset rwd/zh/marginTrading/MI_MARGN on www.twse.com.tw, both called with the target date as a YYYYMMDD parameter and the all-stocks select type; the TPEx institutional daily trading dataset openapi/v1/tpex_3insti_daily_trading and the TPEx mainboard margin balance dataset openapi/v1/tpex_mainboard_margin_balance on www.tpex.org.tw. Requests to www.twse.com.tw SHALL follow the existing TWSE request throttling requirement as part of the same serial sequence; the TPEx host is outside that throttle.

The agent SHALL validate the trading date embedded in each response (ROC calendar format) against the expected previous trading day in Taipei time, and SHALL treat a response carrying an older date as a failed fetch rather than silently using stale chip data. The agent SHALL report, per source, how many records were returned and the embedded data date.

The skill SHALL require the agent to save each response to a file in the working directory and to read chip values only by querying those files for the symbols named in the report's stock_news entries, using file-query tools such as jq or grep. The agent SHALL NOT load a full-market chip response into the conversation context. Date validation and count reporting SHALL likewise be satisfied by querying the saved files.

When a chip source fails after one retry (TWSE retries under the throttling requirement's sixty-second rule), the agent SHALL proceed to judge and publish without that source's chip data, list the gap in its run report, and disclose nothing about the gap in report content, consistent with the source collection completeness gate's internal-disclosure rule.

#### Scenario: Chip data fetched and dated

- **WHEN** the agent fetches the four chip datasets before judgement on a normal trading morning
- **THEN** each response's embedded date matches the previous Taipei trading day and the agent reports record counts and the data date per source

#### Scenario: Stale date treated as failure

- **WHEN** the TPEx institutional response carries an embedded date older than the previous Taipei trading day after the retry
- **THEN** the agent treats that fetch as failed, composes TPEx-listed entries without institutional chip values, and lists the gap in its run report with nothing disclosed in report content

#### Scenario: TWSE chip requests join the throttle sequence

- **WHEN** the agent fetches the two www.twse.com.tw chip datasets in the same run as other TWSE requests
- **THEN** all TWSE requests run serially with at least five seconds between them and no TWSE fetches run in parallel

#### Scenario: Chip values queried from saved files by symbol

- **WHEN** the agent needs chip values for the stocks in stock_news
- **THEN** the agent queries the saved response files for those symbols only, and no full-market chip payload is read into the conversation context

### Requirement: Chip block composition

For every stock_news entry whose market's chip fetches succeeded, the skill SHALL direct the agent to fill the entry's chips block from the saved files of that stock's market (TWSE files for listed stocks, TPEx files for OTC stocks) as follows: date is the validated chip data date in YYYY-MM-DD; foreign_net is the sum of the foreign-investors-excluding-foreign-dealers net column and the foreign-dealers net column, in shares; trust_net and dealer_net are the investment trust and dealer net columns, in shares; total_net is taken directly from the source's three-institutional total column, in shares, not recomputed; margin_change and short_change are the current-day balance minus the previous-day balance from the margin dataset, in trading units (lots).

When a symbol is absent from a chip file (such as a newly listed stock or an instrument without margin trading), the agent SHALL omit the chips block or the affected fields rather than filling zeros or estimates. When a market's chip fetches failed, the agent SHALL omit the chips block for that market's entries entirely.

#### Scenario: Listed stock composed from TWSE files

- **WHEN** a listed stock appears in stock_news and both TWSE chip fetches succeeded with validated dates
- **THEN** its chips block carries the validated date, the summed foreign net, trust net, dealer net, the source's institutional total, and the margin and short balance changes computed as today minus previous day

##### Example: foreign net summation

- **GIVEN** T86 columns for a symbol: foreign-excluding-foreign-dealers net = 54,758,664 shares, foreign-dealers net = 0 shares
- **WHEN** the agent composes the chips block
- **THEN** foreign_net is 54758664

#### Scenario: Symbol absent from a chip file

- **WHEN** a stock_news symbol does not appear in the margin dataset
- **THEN** the entry's chips block omits margin_change and short_change instead of carrying zeros

#### Scenario: Failed market omits chips and still publishes

- **WHEN** both TWSE chip fetches fail after retry while TPEx fetches succeed
- **THEN** listed-stock entries carry no chips block, OTC entries carry theirs, the report is published, and the run report lists the TWSE chip gap

### Requirement: Chip citation discipline in judgement

Chip data SHALL serve as corroboration for judgements already grounded in news from the collection window, consistent with the stock-pass ranking criterion that treats fund-flow signals as corroboration only. The skill SHALL forbid creating calls or stock_news entries from chip movement alone without supporting news. When report content cites a chip figure, the citation MUST state the chip data date. Entries whose chips block is omitted SHALL be judged from news and price context alone with no placeholder language about missing chip data.

#### Scenario: No call from chip movement alone

- **WHEN** a stock shows heavy institutional net buying but no news about it exists in the collection window
- **THEN** no calls or stock_news entry is created for that stock

#### Scenario: Chip citation carries its date

- **WHEN** a stock_news summary cites institutional net buying to corroborate a bullish call
- **THEN** the cited figure states the chip data date
