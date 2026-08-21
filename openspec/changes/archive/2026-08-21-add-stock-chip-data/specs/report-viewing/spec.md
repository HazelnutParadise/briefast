## ADDED Requirements

### Requirement: Chip data display

For every stock_news entry that carries a chips object, the stock news detail section SHALL render a chip block labeled 籌碼面 after the entry's watch points block (or after the summary content when no watch points block renders) and before the source links. The block SHALL contain three horizontal bar rows for foreign investors, investment trust, and dealers. Bar widths SHALL be computed server-side, scaled proportionally to the largest absolute net value among the three rows; a zero net value SHALL render its row with a zero-width bar. Net-buy bars SHALL use the bullish red accent and net-sell bars SHALL use the bearish green accent (Taiwan market convention), with no rounded corners and no shadows, consistent with the established layout language. Each row SHALL show its net value converted from shares to lots (divided by 1,000 and rounded to the nearest integer) with an explicit plus sign for positive values.

Below the bars, the block SHALL render one summary line stating the three-institutional total (in lots, same conversion), the margin balance change and short balance change (in lots, from the chips object's trading-unit values), and the chip data date. When margin_change or short_change is absent from the chips object, the corresponding fragment SHALL be omitted with no placeholder text.

When an entry has no chips object — such as reports stored before the field existed or entries whose chip data was unavailable — the chip block SHALL be omitted entirely, with no placeholder text and no rendering error. Chip block styling SHALL be defined in application code and SHALL NOT be alterable by report payload content.

#### Scenario: Chip block rendered with scaled bars

- **WHEN** a stock_news entry carries chips with foreign_net 54,758,664, trust_net -15,000, and dealer_net 2,063,215
- **THEN** the entry shows a 籌碼面 block where the foreign row renders the widest bar in the bullish red accent, the trust row renders a near-zero-width bar in the bearish green accent, and the dealer row renders a proportionally narrower bar in the bullish red accent

##### Example: bar scaling and lot conversion

| Row | Net (shares) | Bar direction and color | Relative width | Shown value |
|------|--------------|-------------------------|----------------|-------------|
| 外資 | 54,758,664 | net buy, bullish red | 100% | +54,759 張 |
| 投信 | -15,000 | net sell, bearish green | ~0.03% | -15 張 |
| 自營商 | 2,063,215 | net buy, bullish red | ~3.8% | +2,063 張 |

#### Scenario: Summary line with margin data

- **WHEN** an entry's chips carries total_net 56,806,879 shares, margin_change -18,270 lots, short_change 9,056 lots, and date 2026-08-20
- **THEN** the summary line states the total as +56,807 張, 融資 -18,270 張, 融券 +9,056 張, and the data date

#### Scenario: Missing margin fragment omitted

- **WHEN** an entry's chips carries institutional values but no margin_change and no short_change
- **THEN** the summary line shows the institutional total and the data date with no margin or short fragment and no placeholder text

#### Scenario: Legacy entry without chips

- **WHEN** a report stored before the chips field existed is rendered
- **THEN** stock news entries render without a 籌碼面 block and without placeholder content
