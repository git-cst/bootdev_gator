-- +goose up
-- +goose StatementBegin
ALTER TABLE feed
ADD COLUMN last_fetched_at TIMESTAMP DEFAULT NULL;
-- +goose StatementEnd

-- +goose down
ALTER TABLE feed
DROP COLUMN last_fetched_at;