CREATE TABLE IF NOT EXISTS living_places(
    id serial PRIMARY KEY,
    roomID INT REFERENCES rooms(id) ON DELETE CASCADE,
    residentID INT REFERENCES users(id),
    price NUMERIC NOT NULL,
    description TEXT NOT NULL,
    numOFBerths INT NOT NULL,
    deposit NUMERIC NOT NULL
);

CREATE INDEX living_places_price ON living_places(price);
CREATE INDEX living_places_deposit ON living_places(deposit);

