-- +goose up
-- +goose StatementBegin
CREATE TABLE posts(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NUll,
    FOREIGN KEY (feed_id) REFERENCES feed(id)
);
-- +goose StatementEnd

-- +goose down
DROP TABLE posts;