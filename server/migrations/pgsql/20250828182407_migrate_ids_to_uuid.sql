-- +goose Up
-- +goose StatementBegin
drop table coin.action_history;
drop table coin.action_types;

-- Удаляем все существующие foreign key constraints
ALTER TABLE coin.accounts DROP CONSTRAINT accounts_fk;
ALTER TABLE coin.accounts DROP CONSTRAINT accounts_fk_1;
ALTER TABLE coin.accounts DROP CONSTRAINT accounts_fk_2;
ALTER TABLE coin.accounts DROP CONSTRAINT accounts_fk_4;
ALTER TABLE coin.accounts DROP CONSTRAINT accounts_fk_5;
ALTER TABLE coin.accounts DROP CONSTRAINT accounts_fk_6;

ALTER TABLE coin.transactions DROP CONSTRAINT orders_fk;
ALTER TABLE coin.transactions DROP CONSTRAINT orders_fk_1;
ALTER TABLE coin.transactions DROP CONSTRAINT orders_fk_3;
ALTER TABLE coin.transactions DROP CONSTRAINT transactions_account_groups_fk;
ALTER TABLE coin.transactions DROP CONSTRAINT transactions_fk;

ALTER TABLE coin.tags DROP CONSTRAINT tags_fk;

ALTER TABLE coin.tags_to_transaction DROP CONSTRAINT tags_to_transaction_fk;
ALTER TABLE coin.tags_to_transaction DROP CONSTRAINT tags_to_transaction_fk_1;

ALTER TABLE coin.users_to_account_groups DROP CONSTRAINT users_to_account_groups_fk;
ALTER TABLE coin.users_to_account_groups DROP CONSTRAINT users_to_account_groups_fk_1;

ALTER TABLE coin.account_groups DROP CONSTRAINT accounts_groups_fk_1;

ALTER TABLE coin.users DROP CONSTRAINT users_fk;

ALTER TABLE coin.devices DROP CONSTRAINT sessions_fk;


-- Добавляем новые колонки с UUID типом
ALTER TABLE coin.users ADD COLUMN id_new UUID DEFAULT gen_random_uuid() not null;
ALTER TABLE coin.account_groups ADD COLUMN id_new UUID DEFAULT gen_random_uuid() not null;
ALTER TABLE coin.account_groups ADD COLUMN created_by_user_id_new UUID;
ALTER TABLE coin.accounts ADD COLUMN id_new UUID DEFAULT gen_random_uuid() not null;
ALTER TABLE coin.accounts ADD COLUMN account_group_id_new UUID;
ALTER TABLE coin.accounts ADD COLUMN parent_account_id_new UUID;
ALTER TABLE coin.accounts ADD COLUMN created_by_user_id_new UUID;
ALTER TABLE coin.accounts ADD COLUMN icon_id_new UUID;
ALTER TABLE coin.transactions ADD COLUMN id_new UUID DEFAULT gen_random_uuid() not null;
ALTER TABLE coin.transactions ADD COLUMN account_from_id_new UUID;
ALTER TABLE coin.transactions ADD COLUMN account_to_id_new UUID;
ALTER TABLE coin.transactions ADD COLUMN created_by_user_id_new UUID;
ALTER TABLE coin.transactions ADD COLUMN account_group_id_new UUID;
ALTER TABLE coin.tags ADD COLUMN id_new UUID DEFAULT gen_random_uuid() not null;
ALTER TABLE coin.tags ADD COLUMN account_group_id_new UUID;
ALTER TABLE coin.tags ADD COLUMN created_by_user_id_new UUID;
ALTER TABLE coin.tags_to_transaction ADD COLUMN tag_id_new UUID;
ALTER TABLE coin.tags_to_transaction ADD COLUMN transaction_id_new UUID;
ALTER TABLE coin.users_to_account_groups ADD COLUMN user_id_new UUID;
ALTER TABLE coin.users_to_account_groups ADD COLUMN account_group_id_new UUID;
ALTER TABLE coin.devices ADD COLUMN id_new UUID DEFAULT gen_random_uuid() not null;
ALTER TABLE coin.devices ADD COLUMN user_id_new UUID;
ALTER TABLE coin.icons ADD COLUMN id_new UUID DEFAULT gen_random_uuid() not null;

-- Обновляем references для user_id в devices
UPDATE coin.devices SET user_id_new = (
    SELECT u.id_new FROM coin.users u WHERE u.id = coin.devices.user_id
);

-- Обновляем references для created_by_user_id в account_groups
UPDATE coin.account_groups SET created_by_user_id_new = (
    SELECT u.id_new FROM coin.users u WHERE u.id = coin.account_groups.created_by_user_id
);

-- Обновляем references для accounts
UPDATE coin.accounts SET account_group_id_new = (
    SELECT ag.id_new FROM coin.account_groups ag WHERE ag.id = coin.accounts.account_group_id
);

UPDATE coin.accounts SET parent_account_id_new = (
    SELECT a.id_new FROM coin.accounts a WHERE a.id = coin.accounts.parent_account_id
) WHERE parent_account_id IS NOT NULL;

UPDATE coin.accounts SET created_by_user_id_new = (
    SELECT u.id_new FROM coin.users u WHERE u.id = coin.accounts.created_by_user_id
);

UPDATE coin.accounts SET icon_id_new = (
    SELECT i.id_new FROM coin.icons i WHERE i.id = coin.accounts.icon_id
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

UPDATE coin.transactions SET account_group_id_new = (
    SELECT ag.id_new FROM coin.account_groups ag WHERE ag.id = coin.transactions.account_group_id
);

-- Обновляем references для tags
UPDATE coin.tags SET account_group_id_new = (
    SELECT ag.id_new FROM coin.account_groups ag WHERE ag.id = coin.tags.account_group_id
);

-- Обновляем references для tags
UPDATE coin.tags SET created_by_user_id_new = (
    SELECT u.id_new FROM coin.users u WHERE u.id = coin.tags.created_by_user_id
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

ALTER TABLE coin.icons DROP CONSTRAINT idx_16418_primary;
ALTER TABLE coin.tags DROP CONSTRAINT idx_16427_primary;
ALTER TABLE coin.transactions DROP CONSTRAINT idx_16438_primary;
ALTER TABLE coin.accounts DROP CONSTRAINT idx_16391_primary;
ALTER TABLE coin.tags_to_transaction DROP CONSTRAINT tags_to_transaction_pk;
ALTER TABLE coin.users DROP CONSTRAINT idx_16450_primary;
ALTER TABLE coin.users DROP CONSTRAINT users_unique;
ALTER TABLE coin.users DROP CONSTRAINT users_unique_1;
ALTER TABLE coin.users DROP CONSTRAINT users_unique_2;
ALTER TABLE coin.devices DROP CONSTRAINT devices_pk;


-- Удаляем старые колонки и переименовываем новые
-- Devices
ALTER TABLE coin.devices DROP COLUMN id;
ALTER TABLE coin.devices DROP COLUMN user_id;
ALTER TABLE coin.devices RENAME COLUMN id_new TO id;
ALTER TABLE coin.devices RENAME COLUMN user_id_new TO user_id;
ALTER TABLE coin.devices ADD PRIMARY KEY (id);
ALTER TABLE coin.devices ALTER COLUMN user_id SET NOT NULL;

-- Icons
ALTER TABLE coin.icons DROP COLUMN id;
ALTER TABLE coin.icons RENAME COLUMN id_new TO id;
ALTER TABLE coin.icons ADD PRIMARY KEY (id);

-- Users
ALTER TABLE coin.users DROP COLUMN id;
ALTER TABLE coin.users RENAME COLUMN id_new TO id;
ALTER TABLE coin.users ADD PRIMARY KEY (id);

-- Account Groups
ALTER TABLE coin.account_groups DROP COLUMN id;
ALTER TABLE coin.account_groups DROP COLUMN created_by_user_id;
ALTER TABLE coin.account_groups RENAME COLUMN id_new TO id;
ALTER TABLE coin.account_groups RENAME COLUMN created_by_user_id_new TO created_by_user_id;
ALTER TABLE coin.account_groups ADD PRIMARY KEY (id);
ALTER TABLE coin.account_groups ALTER COLUMN created_by_user_id SET NOT NULL;

-- Accounts
ALTER TABLE coin.accounts DROP COLUMN id;
ALTER TABLE coin.accounts DROP COLUMN account_group_id;
ALTER TABLE coin.accounts DROP COLUMN parent_account_id;
ALTER TABLE coin.accounts DROP COLUMN created_by_user_id;
ALTER TABLE coin.accounts DROP COLUMN icon_id;
ALTER TABLE coin.accounts RENAME COLUMN id_new TO id;
ALTER TABLE coin.accounts RENAME COLUMN account_group_id_new TO account_group_id;
ALTER TABLE coin.accounts RENAME COLUMN parent_account_id_new TO parent_account_id;
ALTER TABLE coin.accounts RENAME COLUMN created_by_user_id_new TO created_by_user_id;
ALTER TABLE coin.accounts RENAME COLUMN icon_id_new TO icon_id;
ALTER TABLE coin.accounts ADD PRIMARY KEY (id);
ALTER TABLE coin.accounts ALTER COLUMN account_group_id SET NOT NULL;
ALTER TABLE coin.accounts ALTER COLUMN created_by_user_id SET NOT NULL;
ALTER TABLE coin.accounts ALTER COLUMN icon_id SET NOT NULL;

-- Transactions
ALTER TABLE coin.transactions DROP COLUMN id;
ALTER TABLE coin.transactions DROP COLUMN account_from_id;
ALTER TABLE coin.transactions DROP COLUMN account_to_id;
ALTER TABLE coin.transactions DROP COLUMN created_by_user_id;
ALTER TABLE coin.transactions DROP COLUMN account_group_id;
ALTER TABLE coin.transactions RENAME COLUMN id_new TO id;
ALTER TABLE coin.transactions RENAME COLUMN account_from_id_new TO account_from_id;
ALTER TABLE coin.transactions RENAME COLUMN account_to_id_new TO account_to_id;
ALTER TABLE coin.transactions RENAME COLUMN created_by_user_id_new TO created_by_user_id;
ALTER TABLE coin.transactions RENAME COLUMN account_group_id_new TO account_group_id;
ALTER TABLE coin.transactions ADD PRIMARY KEY (id);
ALTER TABLE coin.transactions ALTER COLUMN account_from_id SET NOT NULL;
ALTER TABLE coin.transactions ALTER COLUMN account_to_id SET NOT NULL;
ALTER TABLE coin.transactions ALTER COLUMN created_by_user_id SET NOT NULL;
ALTER TABLE coin.transactions ALTER COLUMN account_group_id SET NOT NULL;

-- Tags
ALTER TABLE coin.tags DROP COLUMN id;
ALTER TABLE coin.tags DROP COLUMN account_group_id;
ALTER TABLE coin.tags DROP COLUMN created_by_user_id;
ALTER TABLE coin.tags RENAME COLUMN id_new TO id;
ALTER TABLE coin.tags RENAME COLUMN account_group_id_new TO account_group_id;
ALTER TABLE coin.tags RENAME COLUMN created_by_user_id_new TO created_by_user_id;
ALTER TABLE coin.tags ADD PRIMARY KEY (id);
ALTER TABLE coin.tags ALTER COLUMN account_group_id SET NOT NULL;
ALTER TABLE coin.tags ALTER COLUMN created_by_user_id SET NOT NULL;

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

-- Добавляем новые типы данных
CREATE TYPE coin.account_type AS ENUM ('debt', 'earnings', 'expense', 'regular', 'balancing');
alter table coin.accounts rename column type_signatura to account_type;
ALTER TABLE coin.accounts ALTER column account_type TYPE coin.account_type USING account_type::coin.account_type;
drop table coin.account_types;

CREATE TYPE coin.transaction_type AS ENUM ('balancing', 'consumption', 'income', 'transfer');
alter table coin.transactions rename column type_signatura to transaction_type;
ALTER TABLE coin.transactions ALTER column transaction_type TYPE coin.transaction_type USING transaction_type::coin.transaction_type;
drop table coin.transaction_types;

alter table coin.currencies rename column signatura to slug;
alter table coin.accounts rename column currency_signatura to currency;
alter table coin.users rename column default_currency_signatura to default_currency;



-- Восстанавливаем foreign key constraints
ALTER TABLE coin.users ADD CONSTRAINT fk_users_currencies FOREIGN KEY (default_currency) REFERENCES coin.currencies (slug);

ALTER TABLE coin.accounts ADD CONSTRAINT fk_accounts_currencies FOREIGN KEY (currency) REFERENCES coin.currencies (slug);
ALTER TABLE coin.accounts ADD CONSTRAINT fk_accounts_account_groups FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id);
ALTER TABLE coin.accounts ADD CONSTRAINT fk_accounts_users FOREIGN KEY (created_by_user_id) REFERENCES coin.users (id);
ALTER TABLE coin.accounts ADD CONSTRAINT fk_accounts_accounts FOREIGN KEY (parent_account_id) REFERENCES coin.accounts (id);
ALTER TABLE coin.accounts ADD CONSTRAINT fk_accounts_icons FOREIGN KEY (icon_id) REFERENCES coin.icons (id);

ALTER TABLE coin.transactions ADD CONSTRAINT fk_transactions_accounts_account_from_id FOREIGN KEY (account_from_id) REFERENCES coin.accounts (id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE coin.transactions ADD CONSTRAINT fk_transactions_accounts_account_to_id FOREIGN KEY (account_to_id) REFERENCES coin.accounts (id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE coin.transactions ADD CONSTRAINT fk_transactions_users FOREIGN KEY (created_by_user_id) REFERENCES coin.users (id);
ALTER TABLE coin.transactions ADD CONSTRAINT fk_transactions_account_groups FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id);

ALTER TABLE coin.account_groups ADD CONSTRAINT fk_account_groups_users FOREIGN KEY (created_by_user_id) REFERENCES coin.users (id);

ALTER TABLE coin.devices ADD CONSTRAINT fk_devices_users FOREIGN KEY (user_id) REFERENCES coin.users (id);

ALTER TABLE coin.tags ADD CONSTRAINT fk_tags_account_groups FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id);
ALTER TABLE coin.tags ADD CONSTRAINT fk_tags_users FOREIGN KEY (created_by_user_id) REFERENCES coin.users (id);


ALTER TABLE coin.tags_to_transaction ADD CONSTRAINT fk_tags_to_transctions_tags FOREIGN KEY (tag_id) REFERENCES coin.tags (id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE coin.tags_to_transaction ADD CONSTRAINT fk_tags_to_transctions_transactions FOREIGN KEY (transaction_id) REFERENCES coin.transactions (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE coin.users_to_account_groups ADD CONSTRAINT fk_users_to_account_groups_account_groups FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id);
ALTER TABLE coin.users_to_account_groups ADD CONSTRAINT fk_users_to_account_groups_users FOREIGN KEY (user_id) REFERENCES coin.users (id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Rollback всех изменений будет крайне сложным, поэтому лучше создать backup перед миграцией
-- В случае отката, лучше всего восстановить из backup

RAISE EXCEPTION 'Rolling back UUID migration is not supported. Please restore from backup.';

-- +goose StatementEnd
