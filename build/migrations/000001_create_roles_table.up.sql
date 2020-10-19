CREATE TABLE IF NOT EXISTS roles(
    id INTEGER PRIMARY KEY,
    name VARCHAR(63) NOT NULL,
    description VARCHAR(255) NOT NULL
);

INSERT INTO roles(id, name, description) VALUES(1, 'Арендодатель', 'Размещает квартиру на сайте с целью ее дальнейшей сдачи.');
INSERT INTO roles(id, name, description) VALUES(2, 'Пользователь', 'Ищет на сайте квартиры, в которые можно поселиться.');