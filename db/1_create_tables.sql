CREATE TABLE geo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    latitude REAL,
    longitude REAL,
    access_date DATETIME
);

CREATE TABLE weather (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    geo_id INTEGER NOT NULL,
    current REAL,
    daily_dates TEXT,
    daily_min TEXT,
    daily_max TEXT,
    latitude REAL,
    longitude REAL,
    FOREIGN KEY (geo_id) REFERENCES geo(id)
);