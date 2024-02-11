CREATE TABLE contacts (
  id INTEGER,
  last_name TEXT,
  first_name TEXT,
  phone TEXT,
  email TEXT UNIQUE,
  PRIMARY KEY(id)
);

--
-- insert some fixtures
--
INSERT INTO contacts (last_name, first_name, phone, email)
    VALUES
    ('Abrahams', 'Alice', '1-234-567711', 'alice.abrahams@manundone.co.uk'),
    ('Booqie', 'Bob', '1-917-4890931', 'bob.bookie@bookiebob.com')
    ;
