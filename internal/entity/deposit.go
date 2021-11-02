package entity

import "github.com/google/uuid"

// Deposit represents a user's account in the database.
type Deposit struct {
	// OwnerId is a UUID of the user which this Deposit belongs to. Serves as primary key in the database.
	OwnerId uuid.UUID `json:"owner_id" db:"pk"`
	// Balance is an amount of money which is available to this user. Non-negative.
	Balance int64 `json:"balance"`
}
