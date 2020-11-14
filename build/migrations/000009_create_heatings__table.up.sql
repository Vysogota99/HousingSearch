CREATE TABLE IF NOT EXISTS heatings (
    id serial PRIMARY KEY,
    name VARCHAR(64),
    description TEXT
);

INSERT INTO heatings(name, description)
VALUES ('Центральное', ''),
        ('Отдельное', '');


ALTER TABLE flats ADD FOREIGN KEY(heating) REFERENCES heatings(id);
