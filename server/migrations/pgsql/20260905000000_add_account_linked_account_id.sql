-- +goose Up
-- +goose StatementBegin
ALTER TABLE coin.accounts ADD COLUMN linked_account_id UUID;
COMMENT ON COLUMN coin.accounts.linked_account_id IS 'Идентификатор счёта-моста, с которым связан этот счёт (1:1)';

ALTER TABLE coin.accounts ADD CONSTRAINT fk_accounts_linked_account_id FOREIGN KEY (linked_account_id) REFERENCES coin.accounts (id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_linked_account_id ON coin.accounts (linked_account_id) WHERE linked_account_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS coin.idx_accounts_linked_account_id;
ALTER TABLE coin.accounts DROP CONSTRAINT IF EXISTS fk_accounts_linked_account_id;
ALTER TABLE coin.accounts DROP COLUMN IF EXISTS linked_account_id;
-- +goose StatementEnd
