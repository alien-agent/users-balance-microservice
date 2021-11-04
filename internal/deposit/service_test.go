package deposit

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/internal/requests"
	"users-balance-microservice/pkg/log"
)

var (
	databaseError   = errors.New("database error")
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
		}, exchangeService, logger,
	)

	// initial count
	count, err := s.Count(ctx)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 1, count)
	}

	// get existing deposit's balance in RUB
	balance, err := s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id1.String()})
	if assert.NoError(t, err) {
		assert.EqualValues(t, 1000, balance)
	}

	// get existing deposit's balance in USD (fake exchange rate RUB/USD=0.1 is used)
	balance, err = s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id1.String(), Currency: "USD"})
	if assert.NoError(t, err) {
		assert.EqualValues(t, 100, balance)
	}

	// get non-existing deposit's balance - 0 is returned regardless of currency, new deposit is not created.
	balance, err = s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id2.String(), Currency: "EUR"})
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
		}, exchangeService, logger,
	)

	// initial count
	count, err := s.Count(ctx)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 1, count)
	}

	// update balance top-up success
	err = s.Update(ctx, requests.UpdateBalanceRequest{OwnerId: id1.String(), Amount: 500, Description: "visa top-up"})
	if assert.NoError(t, err) {
		balance, err := s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1500, balance)
		}
	}

	// update balance withdrawal success
	err = s.Update(ctx, requests.UpdateBalanceRequest{
		OwnerId:     id1.String(),
		Amount:      -500,
		Description: "monthly subscription",
	})
	if assert.NoError(t, err) {
		balance, err := s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1000, balance)
		}
	}

	// update balance top-up non-existing deposit -> new deposit is created
	err = s.Update(ctx, requests.UpdateBalanceRequest{
		OwnerId:     id2.String(),
		Amount:      2000,
		Description: "mastercard top-up",
	})
	if assert.NoError(t, err) {
		count, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 2, count)
		}

		balance, err := s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id2.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 2000, balance)
		}
	}

	// update balance withdrawal insufficient balance -> failure
	err = s.Update(ctx, requests.UpdateBalanceRequest{
		OwnerId:     id1.String(),
		Amount:      -250000,
		Description: "hack3r attack",
	})
	if assert.Error(t, err) {
		balance, err := s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1000, balance)
		}
	}

	// update balance invalid owner_id -> failure
	count, _ = s.Count(ctx)
	err = s.Update(ctx, requests.UpdateBalanceRequest{OwnerId: "123-456-789", Amount: 2000})
	if assert.Error(t, err) {
		count2, _ := s.Count(ctx)
		assert.EqualValues(t, 0, count2-count)
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
		}, exchangeService, logger,
	)

	// transfer success
	err := s.Transfer(ctx, requests.TransferRequest{
		SenderId:    id2.String(),
		RecipientId: id1.String(),
		Amount:      300,
		Description: "thanks for dinner!",
	})
	if assert.NoError(t, err) {
		balance, err := s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1300, balance)
		}

		balance, err = s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id2.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1700, balance)
		}
	}

	// transfer from existing to non-existing deposit success
	err = s.Transfer(ctx, requests.TransferRequest{
		SenderId:    id2.String(),
		RecipientId: id3.String(),
		Amount:      700,
		Description: "happy birthday!",
	})
	if assert.NoError(t, err) {
		balance, err := s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id2.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1000, balance)
		}

		balance, err = s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id3.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 700, balance)
		}
	}

	// transfer insufficient funds failure
	err = s.Transfer(ctx, requests.TransferRequest{
		SenderId:    id2.String(),
		RecipientId: id1.String(),
		Amount:      300000,
		Description: "thanks for dinner!",
	})
	if assert.Error(t, err) {
		balance, err := s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1300, balance)
		}

		balance, err = s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id2.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1000, balance)
		}
	}

	// transfer negative amount failure
	err = s.Transfer(ctx, requests.TransferRequest{
		SenderId:    id2.String(),
		RecipientId: id1.String(),
		Amount:      -300,
		Description: "hacker attack",
	})
	if assert.Error(t, err) {
		balance, err := s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id1.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1300, balance)
		}

		balance, err = s.GetBalance(ctx, requests.GetBalanceRequest{OwnerId: id2.String()})
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1000, balance)
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
	// simulate database error
	if deposit.OwnerId.String() == "11111111-1111-1111-1111-111111111111" {
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

// Fake exchange rates service provides exchange ratio=0.1 regardless of currency code.
type mockExchangeRatesService struct{}

func (s mockExchangeRatesService) Get(code string) (float32, error) {
	return 0.1, nil
}
