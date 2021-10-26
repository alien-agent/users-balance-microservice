package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id              int64     `json:"id,omitempty" db:"pk"`
	SenderId        uuid.UUID    `json:"sender_id,omitempty"`
	RecipientId     uuid.UUID    `json:"recipient_id,omitempty"`
	Amount          int64     `json:"amount"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date,omitempty"`
}
