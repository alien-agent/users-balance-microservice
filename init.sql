CREATE TABLE IF NOT EXISTS Deposit(
    owner_id UUID PRIMARY KEY,
    balance INT NOT NULL,

    CONSTRAINT chk_balance_not_negative
    CHECK(balance >= 0) /* super-safe :) */
);

CREATE TABLE IF NOT EXISTS Transaction(
    id serial PRIMARY KEY,
    sender_id UUID NULL,
    recipient_id UUID NULL,
    amount INT NOT NULL,
    description VARCHAR(100) NULL,
    transaction_date TIMESTAMP NOT NULL,

    CONSTRAINT chk_amount_not_negative
    CHECK(amount >= 0)
);