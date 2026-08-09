CREATE TABLE reports (
    date TEXT PRIMARY KEY,
    headline TEXT NOT NULL,
    payload BLOB NOT NULL,
    generated_at TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE api_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL,
    revoked_at TEXT
);

CREATE TABLE update_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    at TEXT NOT NULL,
    api_key_id INTEGER,
    key_name TEXT NOT NULL,
    report_date TEXT NOT NULL DEFAULT '',
    action TEXT NOT NULL,
    detail TEXT NOT NULL DEFAULT '',
    FOREIGN KEY (api_key_id) REFERENCES api_keys(id) ON DELETE RESTRICT
);

CREATE INDEX update_log_api_key_id_idx ON update_log(api_key_id);
CREATE INDEX update_log_at_idx ON update_log(at DESC, id DESC);
