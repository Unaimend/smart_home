CREATE TABLE IF NOT EXISTS your_table (
  column1 TEXT,
  column2 TEXT
);


CREATE TABLE temperature_data (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMPTZ DEFAULT now(),
  temperature DOUBLE PRECISION NOT NULL,
  unit TEXT CHECK (unit IN ('C', 'F')) NOT NULL
);


CREATE TABLE temperature_data (
  id SERIAL PRIMARY KEY,
  timestamp TIMESTAMPTZ DEFAULT now(),
  humidity DOUBLE PRECISION NOT NULL,
  unit TEXT CHECK (unit IN ('C', 'F')) NOT NULL
);
