CREATE TABLE IF NOT EXISTS living_places(
    id serial PRIMARY KEY,
    roomID INT REFERENCES rooms(id) ON DELETE CASCADE,
    residentID INT REFERENCES users(id),
    price NUMERIC DEFAULT 0,
    description TEXT DEFAULT '',
    numOFBerths INT NOT NULL,
    deposit NUMERIC DEFAULT 0
);

CREATE INDEX living_places_price ON living_places(price);
CREATE INDEX living_places_deposit ON living_places(deposit);

