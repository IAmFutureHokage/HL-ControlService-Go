package migrations

const CreateTableControlValue = `
CREATE TABLE IF NOT EXISTS control_values (
    id UUID PRIMARY KEY,
    post_code TEXT NOT NULL,
    type INTEGER NOT NULL,
    date_start TIMESTAMP NOT NULL,
    value INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS waterlevels (
    id UUID PRIMARY KEY,
    post_code TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    waterlevel INTEGER NOT NULL
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'control_values_unique_constraint'
    ) THEN
        ALTER TABLE control_values
        ADD CONSTRAINT control_values_unique_constraint UNIQUE (post_code, type, date_start);
    END IF;
END
$$;

CREATE INDEX IF NOT EXISTS idx_control_values_post_code_type_date ON control_values(post_code, type, date_start);
CREATE INDEX IF NOT EXISTS idx_control_values_date_start ON control_values(date_start);
`
