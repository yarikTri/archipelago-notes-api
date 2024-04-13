CREATE TABLE IF NOT EXISTS `user` (
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    login           VARCHAR(64)     UNIQUE NOT NULL,
    password_hash   VARCHAR(128)    NOT NULL,
    name            VARCHAR(64)     DEFAULT 'Name Surname' NOT NULL,
);

CREATE TABLE IF NOT EXISTS user_root_dir (
    user_id     UUID REFERENCES "user" (id) ON DELETE CASCADE UNIQUE NOT NULL,
    root_dir_id INT REFERENCES dir (id) ON DELETE CASCADE UNIQUE NOT NULL
);
