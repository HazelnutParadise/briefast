## MODIFIED Requirements

### Requirement: Batched collection with per-batch counts

The skill SHALL split collection into four named batches, each naming the endpoints it covers and each reporting its own counts. Batch one SHALL cover the six cnyes category endpoints and report, per category, how many items were returned and how many of those were unread. Batch two SHALL cover the two TWSE datasets and report, per dataset, how many records fall inside the collection window. Batch three SHALL cover the six CTEE category feeds and report unread counts per category. Batch four SHALL cover the CNA, TechNews, and LTN feeds and report unread counts per source.

The four batches MAY run in parallel with one another, since they target independent hosts. Within batch two the two TWSE requests SHALL remain serial under the existing interval rule, because that limit applies per host rather than per batch. After all batches return, the skill SHALL require a consolidated table covering every source, and SHALL forbid starting judgement until that table is complete — an absent entry SHALL be treated as an unattempted source, not as an absence of news. A batch whose sources remain unreachable after the retry SHALL be recorded as failed in the consolidated table and handled under the source collection completeness gate: publication proceeds from what was collected, including when the failed batch is the primary source batch, and every missing source is listed in the run report.

#### Scenario: Counts reported per batch

- **WHEN** the agent runs the four batches concurrently
- **THEN** each batch reports its own per-endpoint counts and judgement waits until all four have reported

#### Scenario: Missing batch blocks judgement

- **WHEN** the consolidated table lacks an entry for the CTEE batch
- **THEN** judgement does not begin and the agent resolves the gap by fetching that batch or recording it as a failed source

#### Scenario: Primary batch failure does not stop the run

- **WHEN** the cnyes batch fails after its retry while the other batches returned successfully
- **THEN** the agent composes and publishes the report from the remaining batches, lists cnyes in its run report, and the published report content discloses nothing about the outage
