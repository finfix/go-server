-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS coin;


CREATE TABLE IF NOT EXISTS coin.account_permissions
(
    account_type VARCHAR NOT NULL,
    action_type  VARCHAR NOT NULL,
    access       bool    NOT NULL,
    CONSTRAINT account_permissions_pk PRIMARY KEY (account_type, action_type)
);


CREATE TABLE IF NOT EXISTS coin.account_types
(
    signatura VARCHAR(100) NOT NULL, -- Строковый идентификатор типа счета
    name      VARCHAR(100) NULL,     -- Название типа счета
    CONSTRAINT idx_16404_primary PRIMARY KEY (signatura)
);
COMMENT ON TABLE coin.account_types IS 'Типы счетов';
COMMENT ON COLUMN coin.account_types.signatura IS 'Строковый идентификатор типа счета';
COMMENT ON COLUMN coin.account_types.name IS 'Название типа счета';

INSERT INTO coin.account_types (signatura, name)
VALUES ('debt', 'Долговой счет'),
       ('earnings', 'Доходный счет'),
       ('expense', 'Расходный счет'),
       ('regular', 'Обычный счет'),
       ('balancing', 'Балансировочный')
ON CONFLICT (signatura) DO NOTHING;

CREATE TABLE IF NOT EXISTS coin.action_types
(
    type_signatura VARCHAR(100) NOT NULL, -- Название действия
    note           VARCHAR(100) NULL,     -- Заметка о типе действия
    CONSTRAINT idx_16412_primary PRIMARY KEY (type_signatura)
);
COMMENT ON TABLE coin.action_types IS 'Типы действий';
COMMENT ON COLUMN coin.action_types.type_signatura IS 'Название действия';
COMMENT ON COLUMN coin.action_types.note IS 'Заметка о типе действия';

INSERT INTO coin.action_types (type_signatura, note)
VALUES ('account_create', 'Создание счета'),
       ('account_delete', 'Удаление счета'),
       ('account_update', 'Обновление счета'),
       ('transaction_create', 'Создание транзакции'),
       ('transaction_delete', 'Удаление транзакции'),
       ('transaction_update', 'Обновление транзакции'),
       ('user_create', 'Создание пользователя'),
       ('user_update', 'Обновление пользователя')
ON CONFLICT (type_signatura) DO NOTHING;

CREATE TABLE IF NOT EXISTS coin.currencies
(
    signatura VARCHAR(10)    NOT NULL, -- Строковый идентификатор
    name      VARCHAR(100)   NOT NULL, -- Название валюты
    rate      NUMERIC(15, 8) NOT NULL, -- Курс валюты относительно доллара
    symbol    VARCHAR(100)   NOT NULL,
    CONSTRAINT idx_16415_primary PRIMARY KEY (signatura)
);
COMMENT ON TABLE coin.currencies IS 'Валюты';
COMMENT ON COLUMN coin.currencies.signatura IS 'Строковый идентификатор';
COMMENT ON COLUMN coin.currencies.name IS 'Название валюты';
COMMENT ON COLUMN coin.currencies.rate IS 'Курс валюты относительно доллара';

INSERT INTO coin.currencies (signatura, name, rate, symbol)
VALUES ('BUSD', 'Binance USD', 0.99975103, 'BUSD'),
       ('AZN', 'Азербайджанский манат', 1.70000000, '₼'),
       ('KES', 'Кенийский шиллинг', 128.96861463, 'KSh'),
       ('BTC', 'Биткоин', 0.00001684, '฿'),
       ('NZD', 'Новозеландский доллар', 1.65990027, '$'),
       ('MWK', 'Малавийская квача', 1733.32317562, 'MK'),
       ('VEF', 'Венесуэльский боливар', 3658523.65715850, 'Bs'),
       ('ETB', 'Эфиопский быр', 103.02695021, 'B'),
       ('JEP', 'Фунт Джерси', 0.78305781, '£'),
       ('ERN', 'Эритрейская накфа', 15.00000000, 'Nfk'),
       ('AVAX', 'AVAX', 0.04699879, 'AVAX'),
       ('CLP', 'Чилийское песо', 933.54225642, '$'),
       ('MMK', 'Мьянманский кьят', 2094.14377527, 'K'),
       ('NOK', 'Норвежская крона', 10.80301212, 'kr'),
       ('ZMK', 'Замбийская квача, устаревший', 9001.20000000, 'ZK'),
       ('CUP', 'Кубинское песо', 24.00000000, '$'),
       ('MDL', 'Молдавский лей', 17.52701243, 'L'),
       ('AUD', 'Австралийский доллар', 1.51834018, '$'),
       ('BAM', 'Боснийская конвертируемая марка', 1.78722031, 'KM'),
       ('PKR', 'Пакистанская рупия', 278.43254662, '₨'),
       ('USD', 'Доллар США', 1.00000000, '$'),
       ('PEN', 'Перуанский соль', 3.74552061, 'S/'),
       ('TZS', 'Танзанийский шиллинг', 2707.06284508, 'TSh'),
       ('COP', 'Колумбийское песо', 4048.68529477, '$'),
       ('XPD', 'Палладий', 0.00108776, 'XPD'),
       ('RSD', 'Сербский динар', 106.57925557, 'РСД'),
       ('AFN', 'Афгани', 70.92993285, '؋'),
       ('SRD', 'Суринамский доллар', 28.75548573, '$'),
       ('BTN', 'Бутанский нгултрум', 84.02679314, 'Nu'),
       ('CVE', 'Кабо-Верде эскудо', 100.80646233, '$'),
       ('MKD', 'Македонский денар', 56.15260677, 'ден'),
       ('LTC', 'Litecoin', 0.01626250, 'Ł'),
       ('HRK', 'Хорватская куна', 6.62816094, 'kn'),
       ('ALL', 'Албанский лек', 91.30991598, 'L'),
       ('STD', 'Добра Сан-Томе и Принсипи', 22545.17724454, 'Db'),
       ('TOP', 'Тонганская паанга', 2.34943035, 'T$'),
       ('CHF', 'Швейцарский франк', 0.86470011, 'Fr'),
       ('PYG', 'Парагвайский гуарани', 7588.67151978, '₲'),
       ('IDR', 'Индонезийская рупия', 15932.43336788, 'Rp'),
       ('UYU', 'Уругвайское песо', 40.24256584, '$U'),
       ('KMF', 'Коморский франк', 450.69924133, 'Fr'),
       ('BIF', 'Бурундийский франк', 2882.27421162, 'Fr'),
       ('JOD', 'Иорданский динар', 0.71000000, 'د.ا'),
       ('GMD', 'Гамбийский даласи', 57.99727000, 'D'),
       ('DOP', 'Доминиканское песо', 59.60288633, '$'),
       ('SZL', 'Свазилендский лилангени', 18.23786354, 'L'),
       ('XOF', 'Западноафриканский франк', 599.63052122, 'Fr'),
       ('JPY', 'Японская иена', 146.92923934, '¥'),
       ('PGK', 'Папуа-новогвинейская кина', 3.86106056, 'K'),
       ('TTD', 'Доллар Тринидада и Тобаго', 6.79374105, '$'),
       ('XPT', 'Платина', 0.00105993, 'XPT'),
       ('AOA', 'Ангольская кванза', 877.06715268, 'Kz'),
       ('ILS', 'Израильский новый шекель', 3.77614040, '₪'),
       ('ETH', 'Ethereum', 0.00036691, 'Ξ'),
       ('SYP', 'Сирийский фунт', 13002.26570344, '£'),
       ('OP', 'OP', 0.72440759, 'OP'),
       ('STN', 'STN', 22.54518870, 'STN'),
       ('MOP', 'Макаосская патака', 8.01957109, 'P'),
       ('ARS', 'Аргентинское песо', 938.40445782, '$'),
       ('CNY', 'Китайский юань', 7.17254109, '¥'),
       ('GHS', 'Ганский седи', 15.51973270, '₵'),
       ('RON', 'Румынский лей', 4.54770077, 'lei'),
       ('AED', 'ОАЭ Дирхам', 3.67170050, 'د.إ'),
       ('USDT', 'Tether', 1.00006228, '₮'),
       ('MUR', 'Маврикийская рупия', 46.35398796, 'Rs'),
       ('KYD', 'Доллар Каймановых островов', 0.83333000, '$'),
       ('QAR', 'Катарский риал', 3.63894040, 'ر.ق'),
       ('GIP', 'Гибралтарский фунт', 0.78305768, '£'),
       ('USDC', 'USD Coin', 1.00086558, 'USDC'),
       ('EUR', 'Евро', 0.91410012, '€'),
       ('SGD', 'Сингапурский доллар', 1.32267016, '$'),
       ('XCD', 'Восточно-карибский доллар', 2.70000000, '$'),
       ('SAR', 'Саудовский риял', 3.74592038, 'ر.س'),
       ('FJD', 'Фиджийский доллар', 2.26884038, '$'),
       ('SOS', 'Сомалийский шиллинг', 570.93260820, 'S'),
       ('MNT', 'Монгольский тугрик', 3400.03303623, '₮'),
       ('NIO', 'Никарагуанская кордоба', 36.79611249, 'C$'),
       ('BYN', 'Белорусский рубль', 3.27021358, 'Br'),
       ('MRO', 'Мавританская угия', 356.99982800, 'UM'),
       ('LVL', 'Латвийский лат', 0.64264140, 'Ls'),
       ('BBD', 'Барбадосский доллар', 2.00000000, '$'),
       ('DOT', 'Polkadot', 0.21809671, 'DOT'),
       ('THB', 'Таиландский бат', 35.08057642, '฿'),
       ('SCR', 'Сейшельская рупия', 14.74311242, '₨'),
       ('XRP', 'Ripple', 1.75869342, 'XRP'),
       ('RUB', 'Российский рубль', 90.20750112, '₽'),
       ('BHD', 'Бахрейнский динар', 0.37600000, '.د.ب'),
       ('MXN', 'Мексиканское песо', 19.04690366, '$'),
       ('BNB', 'Binance Coin', 0.00192803, 'BNB'),
       ('IMP', 'Фунт острова Мэн', 0.78305739, '£'),
       ('NPR', 'Непальская рупия', 134.11683894, 'रू'),
       ('MGA', 'Малагасийский ариари', 4583.16684183, 'Ar'),
       ('LSL', 'Лоти', 18.24773297, 'L'),
       ('SOL', 'Соль', 0.00683460, 'SOL'),
       ('PHP', 'Филиппинское песо', 57.27873596, '₱'),
       ('KZT', 'Казахский тенге', 478.96030095, '₸'),
       ('CDF', 'Конголезский франк', 2683.69127057, 'Fr'),
       ('CAD', 'Канадский доллар', 1.37386019, '$'),
       ('LRD', 'Либерийский доллар', 195.44475413, '$'),
       ('GYD', 'Гайанский доллар', 208.76843923, '$'),
       ('DAI', 'Dai', 1.00077717, 'DAI'),
       ('VND', 'Вьетнамский донг', 25112.12472347, '₫'),
       ('XAU', 'Золото', 0.00040390, 'XAU'),
       ('XAG', 'Серебро', 0.03568859, 'XAG'),
       ('MRU', 'MRU', 39.65386917, 'MRU'),
       ('UZS', 'Узбекский сум', 12628.95602386, 'soʻm'),
       ('GBP', 'Британский фунт стерлингов', 0.78296012, '£'),
       ('PAB', 'Панамская бальбоа', 0.99872014, 'B/-'),
       ('EGP', 'Египетский фунт', 49.31287701, '£'),
       ('TWD', 'Новый тайваньский доллар', 32.42319390, '$'),
       ('LBP', 'Ливанский фунт', 89518.49019085, '£'),
       ('BZD', 'Белизский доллар', 2.00000000, '$'),
       ('IQD', 'Иракский динар', 1307.72805952, 'ع.د'),
       ('NGN', 'Нигерийская найра', 1583.08827520, '₦'),
       ('DKK', 'Датская крона', 6.82210121, 'kr'),
       ('GNF', 'Гвинейский франк', 8608.04975434, 'Fr'),
       ('LKR', 'Шри-ланкийская рупия', 299.01116274, 'Rs'),
       ('TND', 'Тунисский динар', 3.07414043, 'د.ت'),
       ('GGP', 'Гернсийский фунт', 0.78305725, '£'),
       ('YER', 'Йеменский риал', 249.80350742, '﷼'),
       ('HTG', 'Гаитянский гурд', 133.85297784, 'G'),
       ('MATIC', 'Polygon', 2.34657861, 'MATIC'),
       ('BOB', 'Боливийский боливиано', 6.93319128, 'R$'),
       ('CUC', 'Кубинское конвертируемое песо', 1.00000000, '$'),
       ('DJF', 'Франк Джибути', 177.72100000, 'Fr'),
       ('KRW', 'Южнокорейская вона', 1366.32490782, '₩'),
       ('BSD', 'Багамский доллар', 1.00000000, '$'),
       ('AMD', 'Армянский драм', 387.85240346, '֏'),
       ('ISK', 'Исландская крона', 138.16866515, 'kr'),
       ('HNL', 'Гондурасская лемпира', 24.75720344, 'L'),
       ('KHR', 'Камбоджийский риель', 4083.32352530, '៛'),
       ('VUV', 'Вануатский вату', 119.91099196, 'VT'),
       ('VES', 'VES', 36.58523178, 'VES'),
       ('CZK', 'Чешская крона', 22.99906339, 'Kč'),
       ('RWF', 'Франк Руанды', 1317.60115787, 'Fr'),
       ('BWP', 'Ботсванская пула', 13.51746149, 'P'),
       ('SVC', 'Сальвадорский колон', 8.75000000, '$'),
       ('WST', 'Самоанская тала', 2.75235358, 'WS$'),
       ('HKD', 'Гонконгский доллар', 7.78968085, '$'),
       ('ZMW', 'Замбийская квача', 26.16693262, 'ZK'),
       ('KPW', 'Северокорейская вона', 900.00371699, '₩'),
       ('XAF', 'Центральноафриканский франк', 599.63054924, 'Fr'),
       ('LAK', 'Лаосский кип', 22169.04406969, '₭'),
       ('MYR', 'Малайзийский ринггит', 4.45359073, 'RM'),
       ('HUF', 'Венгерский форинт', 360.01494733, 'Ft'),
       ('ZWL', 'Зимбабвийский доллар', 34349.79381171, '$'),
       ('XPF', 'Французский тихоокеанский франк', 109.01817510, 'Fr'),
       ('NAD', 'Намибийский доллар', 18.17883220, '$'),
       ('TMT', 'Туркменский манат', 3.50000000, 'T'),
       ('AWG', 'Арубанский флорин', 1.79000000, 'ƒ'),
       ('MVR', 'Мальдивская руфия', 15.44629263, 'Rf'),
       ('BDT', 'Бангладешская така', 117.97841118, '৳'),
       ('DZD', 'Алжирский динар', 134.78627055, 'د.ج'),
       ('BGN', 'Болгарский лев', 1.78722025, 'лв'),
       ('OMR', 'Оманский риал', 0.38408006, 'ر.ع.'),
       ('IRR', 'Иранский риал', 41996.24831249, '﷼'),
       ('FKP', 'Фунт Фолклендских островов', 0.78305740, '£'),
       ('SLL', 'Сьерра-леонский леоне', 22419.43706382, 'Le'),
       ('CLF', 'Условная расчетная единица Чили', 0.02443000, 'UF'),
       ('LYD', 'Ливийский динар', 4.80539093, 'ل.د'),
       ('TRY', 'Турецкая лира', 33.48896651, '₺'),
       ('JMD', 'Ямайский доллар', 156.99187712, '$'),
       ('BYR', 'Белорусский рубль, устаревший', 32702.15315877, 'Br'),
       ('MAD', 'Марокканский дирхам', 9.79721102, 'د.م.'),
       ('XDR', 'Специальные права заимствования', 0.74806014, 'SDR'),
       ('BMD', 'Бермудский доллар', 1.00000000, '$'),
       ('BRL', 'Бразильский реал', 5.49223101, 'Bs'),
       ('MZN', 'Мозамбикский метикал', 63.54863078, 'MT'),
       ('INR', 'Индийская рупия', 83.88429665, '₹'),
       ('BND', 'Брунейский доллар', 1.32265025, '$'),
       ('GTQ', 'Гватемальский кетсаль', 7.72624123, 'Q'),
       ('SHP', 'Фунт Святой Елены', 0.78296009, '£'),
       ('TJS', 'Таджикский сомони', 10.57563111, 'SM'),
       ('ZAR', 'Южноафриканский рэнд', 18.24444196, 'R'),
       ('KGS', 'Киргизский сом', 85.33747179, 'сом'),
       ('KWD', 'Кувейтский динар', 0.30621005, 'د.ك'),
       ('SEK', 'Шведская крона', 10.51785209, 'kr'),
       ('PLN', 'Польский злотый', 3.93667047, 'zł'),
       ('CRC', 'Коста-Риканский колон', 528.29879207, '₡'),
       ('GEL', 'Грузинский лари', 2.69346028, '₾'),
       ('SDG', 'Суданский фунт', 601.50000000, '£'),
       ('ANG', 'Нидерландский антильский гульден', 1.78632032, 'ƒ'),
       ('ADA', 'Кардано', 2.95123043, '₳'),
       ('LTL', 'Литовский лит', 3.15724441, 'Lt'),
       ('UAH', 'Украинская гривна', 41.23685460, '₴'),
       ('ARB', 'ARB', 1.69704841, 'ARB'),
       ('UGX', 'Угандийский шиллинг', 3724.85972651, 'USh'),
       ('SBD', 'Доллар Соломоновых островов', 8.33477925, '$')
ON CONFLICT (signatura) DO NOTHING;

CREATE TABLE IF NOT EXISTS coin.icons
(
    id   int8         NOT NULL, -- Строковый идентификатор
    img  VARCHAR(100) NULL,     -- Ссылка на изображение
    name VARCHAR(100) NULL,     -- Название изображения
    CONSTRAINT idx_16418_primary PRIMARY KEY (id)
);

COMMENT ON TABLE coin.icons IS 'Иконки счетов';
COMMENT ON COLUMN coin.icons.id IS 'Строковый идентификатор';
COMMENT ON COLUMN coin.icons.img IS 'Ссылка на изображение';
COMMENT ON COLUMN coin.icons.name IS 'Название изображения';

INSERT INTO coin.icons (id, img, name)
VALUES (1, 'dollar.png', 'Кошелек')
ON CONFLICT (id) DO NOTHING;



CREATE TABLE IF NOT EXISTS coin.transaction_types
(
    signatura VARCHAR(100) NOT NULL, -- Строковый идентификатор типа операции
    name      VARCHAR(100) NULL,     -- Название типа операции
    CONSTRAINT idx_16446_primary PRIMARY KEY (signatura)
);
COMMENT ON TABLE coin.transaction_types IS 'Типы счетов';
COMMENT ON COLUMN coin.transaction_types.signatura IS 'Строковый идентификатор типа операции';
COMMENT ON COLUMN coin.transaction_types.name IS 'Название типа операции';

INSERT INTO coin.transaction_types (signatura, name)
VALUES ('balancing', 'Балансировка'),
       ('consumption', 'Расход'),
       ('income', 'Доход'),
       ('transfer', 'Перевод')
ON CONFLICT (signatura) DO NOTHING;

CREATE TABLE IF NOT EXISTS coin.account_groups
(
    id                 BIGSERIAL    NOT NULL, -- Идектификатор группы
    name               VARCHAR(100) NULL,     -- Название группы
    available_budget   int4         NULL,
    currency_signatura VARCHAR      NOT NULL,
    serial_number      int4         NOT NULL,
    visible            bool         NOT NULL,
    datetime_create    timestamptz  NOT NULL, -- Дата и время создания группы счетов
    created_by_user_id int8         NOT NULL,
    CONSTRAINT idx_16400_primary PRIMARY KEY (id),
    CONSTRAINT accounts_groups_fk_1 FOREIGN KEY (currency_signatura) REFERENCES coin.currencies (signatura)
);
COMMENT ON TABLE coin.account_groups IS 'Группы аккаунтов';
COMMENT ON COLUMN coin.account_groups.id IS 'Идектификатор группы';
COMMENT ON COLUMN coin.account_groups.name IS 'Название группы';
COMMENT ON COLUMN coin.account_groups.datetime_create IS 'Дата и время создания группы счетов';


CREATE TABLE IF NOT EXISTS coin.tags
(
    id                 BIGSERIAL    NOT NULL, -- Идентификатор подкатегории
    name               VARCHAR(100) NOT NULL, -- Название подкатегории
    account_group_id   int8         NOT NULL, -- Идентификатор группы счетов
    created_by_user_id int8         NOT NULL, -- Идентификатор пользователя, создавшего подкатегорию
    datetime_create    timestamptz  NOT NULL,
    CONSTRAINT idx_16427_primary PRIMARY KEY (id),
    CONSTRAINT tags_fk FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id)
);
CREATE INDEX IF NOT EXISTS idx_16427_tags_fk ON coin.tags USING btree (account_group_id);
COMMENT ON TABLE coin.tags IS 'Подкатегории';
COMMENT ON COLUMN coin.tags.id IS 'Идентификатор подкатегории';
COMMENT ON COLUMN coin.tags.name IS 'Название подкатегории';
COMMENT ON COLUMN coin.tags.account_group_id IS 'Идентификатор группы счетов';
COMMENT ON COLUMN coin.tags.created_by_user_id IS 'Идентификатор пользователя, создавшего подкатегорию';


CREATE TABLE IF NOT EXISTS coin.users
(
    id                         BIGSERIAL          NOT NULL, -- Идентификато пользователя
    name                       VARCHAR(100)       NOT NULL, -- Имя пользователя
    email                      VARCHAR(30)        NOT NULL, -- Почтовый адрес пользователя
    password_hash              bytea              NOT NULL, -- Хэш пароля
    time_create                timestamptz        NOT NULL, -- Дата создания аккаунта
    default_currency_signatura VARCHAR            NOT NULL,
    password_salt              bytea              NOT NULL,
    is_admin                   bool DEFAULT FALSE NOT NULL,
    CONSTRAINT idx_16450_primary PRIMARY KEY (id),
    CONSTRAINT users_unique UNIQUE (email),
    CONSTRAINT users_unique_1 UNIQUE (password_hash),
    CONSTRAINT users_unique_2 UNIQUE (password_salt),
    CONSTRAINT users_fk FOREIGN KEY (default_currency_signatura) REFERENCES coin.currencies (signatura)
);
COMMENT ON TABLE coin.users IS 'Участники проекта';
COMMENT ON COLUMN coin.users.id IS 'Идентификато пользователя';
COMMENT ON COLUMN coin.users.name IS 'Имя пользователя';
COMMENT ON COLUMN coin.users.email IS 'Почтовый адрес пользователя';
COMMENT ON COLUMN coin.users.password_hash IS 'Хэш пароля';
COMMENT ON COLUMN coin.users.time_create IS 'Дата создания аккаунта';


CREATE TABLE IF NOT EXISTS coin.users_to_account_groups
(
    user_id          int4 NOT NULL,
    account_group_id int4 NOT NULL,
    CONSTRAINT users_to_users_group_pk PRIMARY KEY (user_id, account_group_id),
    CONSTRAINT users_to_account_groups_fk FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id),
    CONSTRAINT users_to_account_groups_fk_1 FOREIGN KEY (user_id) REFERENCES coin.users (id)
);


CREATE TABLE IF NOT EXISTS coin.accounts
(
    id                     BIGSERIAL                             NOT NULL, -- Идентификатор счета
    budget_amount          NUMERIC DEFAULT '0'::DOUBLE PRECISION NULL,     -- Бюджет на месяц
    name                   VARCHAR(100)                          NOT NULL, -- Название
    icon_id                int8    DEFAULT 1                     NOT NULL, -- Идентификатор иконки
    type_signatura         VARCHAR(100)                          NOT NULL, -- Тип счета
    currency_signatura     VARCHAR(10)                           NOT NULL, -- Строковый идентификатор валюты
    visible                bool    DEFAULT TRUE                  NOT NULL, -- Видимость счета
    account_group_id       int8                                  NOT NULL, -- Идентификатор группы счетов
    accounting_in_header   bool    DEFAULT TRUE                  NOT NULL, -- Учитывать ли счет в шапке
    parent_account_id      int8                                  NULL,     -- Идентификатор привязки к другому счету
    serial_number          int8                                  NOT NULL, -- Порядковый номер счета
    budget_gradual_filling bool    DEFAULT TRUE                  NOT NULL, -- Заполняется ли бюджет постепенно
    is_parent              bool    DEFAULT FALSE                 NOT NULL, -- Является ли счет родительским
    budget_fixed_sum       NUMERIC DEFAULT 0                     NOT NULL,
    budget_days_offset     int8    DEFAULT 0                     NOT NULL,
    datetime_create        timestamptz                           NOT NULL, -- Дата и время создания счета
    created_by_user_id     int8                                  NOT NULL, -- Каким пользователем был создан счет
    accounting_in_charts   bool    DEFAULT TRUE                  NOT NULL, -- Учитывать ли счет в графиках
    CONSTRAINT idx_16391_primary PRIMARY KEY (id),
    CONSTRAINT accounts_fk FOREIGN KEY (account_group_id) REFERENCES coin.account_groups (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT accounts_fk_1 FOREIGN KEY (type_signatura) REFERENCES coin.account_types (signatura) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT accounts_fk_2 FOREIGN KEY (currency_signatura) REFERENCES coin.currencies (signatura) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT accounts_fk_4 FOREIGN KEY (icon_id) REFERENCES coin.icons (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT accounts_fk_5 FOREIGN KEY (parent_account_id) REFERENCES coin.accounts (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT accounts_fk_6 FOREIGN KEY (created_by_user_id) REFERENCES coin.users (id)
);
CREATE INDEX IF NOT EXISTS idx_16391_account_fk ON coin.accounts USING btree (icon_id);
CREATE INDEX IF NOT EXISTS idx_16391_accounts_fk ON coin.accounts USING btree (account_group_id);
CREATE INDEX IF NOT EXISTS idx_16391_accounts_fk_1 ON coin.accounts USING btree (type_signatura);
CREATE INDEX IF NOT EXISTS idx_16391_accounts_fk_2 ON coin.accounts USING btree (currency_signatura);
CREATE INDEX IF NOT EXISTS idx_16391_accounts_fk_5 ON coin.accounts USING btree (parent_account_id);
COMMENT ON TABLE coin.accounts IS 'Счета';
COMMENT ON COLUMN coin.accounts.id IS 'Идентификатор счета';
COMMENT ON COLUMN coin.accounts.budget_amount IS 'Бюджет на месяц';
COMMENT ON COLUMN coin.accounts.name IS 'Название';
COMMENT ON COLUMN coin.accounts.icon_id IS 'Идентификатор иконки';
COMMENT ON COLUMN coin.accounts.type_signatura IS 'Тип счета';
COMMENT ON COLUMN coin.accounts.currency_signatura IS 'Строковый идентификатор валюты';
COMMENT ON COLUMN coin.accounts.visible IS 'Видимость счета';
COMMENT ON COLUMN coin.accounts.account_group_id IS 'Идентификатор группы счетов';
COMMENT ON COLUMN coin.accounts.accounting_in_header IS 'Учитывать ли счет в шапке';
COMMENT ON COLUMN coin.accounts.parent_account_id IS 'Идентификатор привязки к другому счету';
COMMENT ON COLUMN coin.accounts.serial_number IS 'Порядковый номер счета';
COMMENT ON COLUMN coin.accounts.budget_gradual_filling IS 'Заполняется ли бюджет постепенно';
COMMENT ON COLUMN coin.accounts.is_parent IS 'Является ли счет родительским';
COMMENT ON COLUMN coin.accounts.datetime_create IS 'Дата и время создания счета';
COMMENT ON COLUMN coin.accounts.created_by_user_id IS 'Каким пользователем был создан счет';
COMMENT ON COLUMN coin.accounts.accounting_in_charts IS 'Учитывать ли счет в графиках';


CREATE TABLE IF NOT EXISTS coin.action_history
(
    id                    BIGSERIAL    NOT NULL, -- Идентификатор действия
    action_type_signatura VARCHAR(100) NOT NULL, -- Тип действия пользователя
    user_id               int8         NOT NULL, -- Пользователь, который произвел действие
    object_id             int8         NULL,     -- Идентификатор измененного объекта
    note                  VARCHAR(200) NULL,     -- Заметка от администратора
    action_time           timestamptz  NOT NULL, -- Время, когда совершилось действие
    CONSTRAINT idx_16408_primary PRIMARY KEY (id),
    CONSTRAINT action_history_fk FOREIGN KEY (action_type_signatura) REFERENCES coin.action_types (type_signatura) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT newtable_fk FOREIGN KEY (user_id) REFERENCES coin.users (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_16408_action_history_fk ON coin.action_history USING btree (action_type_signatura);
CREATE INDEX IF NOT EXISTS idx_16408_newtable_fk ON coin.action_history USING btree (user_id);
COMMENT ON TABLE coin.action_history IS 'История действий';
COMMENT ON COLUMN coin.action_history.id IS 'Идентификатор действия';
COMMENT ON COLUMN coin.action_history.action_type_signatura IS 'Тип действия пользователя';
COMMENT ON COLUMN coin.action_history.user_id IS 'Пользователь, который произвел действие';
COMMENT ON COLUMN coin.action_history.object_id IS 'Идентификатор измененного объекта';
COMMENT ON COLUMN coin.action_history.note IS 'Заметка от администратора';
COMMENT ON COLUMN coin.action_history.action_time IS 'Время, когда совершилось действие';


CREATE TABLE IF NOT EXISTS coin.devices
(
    id                    BIGSERIAL    NOT NULL,
    refresh_token         VARCHAR(200) NOT NULL,
    device_id             VARCHAR(100) NOT NULL,
    user_id               int8         NOT NULL,
    notification_token    VARCHAR      NULL, -- Токен для системы уведомлений
    application_bundle_id VARCHAR      NOT NULL,
    device_ip_address     VARCHAR      NOT NULL,
    device_user_agent     VARCHAR      NOT NULL,
    device_os_name        VARCHAR      NOT NULL,
    device_os_version     VARCHAR      NOT NULL,
    device_name           VARCHAR      NOT NULL,
    device_model_name     VARCHAR      NOT NULL,
    application_version   VARCHAR      NOT NULL,
    application_build     VARCHAR      NOT NULL,
    CONSTRAINT devices_pk PRIMARY KEY (device_id, user_id),
    CONSTRAINT sessions_fk FOREIGN KEY (user_id) REFERENCES coin.users (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_16422_sessions_fk ON coin.devices USING btree (user_id);
COMMENT ON TABLE coin.devices IS 'Сессии';
COMMENT ON COLUMN coin.devices.notification_token IS 'Токен для системы уведомлений';


CREATE TABLE IF NOT EXISTS coin.transactions
(
    id                   BIGSERIAL         NOT NULL, -- Идентификатор транзакции
    date_transaction     DATE              NOT NULL, -- Дата совершения транзакции
    type_signatura       VARCHAR(100)      NOT NULL, -- Тип транзакции
    amount_from          NUMERIC(17, 7)    NOT NULL, -- Сумма, ушедшая из первого счета
    amount_to            NUMERIC(17, 7)    NOT NULL, -- Сумма, пришедшая во второй счет
    note                 VARCHAR(500)      NOT NULL, -- Заметка
    account_from_id      int8              NOT NULL, -- Откуда взяли деньги
    account_to_id        int8              NOT NULL, -- Куда положили деньги
    is_executed          bool DEFAULT TRUE NOT NULL, -- Исполнена ли операция
    datetime_create      timestamptz       NOT NULL, -- Дата создания транзакции
    accounting_in_charts bool DEFAULT TRUE NOT NULL,
    created_by_user_id   int8              NOT NULL, -- Пользователь, создавший транзакцию
    CONSTRAINT idx_16438_primary PRIMARY KEY (id),
    CONSTRAINT orders_fk FOREIGN KEY (account_from_id) REFERENCES coin.accounts (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT orders_fk_1 FOREIGN KEY (account_to_id) REFERENCES coin.accounts (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT orders_fk_3 FOREIGN KEY (type_signatura) REFERENCES coin.transaction_types (signatura) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT transactions_fk FOREIGN KEY (created_by_user_id) REFERENCES coin.users (id)
);
CREATE INDEX IF NOT EXISTS idx_16438_orders_fk ON coin.transactions USING btree (account_from_id);
CREATE INDEX IF NOT EXISTS idx_16438_orders_fk_1 ON coin.transactions USING btree (account_to_id);
CREATE INDEX IF NOT EXISTS idx_16438_orders_fk_3 ON coin.transactions USING btree (type_signatura);
COMMENT ON TABLE coin.transactions IS 'Транзакции';
COMMENT ON COLUMN coin.transactions.id IS 'Идентификатор транзакции';
COMMENT ON COLUMN coin.transactions.date_transaction IS 'Дата совершения транзакции';
COMMENT ON COLUMN coin.transactions.type_signatura IS 'Тип транзакции';
COMMENT ON COLUMN coin.transactions.amount_from IS 'Сумма, ушедшая из первого счета';
COMMENT ON COLUMN coin.transactions.amount_to IS 'Сумма, пришедшая во второй счет';
COMMENT ON COLUMN coin.transactions.note IS 'Заметка';
COMMENT ON COLUMN coin.transactions.account_from_id IS 'Откуда взяли деньги';
COMMENT ON COLUMN coin.transactions.account_to_id IS 'Куда положили деньги';
COMMENT ON COLUMN coin.transactions.is_executed IS 'Исполнена ли операция';
COMMENT ON COLUMN coin.transactions.datetime_create IS 'Дата создания транзакции';
COMMENT ON COLUMN coin.transactions.created_by_user_id IS 'Пользователь, создавший транзакцию';


CREATE TABLE IF NOT EXISTS coin.tags_to_transaction
(
    transaction_id int8 NOT NULL, -- Идентификатор транзакции
    tag_id         int8 NOT NULL, -- Идентификатор подкатегории
    CONSTRAINT tags_to_transaction_pk PRIMARY KEY (transaction_id, tag_id),
    CONSTRAINT tags_to_transaction_fk FOREIGN KEY (tag_id) REFERENCES coin.tags (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT tags_to_transaction_fk_1 FOREIGN KEY (transaction_id) REFERENCES coin.transactions (id) ON DELETE RESTRICT ON UPDATE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_16433_tags_fk ON coin.tags_to_transaction USING btree (transaction_id);
CREATE INDEX IF NOT EXISTS idx_16433_tags_to_order_fk ON coin.tags_to_transaction USING btree (tag_id);
COMMENT ON TABLE coin.tags_to_transaction IS 'Подкатегории в транзакции';
COMMENT ON COLUMN coin.tags_to_transaction.transaction_id IS 'Идентификатор транзакции';
COMMENT ON COLUMN coin.tags_to_transaction.tag_id IS 'Идентификатор подкатегории';


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


CREATE SCHEMA IF NOT EXISTS settings;


CREATE TABLE IF NOT EXISTS settings.versions
(
    name    VARCHAR NOT NULL,
    version VARCHAR NOT NULL
);

ALTER TABLE settings.versions
    ADD CONSTRAINT versions_pk PRIMARY KEY (name);

INSERT INTO settings.versions (name, version)
VALUES ('ios', 'v1.1.2')
ON CONFLICT (name) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE settings.versions;
DROP SCHEMA settings;
DROP TABLE permissions.account_permissions;
DROP SCHEMA permissions;
DROP TABLE coin.tags_to_transaction;
DROP TABLE coin.transactions;
DROP TABLE coin.devices;
DROP TABLE coin.action_history;
DROP TABLE coin.accounts;
DROP TABLE coin.users_to_account_groups;
DROP TABLE coin.users;
DROP TABLE coin.tags;
DROP TABLE coin.account_groups;
DROP TABLE coin.transaction_types;
DROP TABLE coin.icons;
DROP TABLE coin.currencies;
DROP TABLE coin.action_types;
DROP TABLE coin.account_types;
DROP TABLE coin.account_permissions;
DROP SCHEMA coin;
-- +goose StatementEnd
