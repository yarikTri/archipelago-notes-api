CREATE TABLE tag (
    tag_id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tag_to_note (
    tag_id UUID,
    note_id UUID,
    UNIQUE (tag_id, note_id)
);

CREATE TABLE IF NOT EXISTS tag_to_tag (
    tag_1_id UUID,
    tag_2_id UUID,
    UNIQUE (tag_1_id, tag_2_id)
);
