CREATE TABLE IF NOT EXISTS Deposit(
    owner_id UUID PRIMARY KEY,
    balance INT NOT NULL
);

/*CREATE TABLE IF NOT EXISTS Transactions(
    id serial PRIMARY KEY,
    deposit_id UUID REFERENCES Deposits(UserUUID),
    owner_id UUID NOT NULL,
    amount INT NOT NULL,
    reason VARCHAR(250) NOT NULL,
    partner_uuid UUID,
    transaction_date TIMESTAMP NOT NULL
);*/