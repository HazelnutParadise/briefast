# news-dedup-api Specification

## Purpose
在伺服器保存每日晨報處理過的新聞與跨天事件，避免重複報導。
## Requirements
### Requirement: Dedup endpoints require a Bearer key

Every /api/seen endpoint SHALL require Bearer key authentication and SHALL re-check key status on every request. A missing, invalid, or revoked key SHALL receive 401 and SHALL be recorded in the update log with a seen-rejected-auth action. A request body that is not valid JSON of the documented shape SHALL receive 400 with an errors array and SHALL be recorded in the update log with a seen-rejected-validation action. Successful dedup calls SHALL NOT be written to the update log.

#### Scenario: Revoked key on peek

- **WHEN** a request with a revoked key calls POST /api/seen/peek
- **THEN** the API responds 401 and records a seen-rejected-auth entry

---

### Requirement: Peek returns unread candidates

POST /api/seen/peek SHALL accept `{"items":[…]}` where each item carries at least one of `url`, `title`, or `source_id`. It SHALL identify an item as already seen when its source ID, canonical URL hash, or normalized title hash matches a stored article, checked in that order. The source ID SHALL be the explicit `source_id`, or else the URL host without `www.` joined by a colon to the last numeric path segment, or else the last path segment. The canonical URL SHALL drop the fragment, a trailing slash, and the query keys fbclid, gclid, igshid, ref, share, from, src, spm and any key starting with utm_, pk_, or mtm_. The response SHALL carry ok=true, `new`, and `seen`, where each seen item is the original item plus `seen_by` set to `source_id`, `url`, or `title`.

#### Scenario: Tracking parameters do not hide a seen article

- **WHEN** an article was recorded with URL `https://news.cnyes.com/news/id/6001` and peek receives `https://news.cnyes.com/news/id/6001/?utm_source=line`
- **THEN** the item is returned under seen with `seen_by` source_id or url

#### Scenario: Fresh candidate

- **WHEN** peek receives an item whose source ID, URL, and title match no stored article
- **THEN** the item is returned under new unchanged

---

### Requirement: Record stores processed articles

POST /api/seen/record SHALL accept `{"items":[…]}` where each item may carry url, title, source_id, content or summary, symbols, event_key, published_at, decision, and call. It SHALL skip items whose source ID or canonical URL hash is already stored and SHALL store the rest with the content hash, the hash of the first 100 normalized characters, and a similarity sketch when content is present. A missing decision SHALL be stored as `skipped`. The response SHALL carry ok=true, `added` with the number of stored items, and `duplicates` with the number skipped.

#### Scenario: Skipped article recorded

- **WHEN** record receives an article with decision skipped
- **THEN** a later peek for the same URL returns it under seen

#### Scenario: Duplicate record ignored

- **WHEN** record receives an article whose URL is already stored
- **THEN** the article is not stored twice and the response counts it as skipped

---

### Requirement: Similar finds rewritten copies

POST /api/seen/similar SHALL accept `{"content":"…","threshold":0.25}` with threshold optional and defaulting to 0.25. It SHALL normalize the content by removing HTML tags and bylines, applying NFKC, and collapsing whitespace, build character bigrams without punctuation, and estimate Jaccard similarity against each stored article sketch using bottom-64 MinHash. It SHALL return ok=true, the applied `threshold`, and `matches` with each match's similarity rounded to three decimals, title, url, and event_key, sorted by similarity descending.

#### Scenario: Reprinted article matched

- **WHEN** a stored article and the submitted content are the same body text with a different byline
- **THEN** the stored article is returned with similarity at or above 0.9

#### Scenario: Unrelated article not matched

- **WHEN** the submitted content shares no topic with any stored article
- **THEN** matches is empty at the default threshold

---

### Requirement: Event registry spans days

POST /api/seen/events SHALL register an event key of the form `symbol|event type|period` with an optional note, and SHALL leave an existing key unchanged, answering ok=true with `added` false for an existing key. A key that does not have three non-empty parts separated by `|` SHALL receive 400. POST /api/seen/events/check SHALL return ok=true, the `key`, `reported` true or false, and `record` holding the stored event or null. A registered event SHALL NOT restrict whether the stock appears in calls.

#### Scenario: Continuing event detected

- **WHEN** `2330|月營收|2026-07` was registered yesterday and check is called with the same key today
- **THEN** the response has reported true and the stored record

---

### Requirement: Dedup storage and expiry

Dedup records SHALL be stored in SQLite tables created by a new migration and SHALL survive restarts. Before handling each dedup request, the server SHALL delete articles first seen more than 90 days ago and events first seen more than 30 days ago.

#### Scenario: Expired article forgotten

- **WHEN** an article was first seen 91 days ago and peek receives its URL
- **THEN** the item is returned under new and the stored row has been deleted

