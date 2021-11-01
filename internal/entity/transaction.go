package entity

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a money transfer in our bank system.
// If a Transaction is missing a RecipientId, it is considered a withdrawal transaction (payment for internal or external service).
// If a Transaction is missing a SenderId, it is considered a top-up transaction.
// Otherwise, a Transaction is considered a money transfer between two users within the system.
type Transaction struct {
	// Database id of this Transaction.
	Id int64 `json:"id,omitempty" db:"pk"`
	// UUID of sender's Deposit or nil.
	SenderId uuid.UUID `json:"sender_id,omitempty"`
	// UUID of recipient's Deposit or nil.
	RecipientId uuid.UUID `json:"recipient_id,omitempty"`
	// An amount of rubles subtracted from sender's deposit and added to recipient's deposit. Positive.
	Amount int64 `json:"amount"`
	// The description of this Transaction. Optional.
	Description string `json:"description"`
	// The date and time when this Transaction was made.
	TransactionDate time.Time `json:"transaction_date,omitempty"`
}
