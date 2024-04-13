CREATE EXTENSION IF NOT EXISTS "ltree";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Dirs table
CREATE TABLE IF NOT EXISTS dir (
    id      SERIAL          CONSTRAINT dir_pk PRIMARY KEY,
    name    VARCHAR(64)     DEFAULT '' NOT NULL,
    path    ltree           DEFAULT '' NOT NULL CONSTRAINT dir_pk_path UNIQUE
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

-- Notes table
CREATE TABLE IF NOT EXISTS note (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    automerge_url   VARCHAR(128) NOT NULL,
    title           VARCHAR(64) DEFAULT '',
    dir_id          INT REFERENCES dir (id) ON DELETE CASCADE NOT NULL
);
