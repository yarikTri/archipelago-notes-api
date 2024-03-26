CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Notes (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       VARCHAR(64) DEFAULT '',
    plain_text  TEXT DEFAULT ''
);
