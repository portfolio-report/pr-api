-- Create Extension
CREATE EXTENSION IF NOT EXISTS pg_trgm ;

-- Create Enums
CREATE TYPE "portfolio_account_type" AS ENUM ('deposit', 'securities');
CREATE TYPE "portfolio_transaction_type" AS ENUM ('Payment', 'CurrencyTransfer', 'DepositInterest', 'DepositFee', 'DepositTax', 'SecuritiesOrder', 'SecuritiesDividend', 'SecuritiesFee', 'SecuritiesTax', 'SecuritiesTransfer');
CREATE TYPE "portfolio_transaction_unit_type" AS ENUM ('base', 'tax', 'fee');

-- Create Tables
CREATE TABLE "portfolios_accounts" (
    "type" "portfolio_account_type" NOT NULL,
    "name" VARCHAR NOT NULL,
    "uuid" UUID NOT NULL,
    "currency_code" CHAR(3),
    "active" BOOLEAN NOT NULL,
    "note" VARCHAR NOT NULL DEFAULT E'',
    "portfolio_id" INTEGER NOT NULL,
    "reference_account_uuid" UUID,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ("portfolio_id", "uuid")
);


CREATE TABLE "clientupdates" (
    "id" SERIAL NOT NULL,
    "timestamp" TIMESTAMP(6) NOT NULL,
    "version" TEXT NOT NULL,
    "country" TEXT,
    "useragent" TEXT,

    PRIMARY KEY ("id")
);

CREATE TABLE "currencies" (
    "code" CHAR(3) NOT NULL,

    PRIMARY KEY ("code")
);

CREATE TABLE "events" (
    "id" SERIAL NOT NULL,
    "date" DATE NOT NULL,
    "type" VARCHAR(10) NOT NULL,
    "amount" DECIMAL(10,4),
    "currency_code" CHAR(3),
    "ratio" VARCHAR(10),
    "security_uuid" UUID NOT NULL,

    PRIMARY KEY ("id")
);

CREATE TABLE "exchangerates" (
    "id" SERIAL NOT NULL,
    "base_currency_code" CHAR(3) NOT NULL,
    "quote_currency_code" CHAR(3) NOT NULL,

    PRIMARY KEY ("id")
);

CREATE TABLE "exchangerates_prices" (
    "exchangerate_id" INTEGER NOT NULL,
    "date" DATE NOT NULL,
    "value" DECIMAL(16,8) NOT NULL,

    PRIMARY KEY ("exchangerate_id","date")
);

CREATE TABLE "markets" (
    "code" TEXT NOT NULL,
    "name" TEXT NOT NULL,

    PRIMARY KEY ("code")
);

CREATE TABLE "portfolios" (
    "id" SERIAL NOT NULL,
    "name" VARCHAR NOT NULL,
    "note" VARCHAR NOT NULL,
    "base_currency_code" CHAR(3) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "user_id" INTEGER NOT NULL,

    PRIMARY KEY ("id")
);

CREATE TABLE "portfolios_securities" (
    "name" VARCHAR NOT NULL,
    "uuid" UUID NOT NULL,
    "currency_code" CHAR(3) NOT NULL,
    "isin" VARCHAR NOT NULL DEFAULT E'',
    "wkn" VARCHAR NOT NULL DEFAULT E'',
    "symbol" VARCHAR NOT NULL DEFAULT E'',
    "active" BOOLEAN NOT NULL,
    "note" VARCHAR NOT NULL DEFAULT E'',
    "portfolio_id" INTEGER NOT NULL,
    "security_uuid" UUID,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "calendar" TEXT,
    "feed" TEXT,
    "feed_url" TEXT,
    "latest_feed" TEXT,
    "latest_feed_url" TEXT,
    "attributes" JSONB,
    "events" JSONB,
    "properties" JSONB,

    PRIMARY KEY ("portfolio_id", "uuid")
);

CREATE TABLE "portfolios_securities_prices" (
    "date" DATE NOT NULL,
    "value" DECIMAL(16,8) NOT NULL,
    "portfolio_id" INTEGER NOT NULL,
    "portfolio_security_uuid" UUID NOT NULL,

    PRIMARY KEY ("portfolio_id", "portfolio_security_uuid", "date")
);

CREATE TABLE "securities" (
    "uuid" UUID NOT NULL,
    "name" TEXT,
    "isin" VARCHAR(12),
    "wkn" VARCHAR(6),
    "symbol_xfra" VARCHAR(10),
    "symbol_xnas" VARCHAR(10),
    "symbol_xnys" VARCHAR(10),
    "security_type" TEXT,

    PRIMARY KEY ("uuid")
);

CREATE TABLE "securities_markets" (
    "id" SERIAL NOT NULL,
    "security_uuid" UUID NOT NULL,
    "market_code" TEXT NOT NULL,
    "currency_code" CHAR(3) NOT NULL,
    "first_price_date" DATE,
    "last_price_date" DATE,
    "symbol" VARCHAR(10),
    "update_prices" BOOLEAN NOT NULL,

    PRIMARY KEY ("id")
);

CREATE TABLE "securities_markets_prices" (
    "security_market_id" INTEGER NOT NULL,
    "date" DATE NOT NULL,
    "close" DECIMAL(10,4) NOT NULL,

    PRIMARY KEY ("security_market_id","date")
);

CREATE TABLE "securities_taxonomies" (
    "security_uuid" UUID NOT NULL,
    "taxonomy_uuid" UUID NOT NULL,
    "weight" DECIMAL(5,2) NOT NULL,

    PRIMARY KEY ("taxonomy_uuid","security_uuid")
);

CREATE TABLE "sessions" (
    "token" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "last_activity_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "note" TEXT NOT NULL DEFAULT E'',
    "user_id" INTEGER NOT NULL,

    PRIMARY KEY ("token")
);

CREATE TABLE "taxonomies" (
    "uuid" UUID NOT NULL,
    "parent_uuid" UUID,
    "root_uuid" UUID,
    "name" TEXT NOT NULL,
    "code" TEXT,

    PRIMARY KEY ("uuid")
);

CREATE TABLE "portfolios_transactions" (
    "type" "portfolio_transaction_type" NOT NULL,
    "datetime" TIMESTAMPTZ NOT NULL,
    "shares" DECIMAL(16,8),
    "note" VARCHAR NOT NULL DEFAULT E'',
    "portfolio_id" INTEGER NOT NULL,
    "uuid" UUID NOT NULL,
    "account_uuid" UUID NOT NULL,
    "partner_transaction_uuid" UUID,
    "portfolio_security_uuid" UUID,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ("portfolio_id", "uuid")
);

CREATE TABLE "portfolios_transactions_units" (
    "id" SERIAL NOT NULL,
    "type" "portfolio_transaction_unit_type" NOT NULL,
    "amount" DECIMAL(10,2) NOT NULL,
    "currency_code" CHAR(3) NOT NULL,
    "original_amount" DECIMAL(10,2),
    "original_currency_code" CHAR(3),
    "exchange_rate" DECIMAL(16,8),
    "portfolio_id" INTEGER NOT NULL,
    "transaction_uuid" UUID NOT NULL,

    PRIMARY KEY ("id")
);

CREATE TABLE "users" (
    "id" SERIAL NOT NULL,
    "username" VARCHAR NOT NULL,
    "password" VARCHAR,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "last_seen_at" DATE NOT NULL DEFAULT CURRENT_DATE,
    "is_admin" BOOLEAN NOT NULL DEFAULT false,

    PRIMARY KEY ("id")
);

-- Create Indexes
CREATE INDEX "portfolios_accounts.portfolio_id_index" ON "portfolios_accounts"("portfolio_id");
CREATE INDEX "clientupdates_country" ON "clientupdates"("country");
CREATE INDEX "clientupdates_timestamp" ON "clientupdates"("timestamp");
CREATE INDEX "clientupdates_version" ON "clientupdates"("version");
CREATE INDEX "events_security_uuid" ON "events"("security_uuid");
CREATE UNIQUE INDEX "exchangerates.base_currency_code_quote_currency_code_unique" ON "exchangerates"("base_currency_code", "quote_currency_code");
CREATE INDEX "exchangerates_prices.exchangerate_id_index" ON "exchangerates_prices"("exchangerate_id");
CREATE INDEX "portfolios.user_id_index" ON "portfolios"("user_id");
CREATE INDEX "portfolios_securities.portfolio_id_index" ON "portfolios_securities"("portfolio_id");
CREATE INDEX "portfolios_securities_prices.portfolio_id_portfolio_security_uuid_index" ON "portfolios_securities_prices"("portfolio_id", "portfolio_security_uuid");
CREATE INDEX "securities_isin" ON "securities"("isin");
CREATE INDEX "securities_name" ON "securities"("name");
CREATE INDEX "securities_name_trigram" ON "securities" USING gin("name" gin_trgm_ops);
CREATE INDEX "securities_security_type" ON "securities"("security_type");
CREATE INDEX "securities_symbol_xfra" ON "securities"("symbol_xfra");
CREATE INDEX "securities_symbol_xnas" ON "securities"("symbol_xnas");
CREATE INDEX "securities_symbol_xnys" ON "securities"("symbol_xnys");
CREATE INDEX "securities_wkn" ON "securities"("wkn");
CREATE UNIQUE INDEX "securities_markets_security_uuid_market_code" ON "securities_markets"("security_uuid", "market_code");
CREATE INDEX "securities_markets_security_uuid" ON "securities_markets"("security_uuid");
CREATE INDEX "securities_markets_prices.security_market_id_index" ON "securities_markets_prices"("security_market_id");
CREATE INDEX "sessions.user_id_index" ON "sessions"("user_id");
CREATE INDEX "sessions.last_activity_at_index" ON "sessions"("last_activity_at");
CREATE UNIQUE INDEX "taxonomies_root_uuid_code" ON "taxonomies"("root_uuid", "code");
CREATE INDEX "taxonomies_parent_uuid" ON "taxonomies"("parent_uuid");
CREATE INDEX "taxonomies_root_uuid" ON "taxonomies"("root_uuid");
CREATE UNIQUE INDEX "portfolios_transactions.portfolio_id_partner_transaction_uuid_unique" ON "portfolios_transactions"("portfolio_id", "partner_transaction_uuid");
CREATE INDEX "portfolios_transactions.portfolio_id_account_uuid_index" ON "portfolios_transactions"("portfolio_id", "account_uuid");
CREATE INDEX "portfolios_transactions.portfolio_id_portfolio_security_uuid_index" ON "portfolios_transactions"("portfolio_id", "portfolio_security_uuid");
CREATE INDEX "portfolios_transactions.portfolio_id_index" ON "portfolios_transactions"("portfolio_id");
CREATE INDEX "portfolios_transactions_units.portfolio_id_transaction_uuid_index" ON "portfolios_transactions_units"("portfolio_id", "transaction_uuid");
CREATE UNIQUE INDEX "users.username_unique" ON "users"("username");

-- Add Foreign Keys
ALTER TABLE "portfolios_accounts" ADD FOREIGN KEY ("currency_code") REFERENCES "currencies"("code") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "portfolios_accounts" ADD FOREIGN KEY ("portfolio_id") REFERENCES "portfolios"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "portfolios_accounts" ADD FOREIGN KEY ("portfolio_id", "reference_account_uuid") REFERENCES "portfolios_accounts"("portfolio_id", "uuid") ON UPDATE CASCADE;
ALTER TABLE "events" ADD FOREIGN KEY ("currency_code") REFERENCES "currencies"("code") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "events" ADD FOREIGN KEY ("security_uuid") REFERENCES "securities"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "exchangerates" ADD FOREIGN KEY ("base_currency_code") REFERENCES "currencies"("code") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "exchangerates" ADD FOREIGN KEY ("quote_currency_code") REFERENCES "currencies"("code") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "exchangerates_prices" ADD FOREIGN KEY ("exchangerate_id") REFERENCES "exchangerates"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "portfolios" ADD FOREIGN KEY ("base_currency_code") REFERENCES "currencies"("code") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "portfolios" ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "portfolios_securities" ADD FOREIGN KEY ("currency_code") REFERENCES "currencies"("code") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "portfolios_securities" ADD FOREIGN KEY ("portfolio_id") REFERENCES "portfolios"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "portfolios_securities" ADD FOREIGN KEY ("security_uuid") REFERENCES "securities"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "portfolios_securities_prices" ADD FOREIGN KEY ("portfolio_id", "portfolio_security_uuid") REFERENCES "portfolios_securities"("portfolio_id", "uuid") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "securities_markets" ADD FOREIGN KEY ("security_uuid") REFERENCES "securities"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "securities_markets" ADD FOREIGN KEY ("market_code") REFERENCES "markets"("code") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "securities_markets" ADD FOREIGN KEY ("currency_code") REFERENCES "currencies"("code") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "securities_markets_prices" ADD FOREIGN KEY ("security_market_id") REFERENCES "securities_markets"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "securities_taxonomies" ADD FOREIGN KEY ("security_uuid") REFERENCES "securities"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "securities_taxonomies" ADD FOREIGN KEY ("taxonomy_uuid") REFERENCES "taxonomies"("uuid") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "taxonomies" ADD FOREIGN KEY ("parent_uuid") REFERENCES "taxonomies"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "taxonomies" ADD FOREIGN KEY ("root_uuid") REFERENCES "taxonomies"("uuid") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "portfolios_transactions" ADD FOREIGN KEY ("portfolio_id", "account_uuid") REFERENCES "portfolios_accounts"("portfolio_id", "uuid") ON UPDATE CASCADE;
ALTER TABLE "portfolios_transactions" ADD FOREIGN KEY ("portfolio_id", "partner_transaction_uuid") REFERENCES "portfolios_transactions"("portfolio_id", "uuid") ON UPDATE CASCADE;
ALTER TABLE "portfolios_transactions" ADD FOREIGN KEY ("portfolio_id") REFERENCES "portfolios"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "portfolios_transactions" ADD FOREIGN KEY ("portfolio_id", "portfolio_security_uuid") REFERENCES "portfolios_securities"("portfolio_id", "uuid") ON UPDATE CASCADE;
ALTER TABLE "portfolios_transactions_units" ADD FOREIGN KEY ("currency_code") REFERENCES "currencies"("code") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "portfolios_transactions_units" ADD FOREIGN KEY ("original_currency_code") REFERENCES "currencies"("code") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "portfolios_transactions_units" ADD FOREIGN KEY ("portfolio_id", "transaction_uuid") REFERENCES "portfolios_transactions"("portfolio_id", "uuid") ON DELETE CASCADE ON UPDATE CASCADE;

-- Insert Data
INSERT INTO
  currencies (code)
VALUES
  ('AUD'),
  ('BGN'),
  ('BRL'),
  ('CAD'),
  ('CHF'),
  ('CNY'),
  ('CZK'),
  ('DKK'),
  ('EUR'),
  ('GBP'),
  ('HKD'),
  ('HRK'),
  ('HUF'),
  ('IDR'),
  ('ILS'),
  ('INR'),
  ('ISK'),
  ('JPY'),
  ('KRW'),
  ('MXN'),
  ('MYR'),
  ('NOK'),
  ('NZD'),
  ('PHP'),
  ('PLN'),
  ('RON'),
  ('RUB'),
  ('SEK'),
  ('SGD'),
  ('THB'),
  ('TRY'),
  ('USD'),
  ('ZAR'),
  ('AED'),
  ('GBX');

INSERT INTO
  exchangerates (base_currency_code, quote_currency_code)
VALUES
  ('EUR', 'USD'),
  ('EUR', 'AUD'),
  ('EUR', 'BGN'),
  ('EUR', 'BRL'),
  ('EUR', 'CAD'),
  ('EUR', 'CHF'),
  ('EUR', 'CNY'),
  ('EUR', 'CZK'),
  ('EUR', 'DKK'),
  ('EUR', 'GBP'),
  ('EUR', 'HKD'),
  ('EUR', 'HRK'),
  ('EUR', 'HUF'),
  ('EUR', 'IDR'),
  ('EUR', 'ILS'),
  ('EUR', 'INR'),
  ('EUR', 'ISK'),
  ('EUR', 'JPY'),
  ('EUR', 'KRW'),
  ('EUR', 'MXN'),
  ('EUR', 'MYR'),
  ('EUR', 'NOK'),
  ('EUR', 'NZD'),
  ('EUR', 'PHP'),
  ('EUR', 'PLN'),
  ('EUR', 'RON'),
  ('EUR', 'RUB'),
  ('EUR', 'SEK'),
  ('EUR', 'SGD'),
  ('EUR', 'THB'),
  ('EUR', 'TRY'),
  ('EUR', 'ZAR'),
  ('GBP', 'GBX'),
  ('USD', 'AED');

INSERT INTO
  exchangerates_prices (date, value, exchangerate_id)
SELECT
  x.date::date, x.value::decimal, er.id
FROM
  ( VALUES
      ('1999-01-04', '1.1789', 'EUR', 'USD'),
      ('1999-01-05', '1.1790', 'EUR', 'USD'),
      ('1999-01-06', '1.1743', 'EUR', 'USD'),
      ('1999-01-07', '1.1632', 'EUR', 'USD'),
      ('1999-01-08', '1.1659', 'EUR', 'USD'),
      ('1999-01-11', '1.1569', 'EUR', 'USD'),
      ('1999-01-12', '1.1520', 'EUR', 'USD'),
      ('1999-01-13', '1.1744', 'EUR', 'USD'),
      ('1999-01-14', '1.1653', 'EUR', 'USD'),
      ('1999-01-15', '1.1626', 'EUR', 'USD'),
      ('1999-01-18', '1.1612', 'EUR', 'USD'),
      ('1999-01-19', '1.1616', 'EUR', 'USD'),
      ('1999-01-20', '1.1575', 'EUR', 'USD'),
      ('1999-01-21', '1.1572', 'EUR', 'USD'),
      ('1999-01-22', '1.1567', 'EUR', 'USD'),
      ('1999-01-25', '1.1584', 'EUR', 'USD'),
      ('1999-01-26', '1.1582', 'EUR', 'USD'),
      ('1999-01-27', '1.1529', 'EUR', 'USD'),
      ('1999-01-28', '1.1410', 'EUR', 'USD'),
      ('1999-01-29', '1.1384', 'EUR', 'USD'),
      ('1971-01-01', '100.00', 'GBP', 'GBX'),
      ('1990-01-01', '3.6725', 'USD', 'AED')
  ) x (date, value, base_currency_code, quote_currency_code)
  JOIN exchangerates er ON (
    er.base_currency_code = x.base_currency_code
    AND er.quote_currency_code = x.quote_currency_code
  );
