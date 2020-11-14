CREATE TABLE IF NOT EXISTS stoves (
    id serial PRIMARY KEY,
    name VARCHAR(64),
    description TEXT
);

INSERT INTO stoves(name, description)
VALUES ('Газовая', ''),
        ('Индукционная', ''),
        ('Электрическая', '');

ALTER TABLE flats ADD FOREIGN KEY(stove) REFERENCES stoves(id);