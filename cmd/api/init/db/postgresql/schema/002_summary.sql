CREATE TABLE summ (
	id UUID PRIMARY KEY,
	text_with_role TEXT DEFAULT '',
	role VARCHAR(30) DEFAULT '',
	text TEXT,
	active BOOLEAN,
	platform VARCHAR(30),
	started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	detalization INT
);

CREATE TABLE IF NOT EXISTS summ_to_note (
    summ_id UUID,
	note_id UUID,
	UNIQUE (summ_id, note_id)
);