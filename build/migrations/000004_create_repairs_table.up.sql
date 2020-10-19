CREATE TABLE IF NOT EXISTS repairs (
    id serial PRIMARY KEY,
    name VARCHAR(64),
    description TEXT
);

INSERT INTO repairs(name, description)
VALUES ('Бабушкин', 'Ремонт времен СССР, с коричневыми шкафами коврами на стене'),
        ('Евроремнт', 'Ремонт, сделанный без интузиазма своими руками пару лет назад.'),
        ('Дизайнерский', 'Ремонт, сделанный недавно и со вкусом.')