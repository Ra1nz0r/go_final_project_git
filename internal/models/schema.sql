CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    title TEXT NOT NULL,
    comment TEXT NOT NULL,
    repeat TEXT NOT NULL,
    search TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);