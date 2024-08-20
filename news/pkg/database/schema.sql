CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    title TEXT UNIQUE,
    description TEXT,
    pub_date INTEGER DEFAULT 0,
    source TEXT
);