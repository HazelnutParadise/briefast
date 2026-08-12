## ADDED Requirements

### Requirement: Batched collection with per-batch counts

The skill SHALL split collection into four named batches, each naming the endpoints it covers and each reporting its own counts. Batch one SHALL cover the six cnyes category endpoints and report, per category, how many items were returned and how many of those were unread. Batch two SHALL cover the two TWSE datasets and report, per dataset, how many records fall inside the collection window. Batch three SHALL cover the six CTEE category feeds and report unread counts per category. Batch four SHALL cover the CNA, TechNews, and LTN feeds and report unread counts per source.

The four batches MAY run in parallel with one another, since they target independent hosts. Within batch two the two TWSE requests SHALL remain serial under the existing interval rule, because that limit applies per host rather than per batch. After all batches return, the skill SHALL require a consolidated table covering every source, and SHALL forbid starting judgement until that table is complete — an absent entry SHALL be treated as an unattempted source, not as an absence of news. A failure of the primary source batch SHALL stop the workflow regardless of what the other batches returned.

#### Scenario: Counts reported per batch

- **WHEN** the agent runs the four batches concurrently
- **THEN** each batch reports its own per-endpoint counts and judgement waits until all four have reported

#### Scenario: Missing batch blocks judgement

- **WHEN** the consolidated table lacks an entry for the CTEE batch
- **THEN** judgement does not begin and the agent resolves the gap by fetching that batch or recording it as a failed source

#### Scenario: Primary failure stops the run

- **WHEN** the cnyes batch fails after its retry while the other batches returned successfully
- **THEN** the workflow stops without composing or posting a report

## MODIFIED Requirements

### Requirement: Source collection completeness gate

The skill SHALL require the agent to record the success or failure of every configured news source during collection and to retry each failed source once. A source that remains unreachable after the retry SHALL NOT stop publication; the skill SHALL direct the agent to publish from what it did collect and to list every missing source in its run report, including the primary source when it is among the failures. Collection mechanics SHALL stay internal: the skill SHALL forbid naming the source roster, the batch structure, or any source outage inside report content such as overview_md or any summary field, so the published report never discloses coverage gaps or how collection works. Per-article citations in the sources field are unaffected. The skill SHALL NOT treat a source that returned successfully with no articles in the window as a failure.

#### Scenario: Primary source failure stays internal

- **WHEN** cnyes remains unreachable after one retry
- **THEN** the agent composes the report from the remaining sources, names the outage only in its run report, and the published report contains no mention of the missing source or of thinner coverage

#### Scenario: Secondary source failure disclosed internally

- **WHEN** two secondary sources remain unreachable after one retry and cnyes succeeded
- **THEN** the agent publishes the report and lists both missing sources in its run report, with nothing about them in report content
