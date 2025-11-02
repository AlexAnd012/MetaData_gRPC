CREATE TABLE IF NOT EXISTS metadata (
id uuid    PRIMARY KEY,
filename text    NOT NULL,
size_bytes bigint  NOT NULL DEFAULT 0,
content_type text,
owner_id text,
created_at bigint  NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_metadata_created_at
    ON metadata (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_metadata_owner_created
    ON metadata (owner_id, created_at DESC);
