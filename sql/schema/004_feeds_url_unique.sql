-- +goose Up
ALTER TABLE feeds ADD CONSTRAINT feeds_url_unique UNIQUE (url, user_id);

-- +goose Down
ALTER TABLE feeds DROP CONSTRAINT feeds_url_unique;
