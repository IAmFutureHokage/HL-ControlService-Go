package migrations

const CreateTableControlValue = `
CREATE TABLE IF NOT EXISTS control_values (
    id TEXT PRIMARY KEY,
    post_code TEXT NOT NULL,
    type INTEGER NOT NULL,
    date_start TIMESTAMP NOT NULL,
    value INTEGER NOT NULL
);
`
