package deposit

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/pkg/log"
)

var (
	databaseError = errors.New("database error")

	logger, _       = log.NewForTest()
	exchangeService = mockExchangeRatesService{}
	ctx             = context.Background()
)

func TestService_GetBalance(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	s := NewService(
		&mockDepositRepository{
			items: []entity.Deposit{
				{id1, 1000},
			},
		}, &mockTransactionRepository{}, exchangeService, logger,
	)

	// initial count
	count, err := s.Count(ctx)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 1, count)
	}

	// get existing deposit's balance in RUB
	balance, err := s.GetBalance(ctx, GetBalanceRequest{OwnerId: id1.String()})
	if assert.NoError(t, err) {
		assert.EqualValues(t, 1000, balance)
	}

	// get existing deposit's balance in USD (fake exchange rate RUB/USD=2 is used)
	balance, err = s.GetBalance(ctx, GetBalanceRequest{OwnerId: id1.String(), Currency: "USD"})
	if assert.NoError(t, err) {
		assert.EqualValues(t, 2000, balance)
	}

	// get non-existing deposit's balance - 0 is returned regardless of currency, new deposit is not created.
	balance, err = s.GetBalance(ctx, GetBalanceRequest{OwnerId: id2.String(), Currency: "EUR"})
	if assert.NoError(t, err) {
		assert.EqualValues(t, 0, balance)
	}
}

func TestService_Update(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	s := NewService(
		&mockDepositRepository{
			items: []entity.Deposit{
				{id1, 1000},
			},
		}, &mockTransactionRepository{}, exchangeService, logger,
	)

	// initial count
	count, err := s.Count(ctx)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 1, count)
	}

	// update balance top-up success
	tx, err := s.Update(ctx, UpdateBalanceRequest{id1.String(), 500, "visa top-up"})
	if assert.NoError(t, err) {
		// verify transaction
		assert.EqualValues(t, uuid.Nil, tx.SenderId)
		assert.EqualValues(t, id1, tx.RecipientId)
		assert.EqualValues(t, 500, tx.Amount)
		assert.Equal(t, "visa top-up", tx.Description)

		balance, err := s.GetBalance(ctx, GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1500, balance)
		}
	}

	// update balance withdrawal success
	tx, err = s.Update(ctx, UpdateBalanceRequest{id1.String(), -500, "monthly subscription"})
	if assert.NoError(t, err) {
		// verify transaction
		assert.EqualValues(t, id1, tx.SenderId)
		assert.EqualValues(t, uuid.Nil, tx.RecipientId)
		assert.EqualValues(t, 500, tx.Amount)
		assert.Equal(t, "monthly subscription", tx.Description)

		balance, err := s.GetBalance(ctx, GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1000, balance)
		}
	}

	// update balance withdrawal insufficient balance -> failure
	_, err = s.Update(ctx, UpdateBalanceRequest{id1.String(), -250000, "hacker attack"})
	if assert.Error(t, err) {
		balance, err := s.GetBalance(ctx, GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1000, balance)
		}
	}

	// update balance top-up non-existing deposit -> new deposit is created
	tx, err = s.Update(ctx, UpdateBalanceRequest{id2.String(), 2000, "mastercard top-up"})
	if assert.NoError(t, err) {
		// verify transaction
		assert.EqualValues(t, uuid.Nil, tx.SenderId)
		assert.EqualValues(t, id2, tx.RecipientId)
		assert.EqualValues(t, 2000, tx.Amount)
		assert.Equal(t, "mastercard top-up", tx.Description)

		count, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 2, count)
		}

		balance, err := s.GetBalance(ctx, GetBalanceRequest{OwnerId: id2.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 2000, balance)
		}
	}
}

func TestService_Transfer(t *testing.T) {
	id1, id2, id3 := uuid.New(), uuid.New(), uuid.New()
	s := NewService(
		&mockDepositRepository{
			items: []entity.Deposit{
				{id1, 1000},
				{id2, 2000},
			},
		}, &mockTransactionRepository{}, exchangeService, logger,
	)

	// transfer success
	tx, err := s.Transfer(ctx, TransferRequest{id2.String(), id1.String(), 300, "thanks for dinner!"})
	if assert.NoError(t, err) {
		// verify transaction
		assert.EqualValues(t, id2, tx.SenderId)
		assert.EqualValues(t, id1, tx.RecipientId)
		assert.EqualValues(t, 300, tx.Amount)
		assert.Equal(t, "thanks for dinner!", tx.Description)

		balance, err := s.GetBalance(ctx, GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1300, balance)
		}

		balance, err = s.GetBalance(ctx, GetBalanceRequest{OwnerId: id2.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1700, balance)
		}
	}

	// transfer insufficient funds failure
	_, err = s.Transfer(ctx, TransferRequest{id2.String(), id1.String(), 300000, "thanks for dinner!"})
	if assert.Error(t, err) {
		balance, err := s.GetBalance(ctx, GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1300, balance)
		}

		balance, err = s.GetBalance(ctx, GetBalanceRequest{OwnerId: id2.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1700, balance)
		}
	}

	// transfer from existing to non-existing deposit success
	// transfer success
	tx, err = s.Transfer(ctx, TransferRequest{id2.String(), id3.String(), 700, "happy birthday!"})
	if assert.NoError(t, err) {
		// verify transaction
		assert.EqualValues(t, id2, tx.SenderId)
		assert.EqualValues(t, id3, tx.RecipientId)
		assert.EqualValues(t, 700, tx.Amount)
		assert.Equal(t, "happy birthday!", tx.Description)

		balance, err := s.GetBalance(ctx, GetBalanceRequest{OwnerId: id2.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1000, balance)
		}

		balance, err = s.GetBalance(ctx, GetBalanceRequest{OwnerId: id3.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 700, balance)
		}
	}
}

type mockDepositRepository struct {
	items []entity.Deposit
}

func (m *mockDepositRepository) Get(ctx context.Context, ownerId uuid.UUID) (entity.Deposit, error) {
	for _, item := range m.items {
		if item.OwnerId == ownerId {
			return item, nil
		}
	}
	return entity.Deposit{}, sql.ErrNoRows
}

func (m *mockDepositRepository) Create(ctx context.Context, deposit entity.Deposit) error {
	if deposit.Balance < 0 {
		return databaseError
	}
	m.items = append(m.items, deposit)
	return nil
}

func (m *mockDepositRepository) Update(ctx context.Context, deposit entity.Deposit) error {
	if deposit.Balance < 0 {
		return databaseError
	}

	for i, item := range m.items {
		if item.OwnerId == deposit.OwnerId {
			m.items[i] = deposit
			return nil
		}
	}

	return m.Create(ctx, deposit)
}

func (m *mockDepositRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
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

// Fake exchange rates service provides exchange ratio=2 regardless of currency code.
type mockExchangeRatesService struct{}

func (s mockExchangeRatesService) Get(code string) (float32, error) {
	return 2, nil
}
