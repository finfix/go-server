-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS coin.pending_linked_transfers
(
    id                    UUID PRIMARY KEY,                              -- Идентификатор переноса
    status                TEXT                     NOT NULL DEFAULT 'pending', -- Статус переноса (pending/completed/ignored)
    source_transaction_id UUID                     NOT NULL,             -- Транзакция-инициатор
    source_account_id     UUID                     NOT NULL,             -- Счёт-мост со стороны источника
    target_account_id     UUID                     NOT NULL,             -- Счёт-мост со стороны получателя
    account_group_id      UUID                     NOT NULL,             -- Группа-источник (source_account_id/source_transaction_id)
    datetime_create       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    CONSTRAINT fk_pending_linked_transfers_transactions FOREIGN KEY (source_transaction_id) REFERENCES coin.transactions (id),
    CONSTRAINT fk_pending_linked_transfers_source_accounts FOREIGN KEY (source_account_id) REFERENCES coin.accounts (id),
    CONSTRAINT fk_pending_linked_transfers_target_accounts FOREIGN KEY (target_account_id) REFERENCES coin.accounts (id),
    CONSTRAINT fk_pending_linked_transfers_account_groups FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id)
);
COMMENT ON TABLE coin.pending_linked_transfers IS 'Требования довнесения транзакции через счёт-мост (см. coin.accounts.linked_account_id)';
COMMENT ON COLUMN coin.pending_linked_transfers.status IS 'pending / completed / ignored';
COMMENT ON COLUMN coin.pending_linked_transfers.source_transaction_id IS 'Транзакция-инициатор';
COMMENT ON COLUMN coin.pending_linked_transfers.source_account_id IS 'Счёт-мост со стороны источника';
COMMENT ON COLUMN coin.pending_linked_transfers.target_account_id IS 'Счёт-мост со стороны получателя';
COMMENT ON COLUMN coin.pending_linked_transfers.account_group_id IS 'Группа-источник';

CREATE INDEX IF NOT EXISTS idx_pending_linked_transfers_account_group_id ON coin.pending_linked_transfers (account_group_id);
CREATE INDEX IF NOT EXISTS idx_pending_linked_transfers_target_account_id ON coin.pending_linked_transfers (target_account_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS coin.pending_linked_transfers;
-- +goose StatementEnd
