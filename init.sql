CREATE TABLE IF NOT EXISTS your_table (
  column1 TEXT,
  column2 TEXT
);


CREATE TABLE IF NOT EXISTS temperature (
  id SERIAL PRIMARY KEY,
  location VARCHAR(256),
  timestamp TIMESTAMPTZ DEFAULT now(),
  temperature DOUBLE PRECISION NOT NULL,
  unit TEXT CHECK (unit IN ('C', 'F')) NOT NULL
);


CREATE TABLE IF NOT EXISTS humidity (
  id SERIAL PRIMARY KEY,
  location VARCHAR(256),
  timestamp TIMESTAMPTZ DEFAULT now(),
  humidity DOUBLE PRECISION NOT NULL,
  unit TEXT CHECK (unit IN ('C', 'F')) NOT NULL
);
