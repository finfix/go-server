-- +goose Up
-- +goose StatementBegin
ALTER TABLE coin.devices ADD COLUMN last_synced_audit_log_id INT8 NOT NULL DEFAULT 0;
COMMENT ON COLUMN coin.devices.last_synced_audit_log_id IS 'Идентификатор последней записи аудит-лога, которую устройство подтвердило (Confirm) как синхронизированную';

ALTER TABLE coin.devices ADD COLUMN pending_checkpoint INT8;
COMMENT ON COLUMN coin.devices.pending_checkpoint IS 'Чекпоинт, выданный последним вызовом Sync, ожидающий подтверждения от устройства';

ALTER TABLE coin.devices ADD COLUMN pending_sync_token UUID;
COMMENT ON COLUMN coin.devices.pending_sync_token IS 'Опорный токен последнего вызова Sync, подтверждающий, что именно этот ответ был получен и применен устройством';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coin.devices DROP COLUMN IF EXISTS pending_sync_token;
ALTER TABLE coin.devices DROP COLUMN IF EXISTS pending_checkpoint;
ALTER TABLE coin.devices DROP COLUMN IF EXISTS last_synced_audit_log_id;
-- +goose StatementEnd
