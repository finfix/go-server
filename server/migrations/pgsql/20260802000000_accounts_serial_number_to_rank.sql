-- +goose Up
-- +goose StatementBegin
ALTER TABLE coin.accounts RENAME COLUMN serial_number TO rank;
ALTER TABLE coin.accounts ALTER COLUMN rank TYPE VARCHAR(255) USING lpad(rank::text, 20, '0');
COMMENT ON COLUMN coin.accounts.rank IS 'Ранг для сортировки счетов (лексикографический, задаётся клиентом)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coin.accounts ALTER COLUMN rank TYPE int8 USING dense_rank() OVER (PARTITION BY account_group_id ORDER BY rank);
ALTER TABLE coin.accounts RENAME COLUMN rank TO serial_number;
COMMENT ON COLUMN coin.accounts.serial_number IS 'Порядковый номер счета';
-- +goose StatementEnd
