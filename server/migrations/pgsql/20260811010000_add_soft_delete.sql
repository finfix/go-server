-- +goose Up
-- +goose StatementBegin
ALTER TABLE coin.accounts ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE coin.account_groups ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE coin.tags ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE coin.transactions ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT false;

COMMENT ON COLUMN coin.accounts.is_deleted IS 'Признак мягкого удаления счета';
COMMENT ON COLUMN coin.account_groups.is_deleted IS 'Признак мягкого удаления группы счетов';
COMMENT ON COLUMN coin.tags.is_deleted IS 'Признак мягкого удаления подкатегории';
COMMENT ON COLUMN coin.transactions.is_deleted IS 'Признак мягкого удаления транзакции';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coin.accounts DROP COLUMN IF EXISTS is_deleted;
ALTER TABLE coin.account_groups DROP COLUMN IF EXISTS is_deleted;
ALTER TABLE coin.tags DROP COLUMN IF EXISTS is_deleted;
ALTER TABLE coin.transactions DROP COLUMN IF EXISTS is_deleted;
-- +goose StatementEnd
