package deposit

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type validationTestcase struct {
	name      string
	model     Request
	wantError bool
}

func testValidation(t *testing.T, tests []validationTestcase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestGetBalanceRequest_Validate(t *testing.T) {
	id1 := uuid.NewString()
	testValidation(t, []validationTestcase{
		{"success", GetBalanceRequest{OwnerId: id1}, false},
		{"success with currency", GetBalanceRequest{OwnerId: id1, Currency: "EUR"}, false},
		{"missing OwnerId", GetBalanceRequest{OwnerId: ""}, true},
		{"invalid OwnerId", GetBalanceRequest{OwnerId: "12712912"}, true},
		{"fail invalid currency", GetBalanceRequest{OwnerId: id1, Currency: "EURUSDPLT"}, true},
	})
}

func TestUpdateBalanceRequest_Validate(t *testing.T) {
	id1 := uuid.NewString()
	testValidation(t, []validationTestcase{
		{"success positive amount", UpdateBalanceRequest{id1, 500, "visa"}, false},
		{"success negative amount", UpdateBalanceRequest{id1, -500, "mastercard"}, false},
		{"success no description", UpdateBalanceRequest{id1, 500, ""}, false},
		{"fail zero amount", UpdateBalanceRequest{id1, 0, ""}, true},
		{"fail invalid OwnerId", UpdateBalanceRequest{"i'm invalid", 500, ""}, true},
		{"fail too long description", UpdateBalanceRequest{id1, 500, strings.Repeat("test", 100)}, true},
	})
}

func TestTransferRequest_Validate(t *testing.T) {
	id1, id2 := uuid.NewString(), uuid.NewString()
	testValidation(t, []validationTestcase{
		{"success no description", TransferRequest{id1, id2, 500, ""}, false},
		{"success with description", TransferRequest{id1, id2, 500, "thanks for dinner"}, false},
		{"fail negative amount", TransferRequest{id1, id2, -500, ""}, true},
		{"fail missing missing SenderId", TransferRequest{"", id2, 500, ""}, true},
		{"fail missing missing RecipientId", TransferRequest{id1, "", 500, ""}, true},
		{"fail SenderId invalid", TransferRequest{"124124-12412-12412", id2, 500, ""}, true},
		{"fail RecipientId invalid", TransferRequest{id1, "982312-124-124-43", 500, ""}, true},
		{"fail description too long", TransferRequest{id1, id2, 500, strings.Repeat("test", 100)}, true},
	})
}

func TestGetHistoryRequest_Validate(t *testing.T) {
	id1 := uuid.NewString()
	testValidation(t, []validationTestcase{
		{"success only OwnerId", GetHistoryRequest{OwnerId: id1}, false},
		{"success with ordering", GetHistoryRequest{OwnerId: id1, OrderBy: "amount", OrderDirection: "ASC"}, false},
		{"success with limit&offset", GetHistoryRequest{OwnerId: id1, Offset: 10, Limit: 5}, false},
		{"success all params", GetHistoryRequest{id1, 10, 5, "transaction_date", "DESC"}, false},
		{"fail missing OwnerId", GetHistoryRequest{OwnerId: ""}, true},
		{"fail invalid OwnerId", GetHistoryRequest{OwnerId: "128312-1241-12"}, true},
		{"fail invalid OrderBy", GetHistoryRequest{OwnerId: id1, OrderBy: "owner_id"}, true},
		{
			"fail invalid OrderDirection",
			GetHistoryRequest{OwnerId: id1, OrderBy: "amount", OrderDirection: "MEDIAN"},
			true,
		},
		{"fail negative offset", GetHistoryRequest{OwnerId: id1, Offset: -10}, true},
		{"fail negative limit", GetHistoryRequest{OwnerId: id1, Limit: -5}, true},
	})
}
