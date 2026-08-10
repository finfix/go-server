-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS coin.audit_log
(
    id             SERIAL PRIMARY KEY,                    -- Числовой идентификатор записи
    entity         VARCHAR(100) NOT NULL,                 -- Название сущности (transaction, account и т.д.)
    method         VARCHAR(20)  NOT NULL,                 -- Метод изменения (create, update, delete)
    entity_id      VARCHAR(100) NOT NULL,                 -- Идентификатор сущности (без FK)
    snapshot_before JSONB,                                -- Слепок сущности до изменения
    snapshot_after  JSONB,                                -- Слепок сущности после изменения
    user_id        UUID         NOT NULL,                 -- Пользователь, совершивший действие
    device_id      VARCHAR(255) NOT NULL,                 -- Устройство, с которого совершено действие
    datetime_create TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() -- Дата и время изменения
);
COMMENT ON TABLE coin.audit_log IS 'Аудит изменяющих действий пользователей';
COMMENT ON COLUMN coin.audit_log.entity IS 'Название сущности (transaction, account и т.д.)';
COMMENT ON COLUMN coin.audit_log.method IS 'Метод изменения (create, update, delete)';
COMMENT ON COLUMN coin.audit_log.entity_id IS 'Идентификатор сущности (без FK)';
COMMENT ON COLUMN coin.audit_log.snapshot_before IS 'Слепок сущности до изменения';
COMMENT ON COLUMN coin.audit_log.snapshot_after IS 'Слепок сущности после изменения';
COMMENT ON COLUMN coin.audit_log.user_id IS 'Пользователь, совершивший действие';
COMMENT ON COLUMN coin.audit_log.device_id IS 'Устройство, с которого совершено действие';
COMMENT ON COLUMN coin.audit_log.datetime_create IS 'Дата и время изменения';

CREATE INDEX IF NOT EXISTS idx_audit_log_entity_entity_id ON coin.audit_log (entity, entity_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS coin.audit_log;
-- +goose StatementEnd
