CREATE TABLE Routes (
    id   SERIAL PRIMARY KEY,
    name TEXT   NOT NULL
);

CREATE TABLE Stations (
    id   SERIAL PRIMARY KEY,
    name TEXT   NOT NULL
);

CREATE TABLE Routes_Stations (
    route_id   INT REFERENCES Routes(id)   ON DELETE CASCADE NOT NULL,
    station_id INT REFERENCES Stations(id) ON DELETE CASCADE NOT NULL,
    seq_number INT                                           NOT NULL,

    PRIMARY KEY(route_id, station_id)
);

CREATE TABLE Tickets (
    id       SERIAL    PRIMARY KEY,
    route_id INT       REFERENCES Routes(id) ON DELETE CASCADE NOT NULL,
    -- TODO: user_id INT REFERENCES users(id) ON DELETE SET NULL,
    started_at TIMESTAMPTZ DEFAULT NOW()                       NOT NULL
);
