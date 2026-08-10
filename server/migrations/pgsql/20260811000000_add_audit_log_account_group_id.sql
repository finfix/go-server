-- +goose Up
-- +goose StatementBegin
ALTER TABLE coin.audit_log
    ADD COLUMN account_group_id UUID; -- Группа счетов, в которой изменена сущность (NULL для глобальных сущностей)
COMMENT ON COLUMN coin.audit_log.account_group_id IS 'Группа счетов, в которой изменена сущность (NULL для глобальных сущностей)';

CREATE INDEX IF NOT EXISTS idx_audit_log_account_group_id ON coin.audit_log (account_group_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS coin.idx_audit_log_account_group_id;
ALTER TABLE coin.audit_log
    DROP COLUMN IF EXISTS account_group_id;
-- +goose StatementEnd
