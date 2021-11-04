/*
For testing purposes, database failure is simulated if Deposit.OwnerId == "11111111-1111-1111-1111-111111111111" or
if Transaction.SenderId or Transaction.RecipientId are equal to "22222222-2222-2222-2222-222222222222".
*/
package deposit

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/internal/exchangerates"
	"users-balance-microservice/internal/test"
	"users-balance-microservice/internal/transaction"
	"users-balance-microservice/pkg/dbcontext"
	"users-balance-microservice/pkg/log"
)

const invalidIdResponse = `{"status":400,"message":"There is some problem with the data you submitted.","details":[{"field":"owner_id","error":"must be a valid UUID"}]}`
const badRequestResponse = `{"status":400,"message":"Your request is in a bad format."}`

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	depositRepo := &mockDepositRepository{
		items: []entity.Deposit{
			{uuid.MustParse("615f3e76-37d3-11ec-8d3d-0242ac130003"), 1000},
		},
	}
	transactionRepo := mockTransactionRepository{
		items: []entity.Transaction{},
	}
	exchangeService := exchangerates.NewService(30*time.Minute, logger)
	RegisterHandlers(
		router.Group(""),
		NewService(depositRepo, exchangeService, logger),
		transaction.NewService(&transactionRepo, logger),
		logger,
		dbcontext.New(nil),
	)

	tests := []test.APITestCase{
		{
			"get balance existing",
			"POST",
			"/deposits/balance",
			`{"owner_id": "615f3e76-37d3-11ec-8d3d-0242ac130003"}`,
			http.StatusOK,
			`1000`,
		},
		{
			"get balance non-existing",
			"POST",
			"/deposits/balance",
			`{"owner_id": "8c5593a0-37d3-11ec-8d3d-0242ac130003"}`,
			http.StatusOK,
			`0`,
		},
		{
			"get balance invalid owner_id",
			"POST",
			"/deposits/balance",
			`{"owner_id": "0123456789"}`,
			http.StatusBadRequest,
			invalidIdResponse,
		},
		{
			"get balance invalid request",
			"POST",
			"/deposits/balance",
			`{"owner_id": `,
			http.StatusBadRequest,
			badRequestResponse,
		},
		{
			"get balance invalid method",
			"GET",
			"/deposits/balance",
			"",
			http.StatusMethodNotAllowed,
			"",
		},
		{
			"update balance success",
			"POST",
			"/deposits/update",
			`{"owner_id": "615f3e76-37d3-11ec-8d3d-0242ac130003", "amount": 500}`,
			http.StatusOK,
			"1231",
		},
		{
			"get balance updated",
			"POST",
			"/deposits/balance",
			`{"owner_id": "615f3e76-37d3-11ec-8d3d-0242ac130003"}`,
			http.StatusOK,
			`1500`,
		},
		{
			"update balance success positive amount",
			"POST",
			"/deposits/update",
			`{"owner_id":"615f3e76-37d3-11ec-8d3d-0242ac130003","amount":500,"description":"visa top-up"}`,
			http.StatusOK,
			"",
		},
		{
			"update balance success negative amount",
			"POST",
			"/deposits/update",
			`{"owner_id":"615f3e76-37d3-11ec-8d3d-0242ac130003","amount":-500,"description":"visa top-up"}`,
			http.StatusOK,
			"",
		},
		{
			"update balance fail not enough funds",
			"POST",
			"/deposits/update",
			`{"owner_id": "615f3e76-37d3-11ec-8d3d-0242ac130003", "amount": -55000}`,
			http.StatusBadRequest,
			"",
		},
		{
			// Here the Deposit is updated, but the Transaction won't be created, so
			"update balance database failure",
			"POST",
			"/deposits/update",
			`{"owner_id":"22222222-2222-2222-2222-222222222222","amount":500,"description":"visa top-up"}`,
			http.StatusInternalServerError,
			"",
		},
		{
			"get balance after failure unchanged",
			"POST",
			"/deposits/balance",
			`{"owner_id":"22222222-2222-2222-2222-222222222222"}`,
			http.StatusOK,
			"0",
		},
	}

	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}

type mockTransactionRepository struct {
	items          []entity.Transaction
	lastInsertedId int64
}

func (m *mockTransactionRepository) Create(ctx context.Context, tx *entity.Transaction) error {
	if tx.Amount < 0 {
		return databaseError
	}
	// simulate database failure
	if tx.SenderId.String() == "22222222-2222-2222-2222-222222222222" ||
		tx.RecipientId.String() == "22222222-2222-2222-2222-222222222222" {
		return databaseError
	}

	tx.Id = m.lastInsertedId
	m.lastInsertedId++
	m.items = append(m.items, *tx)
	return nil
}

// Offset, limit and order are ignored for simplicity
func (m *mockTransactionRepository) GetForUser(ctx context.Context, ownerId uuid.UUID, order string, offset, limit int) ([]entity.Transaction, error) {
	var result []entity.Transaction

	for _, tx := range m.items {
		if tx.SenderId == ownerId || tx.RecipientId == ownerId {
			result = append(result, tx)
		}
	}

	return result, nil
}

func (m *mockTransactionRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}
