## MODIFIED Requirements

### Requirement: Conference announcement detection

The skill SHALL direct the agent, after the five collection batches complete, to scan the TWSE material-information dataset opendata/t187ap04_L and the TPEx dataset openapi/v1/mopsfin_t187ap04_O already fetched in batch 2 for entries whose 符合條款 field is 第12款, and to parse from each entry's 說明 text the conference date (the value after 1.召開法人說明會之日期：, ROC calendar converted to YYYY-MM-DD), the time (after 2.召開法人說明會之時間：, normalised to HH:MM or empty), the venue (after 3.召開法人說明會之地點：), and the summary message (after 4.法人說明會擇要訊息：). The agent SHALL report how many 第12款 entries were found per market. Conference detection SHALL NOT be restricted to the news collection time window when an announced conference remains in the future.

The skill SHALL classify as invited an entry whose subject, summary, or official calendar other-notes explicitly states that the company was invited to a broker-hosted conference or explicitly states that a broker hosts that conference. The words 參加, 論壇, Summit, and Conference alone, or a broker name appearing only in the venue, SHALL NOT suffice to classify a conference as invited. When the evidence is ambiguous or contradictory, the agent SHALL inspect the original official announcement; if classification remains unresolved, it SHALL report the entry as unresolved and SHALL NOT register it. Date ranges SHALL NOT be registered as single-date conferences. Only confirmed company-held conferences whose date is on or after the current Taipei date SHALL be registered. The agent SHALL query GET /api/conferences first and SHALL NOT re-post a brief for a (symbol, held_on) pair already listed in upcoming unless the parsed date, time, or venue differs from the stored values.

#### Scenario: Company-held conference registered

- **WHEN** a 第12款 entry for 1314 reads 說明本公司2026年第二季財務暨營運概況 with date 115/09/16 and no invitation evidence
- **THEN** the agent classifies it as company-held and proceeds to compose and post a brief for (1314, 2026-09-16) if that date has not passed in Taipei

#### Scenario: Invited forum skipped

- **WHEN** a 第12款 entry reads 本公司受邀參加瑞銀證券舉辦之投資論壇
- **THEN** the agent classifies it as an invited appearance, registers nothing, and counts it in the run report

##### Example: classification table

| Subject or summary | Date field | Classification |
|--------------------|------------|----------------|
| 說明本公司2026年第二季財務暨營運概況 | 115/09/16 | company-held if not past, register |
| 本公司受邀參加群益證券舉辦之法人說明會 | 115/09/11 | invited, skip |
| 公告本公司將於115年09月11日舉行法人說明會 | 115/09/11 | company-held if not past, register |
| 本公司舉辦第三季線上法說會，參加方式見官網 | 115/10/12 | company-held, register |
| 德信綜合證券舉辦之法人說明會 | 115/09/29 | broker-hosted, skip |
| 2026/9/11受邀參加CLSA; 2026/9/16受邀參加KeyBanc | 115/09/11 ~ 115/10/13 | invited and range, skip |

#### Scenario: Already registered conference not re-posted

- **WHEN** GET /api/conferences lists (2330, 2026-10-16) in upcoming with the same date, time, and venue as the parsed entry
- **THEN** the agent posts nothing for that pair

#### Scenario: Ambiguous organizer reported

- **WHEN** a conference lists a broker only as venue and neither the subject, summary, nor other-notes names its organizer
- **THEN** the agent checks the original official announcement and, if organizer evidence is still absent, reports the entry as unresolved without registering it

## ADDED Requirements

### Requirement: Official conference calendar reconciliation

On every trading-day run, the skill SHALL direct the agent to fetch the official MOPS 法人說明會一覽表 for listed and OTC companies for the current Taipei month and the next month. It SHALL submit the form fields step=1, firstin=1, off=1, TYPEK=sii or otc, ROC year, two-digit month, and empty co_id to the official ajax_t100sb02_1 endpoint. It SHALL parse the myTable rows and compare future single-date rows against GET /api/conferences by market, symbol, and date, deduplicating repeated official rows for the same key. Confirmed company-held rows missing from upcoming SHALL enter the existing brief composition and POST flow even when their announcements are older than the news window or absent from the fetched major-information datasets. The agent SHALL report row counts for each of the four requests, missing confirmed company-held entries, invited entries, unresolved entries, and each fetch or parse failure. A failed request or missing table SHALL NOT be interpreted as zero scheduled conferences.

#### Scenario: Older announcement recovered from calendar

- **WHEN** the OTC calendar contains 星亞 (7753) for 115/10/14 with summary 本公司舉辦財務業務概況座談會, the run's major-information data contains the same 第12款 announcement, and upcoming lacks (7753, 2026-10-14)
- **THEN** the agent treats the entry as company-held and composes a brief for the existing POST /api/conference flow regardless of the announcement's age

#### Scenario: Wording about attendance does not exclude self-held meeting

- **WHEN** the listed calendar contains 南亞科 (2408) for 115/10/12 with summary 2026年第3季營運狀況說明 and 參加方式：請參見官網
- **THEN** the agent does not classify the entry as invited merely because 參加方式 appears

#### Scenario: Broker invitation stays excluded

- **WHEN** a calendar row's other-notes state 本公司受邀參加永豐證券法人說明會
- **THEN** the agent excludes that row from the brief flow even if its summary only says 本公司營運及財務業務狀況說明

#### Scenario: Calendar request fails

- **WHEN** one of the four calendar responses fails or lacks myTable while the other three can be parsed
- **THEN** the agent processes the three valid responses, reports the failed market and month, and does not report reconciliation as complete
