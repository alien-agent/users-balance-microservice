package entity

import "github.com/google/uuid"

type Deposit struct {
	OwnerId uuid.UUID `json:"owner_id" db:"pk"`
	Balance int64     `json:"balance"`
}
