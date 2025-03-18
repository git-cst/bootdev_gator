-- +goose up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose down
DROP TABLE users;
DROP EXTENSION IF EXISTS "uuid-ossp";