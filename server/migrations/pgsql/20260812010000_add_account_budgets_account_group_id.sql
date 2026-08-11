-- +goose Up
-- +goose StatementBegin
ALTER TABLE coin.account_budgets ADD COLUMN account_group_id UUID;

UPDATE coin.account_budgets ab
SET account_group_id = a.account_group_id
FROM coin.accounts a
WHERE a.id = ab.account_id;

ALTER TABLE coin.account_budgets ALTER COLUMN account_group_id SET NOT NULL;
ALTER TABLE coin.account_budgets ADD CONSTRAINT fk_account_budgets_account_groups FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id);

COMMENT ON COLUMN coin.account_budgets.account_group_id IS 'Идентификатор группы счетов (денормализовано со счета для быстрой фильтрации и аудита)';

CREATE INDEX IF NOT EXISTS idx_account_budgets_account_group_id ON coin.account_budgets (account_group_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coin.account_budgets DROP COLUMN IF EXISTS account_group_id;
-- +goose StatementEnd
