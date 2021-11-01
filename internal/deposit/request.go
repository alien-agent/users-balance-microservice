package deposit

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Request represents a JSON data of an API request.
type Request interface {
	// Validate validates the request's fields.
	Validate() error
}

// GetBalanceRequest represents a request to get balance of specific user.
type GetBalanceRequest struct {
	OwnerId string `json:"owner_id"`
}

// Validate validates the GetBalanceRequest fields.
func (r GetBalanceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerId, validation.Required, is.UUID),
	)
}

// UpdateBalanceRequest represents a request to update user's balance.
type UpdateBalanceRequest struct {
	OwnerId     string `json:"owner_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description,omitempty"`
}

func (r UpdateBalanceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerId, validation.Required, is.UUID),
		validation.Field(&r.Amount, validation.Required),
		validation.Field(&r.Description, validation.Length(0, 100)),
	)
}

// TransferRequest represents a request to transfer money from one user to another.
type TransferRequest struct {
	SenderId    string `json:"sender_id"`
	RecipientId string `json:"recipient_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

// Validate validates the TransferRequest fields.
func (r TransferRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SenderId, validation.Required, is.UUID),
		validation.Field(&r.RecipientId, validation.Required, is.UUID),
		validation.Field(&r.Amount, validation.Required, validation.Min(0).Exclusive()),
		validation.Field(&r.Description, validation.Length(0, 100)),
	)
}

// GetHistoryRequest represents a request to get a list of all user's transactions: top-ups, withdrawals and transfers.
type GetHistoryRequest struct {
	OwnerId        string `json:"owner_id"`
	Offset         int    `json:"offset,omitempty"`
	Limit          int    `json:"limit,omitempty"`
	OrderBy        string `json:"order_by,omitempty"`
	OrderDirection string `json:"order_direction,omitempty"`
}

// Validate validates the GetHistoryRequest.
func (r GetHistoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OwnerId, validation.Required, is.UUID),
		validation.Field(&r.Offset, validation.Min(0)),
		validation.Field(&r.Limit, validation.Min(1)),
		validation.Field(&r.OrderBy, validation.In("transaction_date", "amount")),
		validation.Field(&r.OrderDirection, validation.In("ASC", "DESC")),
	)
}
