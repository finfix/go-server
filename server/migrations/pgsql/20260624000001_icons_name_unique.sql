-- +goose Up
-- +goose StatementBegin
ALTER TABLE coin.icons ADD CONSTRAINT icons_name_unique UNIQUE (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coin.icons DROP CONSTRAINT icons_name_unique;
-- +goose StatementEnd