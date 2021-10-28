package deposit

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/pkg/log"
)

var errCRUD = errors.New("error crud")

func TestGetBalanceRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     GetBalanceRequest
		wantError bool
	}{
		{"success", GetBalanceRequest{uuid.NewString()}, false},
		{"missing OwnerId", GetBalanceRequest{""}, true},
		{"invalid OwnerId", GetBalanceRequest{"12712912"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRepository{}, logger)

	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// get non-existing deposit -> new deposit created with balance=0
	ownerId := uuid.New()
	deposit, err := s.Get(ctx, GetBalanceRequest{OwnerId: ownerId.String()})
	assert.Nil(t, err)
	assert.Equal(t, ownerId, deposit.OwnerId)
	assert.Equal(t, int64(0), deposit.Balance)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in get
	_, err = s.Get(ctx, GetBalanceRequest{OwnerId: "it is invalid UUID"})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// get existing deposit -> its balance is returned
	deposit, err = s.Get(ctx, GetBalanceRequest{OwnerId: ownerId.String()})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), deposit.Balance)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
}

type mockRepository struct {
	items []entity.Deposit
}

func (m *mockRepository) Get(ctx context.Context, ownerId uuid.UUID) (entity.Deposit, error) {
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

func (m *mockRepository) Create(ctx context.Context, deposit entity.Deposit) error {
	if deposit.Balance < 0 {
		return errCRUD
	}
	m.items = append(m.items, deposit)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, deposit entity.Deposit) error {
	for i, item := range m.items {
		if item.OwnerId == deposit.OwnerId {
			m.items[i] = deposit
			break
		}
	}
	return nil
}

func (m *mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}
