CREATE TABLE Users (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(32),
    name        VARCHAR(100),
    avatar_path TEXT
)

CREATE TABLE Routes (
    id         SERIAL PRIMARY KEY,
    name       TEXT   NOT NULL,
    state      TEXT   DEFAULT 'VALID' CHECK (state IN ('VALID', 'DELETED')),
    image_path TEXT
);

CREATE TABLE Tickets (
    id         SERIAL      PRIMARY KEY,
    user_id    INT         REFERENCES Users(id)  ON DELETE SET NULL,
    started_at TIMESTAMPTZ DEFAULT NOW()                            NOT NULL,
    state TEXT CHECK (state IN ('DRAFT', 'DELETED', 'CREATED', 'DONE', 'REJECTED'))
);

CREATE TABLE Routes_Tickets (
    route_id  REFERENCES Routes(id)  ON DELETE SET NULL,
    ticket_id REFERENCES Tickets(id) ON DELETE SET NULL,

    PRIMARY KEY (route_id, ticket_id)
)

-- CREATE TABLE Stations (
--     id   SERIAL PRIMARY KEY,
--     name TEXT   NOT NULL
-- );

-- CREATE TABLE Routes_Stations (
--     route_id   INT REFERENCES Routes(id)   ON DELETE CASCADE NOT NULL,
--     station_id INT REFERENCES Stations(id) ON DELETE CASCADE NOT NULL,
--     seq_number INT                                           NOT NULL,

--     PRIMARY KEY(route_id, station_id)
-- );
