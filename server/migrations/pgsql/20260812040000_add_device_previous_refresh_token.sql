-- +goose Up
-- +goose StatementBegin
ALTER TABLE coin.devices ADD COLUMN previous_refresh_token TEXT;
COMMENT ON COLUMN coin.devices.previous_refresh_token IS 'Refresh-токен, действовавший до последней ротации - остаётся валидным до previous_refresh_token_expires_at на случай потери ответа сети при рефреше';

ALTER TABLE coin.devices ADD COLUMN previous_refresh_token_expires_at TIMESTAMPTZ;
COMMENT ON COLUMN coin.devices.previous_refresh_token_expires_at IS 'Момент, до которого previous_refresh_token остаётся валидным';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coin.devices DROP COLUMN IF EXISTS previous_refresh_token_expires_at;
ALTER TABLE coin.devices DROP COLUMN IF EXISTS previous_refresh_token;
-- +goose StatementEnd
