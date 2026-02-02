CREATE TABLE IF NOT EXISTS currencies (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(3) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    sign VARCHAR(3) NOT NULL
);

CREATE TABLE IF NOT EXISTS exchange_rates (
    id BIGSERIAL PRIMARY KEY,
    base_currency_id BIGINT NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    target_currency_id BIGINT NOT NULL REFERENCES currencies(id) ON DELETE CASCADE,
    rate NUMERIC(20, 8) NOT NULL,
    UNIQUE (base_currency_id, target_currency_id)
);

INSERT INTO currencies (code, full_name, sign) VALUES
    ('USD', 'US Dollar', '$'),
    ('EUR', 'Euro', '€'),
    ('GBP', 'British Pound', '£'),
    ('JPY', 'Japanese Yen', '¥'),
    ('CHF', 'Swiss Franc', '₣'),
    ('CNY', 'Chinese Yuan', '¥'),
    ('AUD', 'Australian Dollar', 'A$'),
    ('CAD', 'Canadian Dollar', 'C$'),
    ('SEK', 'Swedish Krona', 'kr'),
    ('NOK', 'Norwegian Krone', 'kr'),
    ('RUB', 'Russian Ruble', '₽'),
    ('BYN', 'Belarusian Ruble', 'Br')
ON CONFLICT (code) DO NOTHING;

INSERT INTO exchange_rates (base_currency_id, target_currency_id, rate) VALUES
    ((SELECT id FROM currencies WHERE code = 'USD'), (SELECT id FROM currencies WHERE code = 'EUR'), 0.92),
    ((SELECT id FROM currencies WHERE code = 'EUR'), (SELECT id FROM currencies WHERE code = 'USD'), 1.09),
    ((SELECT id FROM currencies WHERE code = 'USD'), (SELECT id FROM currencies WHERE code = 'GBP'), 0.78),
    ((SELECT id FROM currencies WHERE code = 'GBP'), (SELECT id FROM currencies WHERE code = 'USD'), 1.28),
    ((SELECT id FROM currencies WHERE code = 'USD'), (SELECT id FROM currencies WHERE code = 'JPY'), 145.30),
    ((SELECT id FROM currencies WHERE code = 'EUR'), (SELECT id FROM currencies WHERE code = 'CHF'), 0.97)
ON CONFLICT (base_currency_id, target_currency_id) DO NOTHING;
