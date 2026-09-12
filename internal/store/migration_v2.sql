CREATE TABLE conferences (
    symbol TEXT NOT NULL,
    held_on TEXT NOT NULL,
    name TEXT NOT NULL,
    market TEXT NOT NULL,
    held_at TEXT NOT NULL DEFAULT '',
    venue TEXT NOT NULL DEFAULT '',
    announced_on TEXT NOT NULL,
    prediction TEXT NOT NULL,
    payload BLOB NOT NULL,
    pre_close_date TEXT,
    pre_close REAL,
    post_close_date TEXT,
    post_close REAL,
    outcome TEXT,
    settled_at TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    PRIMARY KEY (symbol, held_on)
);

CREATE INDEX conferences_held_on_idx ON conferences(held_on);
CREATE INDEX conferences_settled_idx ON conferences(settled_at, held_on);
