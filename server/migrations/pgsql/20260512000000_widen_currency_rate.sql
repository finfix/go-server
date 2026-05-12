-- +goose Up
-- +goose StatementBegin
ALTER TABLE coin.currencies ALTER COLUMN rate TYPE NUMERIC(30, 8);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coin.currencies ALTER COLUMN rate TYPE NUMERIC(15, 8);
-- +goose StatementEnd
