CREATE EXTENSION IF NOT EXISTS "ltree";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- TODO:
--  - trigger on insert on dir - set default_access as parent's one
--  - trigger on update of default_access on dir - set same access to children

-- Users table
CREATE TABLE IF NOT EXISTS "user" (
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    email           VARCHAR(128)    UNIQUE NOT NULL,
    email_confirmed BOOLEAN         DEFAULT false NOT NULL,
    password_hash   VARCHAR(128)    NOT NULL,
    name            VARCHAR(64)     DEFAULT 'Name Surname' NOT NULL
);

-- Dirs table
CREATE TABLE IF NOT EXISTS dir (
    id      SERIAL          CONSTRAINT dir_pk PRIMARY KEY,
    name    VARCHAR(64)     DEFAULT ''  NOT NULL,
    path    LTREE           DEFAULT ''  NOT NULL CONSTRAINT dir_pk_path UNIQUE
    -- creator_id      UUID            REFERENCES `user` (id) NOT NULL,
    -- default_access  CHAR            DEFAULT 'e' NOT NULL,

    -- -- empty, read, write (into notes), modify (CUD notes/dirs) & manage access
    -- CHECK (default_access IN ('e', 'r', 'w', 'm', 'ma'))
);

CREATE INDEX dir_path_idx ON dir USING gist (path);

-- Triggers functions for paths consistency

-- Constraint for existing path
CREATE OR REPLACE FUNCTION dir_before_update_insert_check_path()
    RETURNS TRIGGER AS $dir_before_update_insert_check_path$
DECLARE
parentPath ltree;
    curLabel ltree;
    parentId text;
BEGIN
    -- ! При INSERT | UPDATE обязательно следует передавать родительский path в поле `path`
    -- Конкотенация id в конец `path`
    NEW.path := NEW.path || NEW.id::text;

    parentPath := subpath(NEW.path, 0, -1);
    curLabel := subpath(NEW.path, -1);

    -- последняя метка в пути должна равняться id
    IF (curLabel::text != NEW.id::text) THEN
        RAISE EXCEPTION 'The last path label % must be equal id %', curLabel::text, NEW.id::text;
END IF;

    -- должна существовать родительская запись с подходящим путем, если новая запись не корневая
    IF (parentPath != '') THEN
        parentId := (SELECT id FROM dir WHERE PATH = parentPath);
        if (parentId IS NULL) THEN
            RAISE EXCEPTION 'Parent dir with path % not found', parentPath;
END IF;
END IF;

RETURN NEW;
END;
$dir_before_update_insert_check_path$ LANGUAGE plpgsql;

CREATE TRIGGER tr_dir_before_update_insert_check_path
    BEFORE INSERT OR UPDATE OF path ON dir
    FOR EACH ROW
    EXECUTE FUNCTION dir_before_update_insert_check_path();

-- Update paths for all children
CREATE OR REPLACE FUNCTION dir_after_update_set_children_path()
    RETURNS TRIGGER AS $dir_after_update_set_children_path$
BEGIN
    IF (NEW.path != OLD.path) THEN
UPDATE dir
SET path = NEW.path || subpath(path, nlevel(OLD.path))
WHERE path <@ OLD.path;
END IF;

RETURN NULL;
END;
$dir_after_update_set_children_path$ LANGUAGE plpgsql;

CREATE TRIGGER tr_update_children_dir_path
    AFTER UPDATE OF path ON dir
    FOR EACH ROW
    EXECUTE FUNCTION dir_after_update_set_children_path();

-- Users' root dirs table
CREATE TABLE IF NOT EXISTS user_root_dir (
    user_id     UUID REFERENCES "user" (id) ON DELETE CASCADE UNIQUE NOT NULL,
    root_dir_id INT REFERENCES dir (id) ON DELETE CASCADE UNIQUE NOT NULL
);

-- Notes table
CREATE TABLE IF NOT EXISTS note (
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    automerge_url   VARCHAR(128)    NOT NULL,
    title           VARCHAR(64)     DEFAULT 'Untitled' NOT NULL,
    dir_id          INT             REFERENCES dir (id) ON DELETE CASCADE NOT NULL,
    creator_id      UUID            REFERENCES "user" (id) NOT NULL,
    default_access  VARCHAR(2)      DEFAULT 'e' NOT NULL,

    -- empty, read, write, modify, manage access
    CHECK (default_access IN ('e', 'r', 'w'))
);

-- Users' accesses to notes table
CREATE TABLE IF NOT EXISTS note_access (
    note_id UUID        REFERENCES note (id) NOT NULL,
    user_id UUID        REFERENCES "user" (id) NOT NULL,
    access  VARCHAR(2)  DEFAULT 'r' NOT NULL,

    UNIQUE(note_id, user_id),
    -- empty (for black list), read, write, manage access
    CHECK (access IN ('e', 'r', 'w', 'm', 'ma'))
);
