package deposit

import (
	"context"
	"net/http"
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"
	"users-balance-microservice/internal/entity"
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
	exchangeService := mockExchangeRatesService{}

	db, err := dbx.MustOpen("postgres", "postgres://127.0.0.1/postgres?sslmode=disable&user=postgres&password=postgres")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer func() {
		_ = db.Close()
	}()
	dbCtx := dbcontext.New(db) // in order for TransactionHandler() to work

	RegisterHandlers(
		router.Group(""),
		NewService(depositRepo, exchangeService, logger),
		transaction.NewService(&transactionRepo, logger),
		logger,
		dbCtx.TransactionHandler(),
	)

	tests := []test.APITestCase{
		{
			"get balance success existing Deposit",
			"POST",
			"/deposits/balance",
			`{"owner_id": "615f3e76-37d3-11ec-8d3d-0242ac130003"}`,
			http.StatusOK,
			`1000`,
		},
		{
			"get balance success non-existing Deposit",
			"POST",
			"/deposits/balance",
			`{"owner_id": "8c5593a0-37d3-11ec-8d3d-0242ac130003"}`,
			http.StatusOK,
			`0`,
		},
		{
			"get balance failure invalid owner_id",
			"POST",
			"/deposits/balance",
			`{"owner_id": "0123456789"}`,
			http.StatusBadRequest,
			invalidIdResponse,
		},
		{
			"get balance failure invalid request",
			"POST",
			"/deposits/balance",
			`{"owner_id": `,
			http.StatusBadRequest,
			badRequestResponse,
		},
		{
			"get balance failure invalid method",
			"GET",
			"/deposits/balance",
			"",
			http.StatusMethodNotAllowed,
			"",
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
			"update balance failure not enough funds",
			"POST",
			"/deposits/update",
			`{"owner_id": "615f3e76-37d3-11ec-8d3d-0242ac130003", "amount": -55000}`,
			http.StatusBadRequest,
			"",
		},
		{
			"update balance failure invalid request",
			"POST",
			"/deposits/update",
			`owner_id:"11111111-1111-1111-1111-111111111111"`,
			http.StatusBadRequest,
			"",
		},
		{
			"transfer successful",
			"POST",
			"/deposits/transfer",
			`{"sender_id":"615f3e76-37d3-11ec-8d3d-0242ac130003","recipient_id":"11111111-37d3-11ec-8d3d-0242ac130003","amount":100}`,
			http.StatusOK,
			"",
		},
		{
			"transfer failure sender_id missing",
			"POST",
			"/deposits/transfer",
			`{"recipient_id":"11111111-37d3-11ec-8d3d-0242ac130003","amount":100}`,
			http.StatusBadRequest,
			"",
		},
		{
			"transfer failure invalid request",
			"POST",
			"/deposits/transfer",
			`{recipient_id:11111111-37d3-11ec-8d3d-0242ac130003,"amount":100}`,
			http.StatusBadRequest,
			"",
		},
		{
			"transfer failure insufficient sender's balance",
			"POST",
			"/deposits/transfer",
			`{"sender_id":"615f3e76-37d3-11ec-8d3d-0242ac130003","recipient_id":"11111111-37d3-11ec-8d3d-0242ac130003","amount":1000000}`,
			http.StatusBadRequest,
			"",
		},
		{
			"getHistory success",
			"POST",
			"/deposits/history",
			`{"owner_id":"11112222-3333-4444-5555-666677778888"}`,
			http.StatusOK,
			"",
		},
		{
			"getHistory fail invalid owner_id",
			"POST",
			"/deposits/history",
			`{"owner_id":"123-456-789"}`,
			http.StatusBadRequest,
			"",
		},
		{
			"getHistory fail invalid request",
			"POST",
			"/deposits/history",
			`{owner_id: 123-456-789}`,
			http.StatusBadRequest,
			"",
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

	tx.Id = m.lastInsertedId
	m.lastInsertedId++
	m.items = append(m.items, *tx)
	return nil
}

// Offset, limit and order are ignored for simplicity
func (m *mockTransactionRepository) GetForUser(ctx context.Context, ownerId uuid.UUID, orderBy, orderDirection string, offset, limit int) ([]entity.Transaction, error) {
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
