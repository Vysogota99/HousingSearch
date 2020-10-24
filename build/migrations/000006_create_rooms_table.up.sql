CREATE TABLE IF NOT EXISTS rooms (
    id serial PRIMARY KEY,
    flatID INT REFERENCES flats(id) ON DELETE CASCADE NOT NULL,
    maxResidents INT NOT NULL,
    description TEXT NOT NULL,
    currNumberOfResidents INT NOT NULL,
    numOfWindows INT NOT NULL,
    balcony boolean NOT NULL,
    numOfTables INT NOT NULL,
    numOfChairs INT NOT NULL,
    TV boolean NOT NULL,
    numOfCupboards INT NOT NULL,
    area INT NOT NULL
);