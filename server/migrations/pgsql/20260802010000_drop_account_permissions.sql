-- +goose Up
-- +goose StatementBegin
DROP TABLE permissions.account_permissions;
DROP SCHEMA permissions;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS permissions;

CREATE TABLE IF NOT EXISTS permissions.account_permissions
(
    account_type VARCHAR NOT NULL,
    action_type  VARCHAR NOT NULL,
    access       bool    NOT NULL
);

ALTER TABLE permissions.account_permissions
    ADD CONSTRAINT account_permissions_pk PRIMARY KEY (account_type, action_type);

INSERT INTO permissions.account_permissions (account_type, action_type, access)
VALUES ('general', 'update_currency', TRUE),
       ('general', 'update_parent_account_id', TRUE),
       ('general', 'update_budget', TRUE),
       ('general', 'update_remainder', TRUE),
       ('general', 'create_transaction', TRUE),
       ('regular', 'update_remainder', TRUE),
       ('regular', 'update_currency', TRUE),
       ('regular', 'update_parent_account_id', TRUE),
       ('regular', 'create_transaction', TRUE),
       ('expense', 'update_currency', TRUE),
       ('expense', 'update_parent_account_id', TRUE),
       ('expense', 'create_transaction', TRUE),
       ('earnings', 'update_currency', TRUE),
       ('earnings', 'update_parent_account_id', TRUE),
       ('earnings', 'create_transaction', TRUE),
       ('debt', 'update_remainder', FALSE),
       ('parent', 'update_remainder', TRUE),
       ('parent', 'update_currency', TRUE),
       ('expense', 'update_budget', TRUE),
       ('earnings', 'update_budget', TRUE),
       ('parent', 'update_budget', TRUE),
       ('regular', 'update_budget', FALSE),
       ('expense', 'update_remainder', FALSE),
       ('earnings', 'update_remainder', FALSE),
       ('debt', 'update_budget', TRUE),
       ('debt', 'update_currency', TRUE),
       ('debt', 'update_parent_account_id', TRUE),
       ('debt', 'create_transaction', TRUE),
       ('parent', 'update_parent_account_id', FALSE),
       ('parent', 'create_transaction', FALSE)
ON CONFLICT (account_type, action_type) DO NOTHING;
-- +goose StatementEnd