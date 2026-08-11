-- +goose Up
-- +goose StatementBegin
-- Бэкафилл аудит-лога: фиксируем событие создания для всех объектов, которые существовали до внедрения аудит-лога
-- и поэтому не имеют записи о своём создании. Слепок "после" берём как текущее состояние строки на момент миграции -
-- это единственное доступное приближение, так как исторического состояния на момент создания не сохранилось.

INSERT INTO coin.audit_log (entity, method, entity_id, snapshot_after, user_id, device_id, account_group_id, datetime_create)
SELECT 'accountGroup', 'create', t.id::text, to_jsonb(t.*), t.created_by_user_id, '', t.id, t.datetime_create
FROM coin.account_groups t;

INSERT INTO coin.audit_log (entity, method, entity_id, snapshot_after, user_id, device_id, account_group_id, datetime_create)
SELECT 'account', 'create', t.id::text, to_jsonb(t.*), t.created_by_user_id, '', t.account_group_id, t.datetime_create
FROM coin.accounts t;

INSERT INTO coin.audit_log (entity, method, entity_id, snapshot_after, user_id, device_id, account_group_id, datetime_create)
SELECT 'tag', 'create', t.id::text, to_jsonb(t.*), t.created_by_user_id, '', t.account_group_id, t.datetime_create
FROM coin.tags t;

INSERT INTO coin.audit_log (entity, method, entity_id, snapshot_after, user_id, device_id, account_group_id, datetime_create)
SELECT 'transaction', 'create', t.id::text, to_jsonb(t.*), t.created_by_user_id, '', t.account_group_id, t.datetime_create
FROM coin.transactions t;

INSERT INTO coin.audit_log (entity, method, entity_id, snapshot_after, user_id, device_id, account_group_id, datetime_create)
SELECT 'accountBudget', 'create', t.id::text, to_jsonb(t.*), t.created_by_user_id, '', t.account_group_id, t.datetime_create
FROM coin.account_budgets t;

-- У пользователя нет отдельного "создателя" - создатель это он сам. Пароль и его соль в слепок не попадают.
INSERT INTO coin.audit_log (entity, method, entity_id, snapshot_after, user_id, device_id, account_group_id, datetime_create)
SELECT 'user', 'create', t.id::text,
       jsonb_build_object(
           'id', t.id,
           'name', t.name,
           'email', t.email,
           'timeCreate', t.time_create,
           'defaultCurrency', t.default_currency,
           'isAdmin', t.is_admin
       ),
       t.id, '', NULL, t.time_create
FROM coin.users t;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM coin.audit_log WHERE device_id = '' AND method = 'create';
-- +goose StatementEnd
