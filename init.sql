CREATE TABLE IF NOT EXISTS Deposit(
    owner_id UUID PRIMARY KEY,
    balance INT NOT NULL
);

CREATE TABLE IF NOT EXISTS Transaction(
    id serial PRIMARY KEY,
    sender_id UUID,
    recipient_id UUID,
    amount INT NOT NULL,
    description VARCHAR(100),
    transaction_date TIMESTAMP NOT NULL,
    CONSTRAINT fk_senderid
        FOREIGN KEY (sender_id)
            REFERENCES Deposit(owner_id),
    CONSTRAINT fk_recipientid
        FOREIGN KEY (recipient_id)
            REFERENCES Deposit(owner_id)
);