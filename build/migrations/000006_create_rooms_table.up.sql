CREATE TABLE IF NOT EXISTS rooms (
    id serial PRIMARY KEY,
    flat_id INT REFERENCES flats(id) ON DELETE CASCADE NOT NULL,
    max_residents INT NOT NULL,
    description TEXT DEFAULT '',
    price NUMERIC DEFAULT 0,
    deposit NUMERIC DEFAULT 0,
    curr_number_of_residents INT NOT NULL,
    balcony boolean NOT NULL,
    num_of_tables INT NOT NULL,
    num_of_chairs INT NOT NULL,
    tv boolean NOT NULL,
    furniture boolean NOT NULL,
    area INT NOT NULL,
    windows boolean NOT NULL,
    is_visible boolean DEFAULT TRUE
);