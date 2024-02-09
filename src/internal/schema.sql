CREATE TABLE contacts (
  id INTEGER,
  last_name TEXT,
  first_name TEXT,
  phone TEXT,
  email TEXT UNIQUE,
  PRIMARY KEY(id)
);

