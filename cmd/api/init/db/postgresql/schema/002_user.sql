CREATE TABLE IF NOT EXISTS user (
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    login           VARCHAR(64)     UNIQUE NOT NULL,
    password_hash   VARCHAR(128)    NOT NULL,
    name            VARCHAR(64)     DEFAULT 'Name Surname' NOT NULL
);
