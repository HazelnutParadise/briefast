-- 每日晨報的已讀紀錄。資料可以重新累積，超過保留期由程式直接刪除，
-- 不走軟刪除：文章保留 90 天、事件保留 30 天。
CREATE TABLE seen_articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_id TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL DEFAULT '',
    url_hash TEXT NOT NULL DEFAULT '',
    title TEXT NOT NULL DEFAULT '',
    title_hash TEXT NOT NULL DEFAULT '',
    content_hash TEXT NOT NULL DEFAULT '',
    lead_hash TEXT NOT NULL DEFAULT '',
    sketch TEXT NOT NULL DEFAULT '[]',
    source_domain TEXT NOT NULL DEFAULT '',
    symbols TEXT NOT NULL DEFAULT '[]',
    event_key TEXT NOT NULL DEFAULT '',
    published_at TEXT NOT NULL DEFAULT '',
    decision TEXT NOT NULL,
    stock_call TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX seen_articles_source_id_idx ON seen_articles(source_id) WHERE source_id <> '';
CREATE INDEX seen_articles_url_hash_idx ON seen_articles(url_hash) WHERE url_hash <> '';
CREATE INDEX seen_articles_title_hash_idx ON seen_articles(title_hash) WHERE title_hash <> '';
CREATE INDEX seen_articles_created_at_idx ON seen_articles(created_at);

CREATE TABLE seen_events (
    event_key TEXT PRIMARY KEY,
    note TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX seen_events_created_at_idx ON seen_events(created_at);
