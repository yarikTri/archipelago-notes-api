CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Notes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    automerge_url   VARCHAR(128) NOT NULL,
    title           VARCHAR(64) DEFAULT ''
);
