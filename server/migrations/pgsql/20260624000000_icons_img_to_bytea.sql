-- +goose Up
-- +goose StatementBegin
ALTER TABLE coin.icons ALTER COLUMN img DROP DEFAULT;
ALTER TABLE coin.icons ALTER COLUMN img TYPE BYTEA USING NULL;
COMMENT ON COLUMN coin.icons.img IS 'Содержимое изображения';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coin.icons ALTER COLUMN img TYPE VARCHAR(100) USING NULL;
COMMENT ON COLUMN coin.icons.img IS 'Ссылка на изображение';
-- +goose StatementEnd