-- +goose Up
-- +goose StatementBegin

-- Включение расширения для работы с UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создание функции для генерации UUID v4 если она не существует
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_proc 
        WHERE proname = 'uuid_generate_v4' 
        AND pg_function_is_visible(oid)
    ) THEN
        CREATE FUNCTION uuid_generate_v4() RETURNS UUID 
        AS 'SELECT gen_random_uuid()' 
        LANGUAGE SQL;
    END IF;
END $$;

-- Удаляем все существующие foreign key constraints
ALTER TABLE coin.accounts DROP CONSTRAINT IF EXISTS accounts_fk;
ALTER TABLE coin.accounts DROP CONSTRAINT IF EXISTS accounts_fk_1;
ALTER TABLE coin.accounts DROP CONSTRAINT IF EXISTS accounts_fk_2;
ALTER TABLE coin.accounts DROP CONSTRAINT IF EXISTS accounts_fk_3;
ALTER TABLE coin.accounts DROP CONSTRAINT IF EXISTS accounts_fk_4;
ALTER TABLE coin.accounts DROP CONSTRAINT IF EXISTS accounts_parent_account_fk;

ALTER TABLE coin.transactions DROP CONSTRAINT IF EXISTS orders_fk;
ALTER TABLE coin.transactions DROP CONSTRAINT IF EXISTS orders_fk_1;
ALTER TABLE coin.transactions DROP CONSTRAINT IF EXISTS transactions_fk;

ALTER TABLE coin.tags DROP CONSTRAINT IF EXISTS tags_fk;

ALTER TABLE coin.tags_to_transaction DROP CONSTRAINT IF EXISTS tags_to_transaction_fk;
ALTER TABLE coin.tags_to_transaction DROP CONSTRAINT IF EXISTS tags_to_transaction_fk_1;

ALTER TABLE coin.users_to_account_groups DROP CONSTRAINT IF EXISTS users_to_account_groups_fk;
ALTER TABLE coin.users_to_account_groups DROP CONSTRAINT IF EXISTS users_to_account_groups_fk_1;

-- Добавляем новые колонки с UUID типом
ALTER TABLE coin.users ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE coin.account_groups ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE coin.accounts ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE coin.accounts ADD COLUMN account_group_id_new UUID;
ALTER TABLE coin.accounts ADD COLUMN parent_account_id_new UUID;
ALTER TABLE coin.accounts ADD COLUMN created_by_user_id_new UUID;
ALTER TABLE coin.transactions ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE coin.transactions ADD COLUMN account_from_id_new UUID;
ALTER TABLE coin.transactions ADD COLUMN account_to_id_new UUID;
ALTER TABLE coin.transactions ADD COLUMN created_by_user_id_new UUID;
ALTER TABLE coin.tags ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE coin.tags ADD COLUMN account_group_id_new UUID;
ALTER TABLE coin.tags_to_transaction ADD COLUMN tag_id_new UUID;
ALTER TABLE coin.tags_to_transaction ADD COLUMN transaction_id_new UUID;
ALTER TABLE coin.users_to_account_groups ADD COLUMN user_id_new UUID;
ALTER TABLE coin.users_to_account_groups ADD COLUMN account_group_id_new UUID;

-- Обновляем references для account_groups
UPDATE coin.accounts SET account_group_id_new = (
    SELECT ag.id_new FROM coin.account_groups ag WHERE ag.id = coin.accounts.account_group_id
);

-- Обновляем references для parent accounts
UPDATE coin.accounts SET parent_account_id_new = (
    SELECT a.id_new FROM coin.accounts a WHERE a.id = coin.accounts.parent_account_id
) WHERE parent_account_id IS NOT NULL;

-- Обновляем references для created_by_user_id в accounts
UPDATE coin.accounts SET created_by_user_id_new = (
    SELECT u.id_new FROM coin.users u WHERE u.id = coin.accounts.created_by_user_id
);

-- Обновляем references для transactions
UPDATE coin.transactions SET account_from_id_new = (
    SELECT a.id_new FROM coin.accounts a WHERE a.id = coin.transactions.account_from_id
);

UPDATE coin.transactions SET account_to_id_new = (
    SELECT a.id_new FROM coin.accounts a WHERE a.id = coin.transactions.account_to_id
);

UPDATE coin.transactions SET created_by_user_id_new = (
    SELECT u.id_new FROM coin.users u WHERE u.id = coin.transactions.created_by_user_id
);

-- Обновляем references для tags
UPDATE coin.tags SET account_group_id_new = (
    SELECT ag.id_new FROM coin.account_groups ag WHERE ag.id = coin.tags.account_group_id
);

-- Обновляем references для tags_to_transaction
UPDATE coin.tags_to_transaction SET tag_id_new = (
    SELECT t.id_new FROM coin.tags t WHERE t.id = coin.tags_to_transaction.tag_id
);

UPDATE coin.tags_to_transaction SET transaction_id_new = (
    SELECT tr.id_new FROM coin.transactions tr WHERE tr.id = coin.tags_to_transaction.transaction_id
);

-- Обновляем references для users_to_account_groups
UPDATE coin.users_to_account_groups SET user_id_new = (
    SELECT u.id_new FROM coin.users u WHERE u.id = coin.users_to_account_groups.user_id
);

UPDATE coin.users_to_account_groups SET account_group_id_new = (
    SELECT ag.id_new FROM coin.account_groups ag WHERE ag.id = coin.users_to_account_groups.account_group_id
);

-- Удаляем старые колонки и переименовываем новые
-- Users
ALTER TABLE coin.users DROP COLUMN id;
ALTER TABLE coin.users RENAME COLUMN id_new TO id;
ALTER TABLE coin.users ADD PRIMARY KEY (id);

-- Account Groups
ALTER TABLE coin.account_groups DROP COLUMN id;
ALTER TABLE coin.account_groups RENAME COLUMN id_new TO id;
ALTER TABLE coin.account_groups ADD PRIMARY KEY (id);

-- Accounts
ALTER TABLE coin.accounts DROP COLUMN id;
ALTER TABLE coin.accounts DROP COLUMN account_group_id;
ALTER TABLE coin.accounts DROP COLUMN parent_account_id;
ALTER TABLE coin.accounts DROP COLUMN created_by_user_id;
ALTER TABLE coin.accounts RENAME COLUMN id_new TO id;
ALTER TABLE coin.accounts RENAME COLUMN account_group_id_new TO account_group_id;
ALTER TABLE coin.accounts RENAME COLUMN parent_account_id_new TO parent_account_id;
ALTER TABLE coin.accounts RENAME COLUMN created_by_user_id_new TO created_by_user_id;
ALTER TABLE coin.accounts ADD PRIMARY KEY (id);
ALTER TABLE coin.accounts ALTER COLUMN account_group_id SET NOT NULL;
ALTER TABLE coin.accounts ALTER COLUMN created_by_user_id SET NOT NULL;

-- Transactions  
ALTER TABLE coin.transactions DROP COLUMN id;
ALTER TABLE coin.transactions DROP COLUMN account_from_id;
ALTER TABLE coin.transactions DROP COLUMN account_to_id;
ALTER TABLE coin.transactions DROP COLUMN created_by_user_id;
ALTER TABLE coin.transactions RENAME COLUMN id_new TO id;
ALTER TABLE coin.transactions RENAME COLUMN account_from_id_new TO account_from_id;
ALTER TABLE coin.transactions RENAME COLUMN account_to_id_new TO account_to_id;
ALTER TABLE coin.transactions RENAME COLUMN created_by_user_id_new TO created_by_user_id;
ALTER TABLE coin.transactions ADD PRIMARY KEY (id);
ALTER TABLE coin.transactions ALTER COLUMN account_from_id SET NOT NULL;
ALTER TABLE coin.transactions ALTER COLUMN account_to_id SET NOT NULL;
ALTER TABLE coin.transactions ALTER COLUMN created_by_user_id SET NOT NULL;

-- Tags
ALTER TABLE coin.tags DROP COLUMN id;
ALTER TABLE coin.tags DROP COLUMN account_group_id;
ALTER TABLE coin.tags RENAME COLUMN id_new TO id;
ALTER TABLE coin.tags RENAME COLUMN account_group_id_new TO account_group_id;
ALTER TABLE coin.tags ADD PRIMARY KEY (id);
ALTER TABLE coin.tags ALTER COLUMN account_group_id SET NOT NULL;

-- Tags to Transaction
ALTER TABLE coin.tags_to_transaction DROP COLUMN tag_id;
ALTER TABLE coin.tags_to_transaction DROP COLUMN transaction_id;
ALTER TABLE coin.tags_to_transaction RENAME COLUMN tag_id_new TO tag_id;
ALTER TABLE coin.tags_to_transaction RENAME COLUMN transaction_id_new TO transaction_id;
ALTER TABLE coin.tags_to_transaction ADD PRIMARY KEY (transaction_id, tag_id);
ALTER TABLE coin.tags_to_transaction ALTER COLUMN tag_id SET NOT NULL;
ALTER TABLE coin.tags_to_transaction ALTER COLUMN transaction_id SET NOT NULL;

-- Users to Account Groups
ALTER TABLE coin.users_to_account_groups DROP COLUMN user_id;
ALTER TABLE coin.users_to_account_groups DROP COLUMN account_group_id;
ALTER TABLE coin.users_to_account_groups RENAME COLUMN user_id_new TO user_id;
ALTER TABLE coin.users_to_account_groups RENAME COLUMN account_group_id_new TO account_group_id;
ALTER TABLE coin.users_to_account_groups ADD PRIMARY KEY (user_id, account_group_id);
ALTER TABLE coin.users_to_account_groups ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE coin.users_to_account_groups ALTER COLUMN account_group_id SET NOT NULL;

-- Восстанавливаем foreign key constraints
ALTER TABLE coin.accounts ADD CONSTRAINT accounts_fk_1 FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id);
ALTER TABLE coin.accounts ADD CONSTRAINT accounts_fk_2 FOREIGN KEY (created_by_user_id) REFERENCES coin.users (id);
ALTER TABLE coin.accounts ADD CONSTRAINT accounts_parent_account_fk FOREIGN KEY (parent_account_id) REFERENCES coin.accounts (id);
ALTER TABLE coin.accounts ADD CONSTRAINT accounts_fk FOREIGN KEY (currency_signatura) REFERENCES coin.currencies (signatura);
ALTER TABLE coin.accounts ADD CONSTRAINT accounts_fk_3 FOREIGN KEY (type_signatura) REFERENCES coin.account_types (signatura);
ALTER TABLE coin.accounts ADD CONSTRAINT accounts_fk_4 FOREIGN KEY (icon_id) REFERENCES coin.icons (id);

ALTER TABLE coin.transactions ADD CONSTRAINT orders_fk FOREIGN KEY (account_from_id) REFERENCES coin.accounts (id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE coin.transactions ADD CONSTRAINT orders_fk_1 FOREIGN KEY (account_to_id) REFERENCES coin.accounts (id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE coin.transactions ADD CONSTRAINT transactions_fk FOREIGN KEY (created_by_user_id) REFERENCES coin.users (id);

ALTER TABLE coin.tags ADD CONSTRAINT tags_fk FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id);

ALTER TABLE coin.tags_to_transaction ADD CONSTRAINT tags_to_transaction_fk FOREIGN KEY (tag_id) REFERENCES coin.tags (id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE coin.tags_to_transaction ADD CONSTRAINT tags_to_transaction_fk_1 FOREIGN KEY (transaction_id) REFERENCES coin.transactions (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE coin.users_to_account_groups ADD CONSTRAINT users_to_account_groups_fk FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id);
ALTER TABLE coin.users_to_account_groups ADD CONSTRAINT users_to_account_groups_fk_1 FOREIGN KEY (user_id) REFERENCES coin.users (id);

-- Создаем индексы для производительности
CREATE INDEX IF NOT EXISTS idx_accounts_account_group_id ON coin.accounts USING btree (account_group_id);
CREATE INDEX IF NOT EXISTS idx_accounts_created_by_user_id ON coin.accounts USING btree (created_by_user_id);
CREATE INDEX IF NOT EXISTS idx_accounts_parent_account_id ON coin.accounts USING btree (parent_account_id) WHERE parent_account_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_transactions_account_from_id ON coin.transactions USING btree (account_from_id);
CREATE INDEX IF NOT EXISTS idx_transactions_account_to_id ON coin.transactions USING btree (account_to_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_by_user_id ON coin.transactions USING btree (created_by_user_id);

CREATE INDEX IF NOT EXISTS idx_tags_account_group_id ON coin.tags USING btree (account_group_id);

CREATE INDEX IF NOT EXISTS idx_tags_to_transaction_tag_id ON coin.tags_to_transaction USING btree (tag_id);
CREATE INDEX IF NOT EXISTS idx_tags_to_transaction_transaction_id ON coin.tags_to_transaction USING btree (transaction_id);

CREATE INDEX IF NOT EXISTS idx_users_to_account_groups_user_id ON coin.users_to_account_groups USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_users_to_account_groups_account_group_id ON coin.users_to_account_groups USING btree (account_group_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Rollback всех изменений будет крайне сложным, поэтому лучше создать backup перед миграцией
-- В случае отката, лучше всего восстановить из backup

RAISE EXCEPTION 'Rolling back UUID migration is not supported. Please restore from backup.';

-- +goose StatementEnd
