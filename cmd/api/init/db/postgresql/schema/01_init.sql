-- CREATE TABLE Users (
--     id          SERIAL PRIMARY KEY,
--     username    VARCHAR(32),
--     name        VARCHAR(100),
--     avatar_path TEXT
-- )

-- CREATE TABLE Routes (
--     id               SERIAL  PRIMARY KEY,
--     name             TEXT    NOT NULL,
--     active           BOOLEAN NOT NULL,
--     capacity         INT     NOT NULL,
--     end_stations     TEXT[]  NOT NULL,
--     start_time       TEXT    NOT NULL,
--     end_time         TEXT    NOT NULL,
--     interval_minutes INT     NOT NULL,
--     description      TEXT    NOT NULL,
--     image_path       TEXT    DEFAULT '' NOT NULL
-- );

-- CREATE TABLE Tickets (
--     id         SERIAL      PRIMARY KEY,
--     user_id    INT         REFERENCES Users(id)  ON DELETE SET 0,
--     start_date TEXT        NOT NULL,
--     state TEXT CHECK (state IN ('DRAFT', 'DELETED', 'CREATED', 'DONE', 'REJECTED'))
-- );

-- CREATE TABLE Routes_Tickets (
--     route_id  REFERENCES Routes(id)  ON DELETE SET NULL,
--     ticket_id REFERENCES Tickets(id) ON DELETE SET NULL,

--     PRIMARY KEY (route_id, ticket_id)
-- )

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

INSERT INTO Routes (name, active, capacity, interval_minutes, start_time, end_time, start_station, end_station, image_path, description) VALUES 
    ('Автобус №12', true, 50, 10, '05:00', '23:30', 'ПАТП №1', 'Завод "Призма"', '/static/image/12.jpeg', 'Автобус, позволяющий из окна насладиться всеми главными городскими достопримечательностями'),
    ('Троллейбус №6', true, 50, 10, '05:30', '23:00', 'Завод "Призма"', 'Улица Ворошилова', '/static/image/6.jpeg', 'От Мариевки до Ворошилова с божьей помощью'),
    ('Автобус №3', true, 50, 10, '05:00', '23:30', 'Ж/Д Вокзал', 'Переборы', '/static/image/3.jpeg', 'Всегда пустой автобус');
