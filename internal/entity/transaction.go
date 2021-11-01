package entity

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a single change in user's Deposit.
//
// SenderId and RecipientId are "positional". The transaction Amount (which is positive) is always subtracted from
// sender's deposit and added to recipient's deposit.
//
// If a transaction is missing a RecipientId, it is considered a deposit withdrawal.
// If a transaction is missing a SenderId, it is considered a deposit top-up.
// Otherwise, a transaction is considered a money transfer between two users within the system.
type Transaction struct {
	// Database id of this Transaction.
	Id int64 `json:"id,omitempty" db:"pk"`
	// UUID of sender's Deposit. Optional.
	SenderId uuid.UUID `json:"sender_id,omitempty"`
	// UUID of recipient's Deposit. Optional.
	RecipientId uuid.UUID `json:"recipient_id,omitempty"`
	// An amount of rubles subtracted from sender's deposit and added to recipient's deposit. Positive.
	Amount int64 `json:"amount"`
	// The description of this Transaction. Optional.
	Description string `json:"description"`
	// The date and time when this Transaction was made.
	TransactionDate time.Time `json:"transaction_date,omitempty"`
}
