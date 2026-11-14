-- +goose Up
-- История курсов валют: точка курса действует с valid_at до следующей точки той же валюты.
-- Курс — единиц валюты за 1 USD, как в coin.currencies
CREATE TABLE coin.currency_rates
(
    id               BIGSERIAL   NOT NULL PRIMARY KEY,
    currency         VARCHAR     NOT NULL,                                   -- Код валюты
    valid_at         TIMESTAMPTZ NOT NULL,                                   -- С какого момента действует курс
    rate             NUMERIC     NOT NULL CHECK (rate > 0),                  -- Единиц валюты за 1 USD
    source           VARCHAR     NOT NULL CHECK (source IN ('provider', 'manual')), -- Откуда курс
    account_group_id UUID REFERENCES coin.account_groups (id),               -- Группа, для которой курс; NULL — для всех
    datetime_create  TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMENT ON TABLE coin.currency_rates IS 'История курсов валют к USD';

-- Одна общая точка на валюту и момент, и одна на валюту, момент и группу
CREATE UNIQUE INDEX uq_currency_rates_common ON coin.currency_rates (currency, valid_at) WHERE account_group_id IS NULL;
CREATE UNIQUE INDEX uq_currency_rates_group ON coin.currency_rates (currency, valid_at, account_group_id) WHERE account_group_id IS NOT NULL;

-- Текущие курсы — первая точка истории
INSERT INTO coin.currency_rates (currency, valid_at, rate, source)
SELECT slug, now(), rate, 'provider' FROM coin.currencies WHERE rate > 0;

-- +goose Down
DROP TABLE coin.currency_rates;
