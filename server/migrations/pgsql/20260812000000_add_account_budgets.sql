-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS coin.account_budgets
(
    id                 UUID PRIMARY KEY,                                    -- Идентификатор версии бюджета
    account_id         UUID                     NOT NULL,                   -- Идентификатор счета
    amount             NUMERIC                  NOT NULL DEFAULT 0,         -- Сумма бюджета
    fixed_sum          NUMERIC                  NOT NULL DEFAULT 0,         -- Фиксированная сумма
    days_offset        INT8                     NOT NULL DEFAULT 0,         -- Смещение в днях
    gradual_filling    BOOL                     NOT NULL DEFAULT TRUE,      -- Заполняется ли бюджет постепенно
    effective_from     DATE                     NOT NULL,                   -- Дата, с которой действует эта версия бюджета
    is_deleted         BOOL                     NOT NULL DEFAULT FALSE,     -- Признак мягкого удаления версии
    created_by_user_id UUID                     NOT NULL,                   -- Идентификатор пользователя, создавшего версию
    datetime_create    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),     -- Дата и время создания версии
    CONSTRAINT fk_account_budgets_accounts FOREIGN KEY (account_id) REFERENCES coin.accounts (id),
    CONSTRAINT fk_account_budgets_users FOREIGN KEY (created_by_user_id) REFERENCES coin.users (id)
);
COMMENT ON TABLE coin.account_budgets IS 'История версий бюджета счета (каждое изменение цели - новая версия, действующая с effective_from)';
COMMENT ON COLUMN coin.account_budgets.account_id IS 'Идентификатор счета';
COMMENT ON COLUMN coin.account_budgets.amount IS 'Сумма бюджета';
COMMENT ON COLUMN coin.account_budgets.fixed_sum IS 'Фиксированная сумма';
COMMENT ON COLUMN coin.account_budgets.days_offset IS 'Смещение в днях';
COMMENT ON COLUMN coin.account_budgets.gradual_filling IS 'Заполняется ли бюджет постепенно';
COMMENT ON COLUMN coin.account_budgets.effective_from IS 'Дата, с которой действует эта версия бюджета';
COMMENT ON COLUMN coin.account_budgets.is_deleted IS 'Признак мягкого удаления версии';
COMMENT ON COLUMN coin.account_budgets.created_by_user_id IS 'Идентификатор пользователя, создавшего версию';
COMMENT ON COLUMN coin.account_budgets.datetime_create IS 'Дата и время создания версии';

CREATE INDEX IF NOT EXISTS idx_account_budgets_account_id_effective_from ON coin.account_budgets (account_id, effective_from);

-- Переносим текущие значения бюджета в первую версию истории для каждого счета
INSERT INTO coin.account_budgets (id, account_id, amount, fixed_sum, days_offset, gradual_filling, effective_from, created_by_user_id, datetime_create)
SELECT gen_random_uuid(), id, COALESCE(budget_amount, 0), budget_fixed_sum, budget_days_offset, budget_gradual_filling, datetime_create::date, created_by_user_id, datetime_create
FROM coin.accounts;

ALTER TABLE coin.accounts DROP COLUMN IF EXISTS budget_amount;
ALTER TABLE coin.accounts DROP COLUMN IF EXISTS budget_fixed_sum;
ALTER TABLE coin.accounts DROP COLUMN IF EXISTS budget_days_offset;
ALTER TABLE coin.accounts DROP COLUMN IF EXISTS budget_gradual_filling;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coin.accounts ADD COLUMN IF NOT EXISTS budget_amount NUMERIC DEFAULT 0;
ALTER TABLE coin.accounts ADD COLUMN IF NOT EXISTS budget_fixed_sum NUMERIC DEFAULT 0 NOT NULL;
ALTER TABLE coin.accounts ADD COLUMN IF NOT EXISTS budget_days_offset INT8 DEFAULT 0 NOT NULL;
ALTER TABLE coin.accounts ADD COLUMN IF NOT EXISTS budget_gradual_filling BOOL DEFAULT TRUE NOT NULL;

UPDATE coin.accounts a
SET budget_amount = latest.amount,
    budget_fixed_sum = latest.fixed_sum,
    budget_days_offset = latest.days_offset,
    budget_gradual_filling = latest.gradual_filling
FROM (
    SELECT DISTINCT ON (account_id) account_id, amount, fixed_sum, days_offset, gradual_filling
    FROM coin.account_budgets
    WHERE is_deleted = false
    ORDER BY account_id, effective_from DESC
) latest
WHERE a.id = latest.account_id;

DROP TABLE IF EXISTS coin.account_budgets;
-- +goose StatementEnd
