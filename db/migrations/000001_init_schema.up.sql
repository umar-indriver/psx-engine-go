-- Database Migration: Init PSX Schema

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    tax_status VARCHAR(20) NOT NULL DEFAULT 'FILER',
    subscription_status VARCHAR(20) NOT NULL DEFAULT 'FREE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS portfolios (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'PKR',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_portfolios_user_id ON portfolios(user_id);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    portfolio_id UUID NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    symbol VARCHAR(15) NOT NULL,
    trade_type VARCHAR(4) NOT NULL CHECK (trade_type IN ('BUY', 'SELL')),
    quantity NUMERIC(15, 4) NOT NULL,
    price_per_share NUMERIC(15, 4) NOT NULL,
    commission NUMERIC(15, 4) NOT NULL DEFAULT 0.0000,
    trade_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_transactions_portfolio_symbol ON transactions(portfolio_id, symbol);

CREATE TABLE IF NOT EXISTS dividend_announcements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    symbol VARCHAR(15) NOT NULL,
    payout_per_share NUMERIC(15, 4) NOT NULL,
    book_closure_date DATE NOT NULL,
    payment_date DATE,
    financial_period VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_dividends_symbol_date ON dividend_announcements(symbol, book_closure_date DESC);
