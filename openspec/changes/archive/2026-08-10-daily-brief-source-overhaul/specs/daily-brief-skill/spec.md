## ADDED Requirements

### Requirement: Pinned source endpoints

The skill SHALL pin every news source to a named list endpoint verified to be reachable and to yield full article text, and SHALL NOT direct the agent to discover sources from portal homepages. The pinned set SHALL be: cnyes category list API with the six categories tw_stock, headline, tech, tw_macro, cnyeshouse, and wd_stock, whose list response already carries full article content and stock tags so no detail-page fetch is made, with wd_stock feeding overview_md only; TWSE OpenAPI datasets t187ap04_L (material information) and t187ap05_L (monthly revenue) pulled whole and filtered by the collection window; CTEE rss_web livenews RSS for the six categories policy, stock, finance, industry, house, and tech, with article bodies fetched over plain HTTP and no headless browser; CNA finance RSS; TechNews site feed; and LTN business RSS. The skill SHALL NOT list UDN as a source.

#### Scenario: Collection uses pinned endpoints

- **WHEN** the agent collects news for the day
- **THEN** every list fetch targets one of the pinned endpoints and cnyes article content is taken from the list response without a detail-page fetch

#### Scenario: No headless browser required

- **WHEN** the agent collects CTEE articles
- **THEN** list discovery uses the CTEE RSS endpoints and article bodies are fetched over plain HTTP

### Requirement: TWSE request throttling

The skill SHALL require TWSE OpenAPI requests to run serially with an interval of at least five seconds between requests, and SHALL forbid parallel TWSE fetches. When TWSE responds 429 or refuses the connection, the skill SHALL direct the agent to wait sixty seconds and retry once; when the retry also fails, the agent SHALL treat TWSE as a failed secondary source under the source collection completeness gate.

#### Scenario: Rate limited then recovered

- **WHEN** a TWSE request returns 429 and the retry sixty seconds later succeeds
- **THEN** collection continues normally with no failure disclosure needed for TWSE

#### Scenario: Rate limited twice

- **WHEN** a TWSE request returns 429 and the retry sixty seconds later also fails
- **THEN** the agent proceeds without TWSE data and lists TWSE as a missing source in its run report

### Requirement: Two-pass ranked judgement

The skill SHALL split judgement into two passes run in order: an industry pass using the industry analyst perspective, then a stock pass using the equity analyst perspective. Before writing in each pass, the skill SHALL require the agent to rank the window's news by importance using that pass's criteria, compared in listed order. Industry-pass criteria: (1) supply-demand or price shifts that already happened, (2) policy or regulatory measures about to take binding effect, (3) breadth of impact across the chain, (4) durability of the shift, (5) certainty of the information. Stock-pass criteria: (1) quantifiable earnings impact already visible, (2) factual tier — company disclosure above company outlook above analyst views above supply-chain chatter above rumor, (3) how soon the impact shows up in results, (4) irreversible change in competitive position, (5) fund-flow signals as corroboration only. The ranking SHALL govern only which events lead a section, how much depth each item receives, and ordering within summaries; the skill SHALL forbid using the ranking to cap or exclude items from the report, and the no-quota rule SHALL remain in force.

#### Scenario: Ranking shapes emphasis not inclusion

- **WHEN** the industry pass ranks a confirmed contract-price increase above an unconfirmed capacity rumor
- **THEN** the price increase leads its section with more depth while the rumor still appears with appropriate hedging, and neither is dropped because of rank

#### Scenario: Same news ranked differently per pass

- **WHEN** a policy announcement ranks high in the industry pass but has no quantifiable earnings impact for a specific stock
- **THEN** the industry section leads with it while the stock pass gives it lower emphasis in the affected entries

### Requirement: Analyst summary writing craft

The skill SHALL include writing-craft rules for both judgement passes, all bound by two guards: information may come only from the collected news or data fetched during the current run, and when the required material is absent the sentence is omitted entirely rather than replaced with a placeholder such as a cannot-estimate phrase. Industry-pass craft: figures SHALL carry direction, magnitude, period, and a comparison baseline only as stated in the news; price or supply-demand shifts SHALL be attributed to the supply side or the demand side together with the durability implication, and when the news does not state the driver the summary SHALL say the driver is unclear rather than invent one — this is the only rule where missing material is stated instead of omitted; summaries SHALL describe impact distribution as who benefits, who absorbs cost, and who can pass it through; events covered in depth SHALL state whether the shift is cyclical or structural, consistent with the durability ranking criterion. Stock-pass craft: each entry SHALL open with what the news changes about the judgement, one level deeper than the entry headline, and entries with call none SHALL open with why direction cannot be judged; magnitude conversion against company financials SHALL appear only when the base figures were fetched during the current run; expectation-gap statements SHALL rest only on expectations stated in the news or on events already registered in the dedup script; promotional wording SHALL be translated into facts, with unquantified deals downgraded to stated intent; attribution SHALL name sources only as the news names them, keeping vague subjects vague with the certainty tier lowered accordingly. The skill SHALL mark its example table as illustrative, not as sentence templates to copy.

#### Scenario: Unstated driver acknowledged not invented

- **WHEN** the news reports a price rise without stating whether supply cuts or demand pull drives it
- **THEN** the event summary marks the driver as unclear instead of asserting a cause

#### Scenario: Magnitude omitted without placeholder

- **WHEN** a stock entry cites a contract win but no base revenue figure was fetched this run
- **THEN** the entry carries no magnitude comparison and no cannot-estimate placeholder sentence

#### Scenario: None entry opens with the uncertainty

- **WHEN** a stock entry carries call none
- **THEN** its summary opens by stating why the direction cannot be judged from the day's news

## MODIFIED Requirements

### Requirement: Endpoint configuration via environment

The skill SHALL direct the agent to read the site URL from the BRIEFAST_URL key and the API key from the BRIEFAST_API_KEY key of a .env file at the execution workspace root, in KEY=VALUE format, and SHALL NOT contain a hardcoded URL or key. The skill directory SHALL ship a .env.example template listing the two keys with placeholder values and no real secrets. The skill SHALL instruct the agent to stop and report when the .env file is missing or either key is missing or blank.

#### Scenario: Missing credentials

- **WHEN** the workspace root has no .env file, or BRIEFAST_API_KEY is missing or blank in it, when the skill runs
- **THEN** the workflow stops before collecting and reports the missing configuration

#### Scenario: Template present

- **WHEN** a user prepares a new execution workspace
- **THEN** the skill's .env.example shows the two required keys with placeholder values to copy into the workspace .env
