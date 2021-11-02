package deposit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/internal/exchangerates"
	"users-balance-microservice/pkg/log"
)

var errCRUD = errors.New("error crud")

func Test_service_CRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	exchangeService := exchangerates.NewRatesService(5*time.Minute, logger)
	s := NewService(&mockDepositRepository{}, &mockTransactionRepository{}, exchangeService, logger)

	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// getBalance non-existing deposit -> new deposit created with balance=0
	id1 := uuid.NewString()
	balance, err := s.GetBalance(ctx, GetBalanceRequest{OwnerId: id1})
	assert.Nil(t, err)
	assert.Equal(t, float32(0), balance)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// getBalance invalid UUID -> no new deposit is created, error returned
	_, err = s.GetBalance(ctx, GetBalanceRequest{OwnerId: "it is invalid UUID"})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// update existing deposit - success
	_, err = s.Update(ctx, UpdateBalanceRequest{id1, 500, ""})
	assert.Nil(t, err)
	balance, err = s.GetBalance(ctx, GetBalanceRequest{OwnerId: id1})
	assert.Nil(t, err)
	assert.Equal(t, float32(500), balance)

	// update non-existing deposit - new deposit is created and its balance is updated
	id2 := uuid.NewString()
	_, err = s.Update(ctx, UpdateBalanceRequest{id2, 1200, ""})
	assert.Nil(t, err)
	balance, err = s.GetBalance(ctx, GetBalanceRequest{OwnerId: id2})
	assert.Nil(t, err)
	assert.Equal(t, float32(1200), balance)

	// update deposit, withdrawal amount more than balance -> error, balance is not modified
	_, err = s.Update(ctx, UpdateBalanceRequest{id1, -2500, ""})
	assert.NotNil(t, err)
	balance, err = s.GetBalance(ctx, GetBalanceRequest{OwnerId: id1})
	assert.Nil(t, err)
	assert.Equal(t, float32(500), balance)
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

	dep := entity.Deposit{OwnerId: ownerId, Balance: 0}
	err := m.Create(ctx, dep)
	if err != nil {
		return entity.Deposit{}, err
	}

	return dep, nil
}

func (m *mockDepositRepository) Create(ctx context.Context, deposit entity.Deposit) error {
	if deposit.Balance < 0 {
		return errCRUD
	}
	m.items = append(m.items, deposit)
	return nil
}

func (m *mockDepositRepository) Update(ctx context.Context, deposit entity.Deposit) error {
	for i, item := range m.items {
		if item.OwnerId == deposit.OwnerId {
			m.items[i] = deposit
			break
		}
	}
	return nil
}

func (m *mockDepositRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

type mockTransactionRepository struct {
	items          []entity.Transaction
	lastInsertedId int64
}

func (m *mockTransactionRepository) Create(ctx context.Context, tx *entity.Transaction) error {
	if tx.Amount < 0 {
		return errCRUD
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

func (m *mockTransactionRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}
