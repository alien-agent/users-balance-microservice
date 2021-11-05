package transaction

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"users-balance-microservice/internal/entity"
	"users-balance-microservice/internal/requests"
	"users-balance-microservice/pkg/log"
)

var (
	databaseError = errors.New("database error")
	logger, _     = log.NewForTest()
	ctx           = context.Background()
)

func TestService_CreateUpdateTransaction(t *testing.T) {
	id1 := uuid.New()
	s := NewService(&mockTransactionRepository{}, logger)

	// initial count
	count, err := s.Count(ctx)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 0, count)
	}

	// success positive amount
	tx, err := s.CreateUpdateTransaction(ctx, requests.UpdateBalanceRequest{
		OwnerId:     id1.String(),
		Amount:      1000,
		Description: "visa top-up",
	})
	if assert.NoError(t, err) {
		assert.Equal(t, uuid.Nil, tx.SenderId)
		assert.Equal(t, id1, tx.RecipientId)
		assert.EqualValues(t, 1000, tx.Amount)
		assert.Equal(t, "visa top-up", tx.Description)

		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1, count2-count)
			count++
		}
	}

	// success negative amount
	tx, err = s.CreateUpdateTransaction(ctx, requests.UpdateBalanceRequest{
		OwnerId:     id1.String(),
		Amount:      -1000,
		Description: "monthly subscription",
	})
	if assert.NoError(t, err) {
		assert.Equal(t, id1, tx.SenderId)
		assert.Equal(t, uuid.Nil, tx.RecipientId)
		assert.EqualValues(t, 1000, tx.Amount)
		assert.Equal(t, "monthly subscription", tx.Description)

		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1, count2-count)
			count++
		}
	}

	// fail invalid ownerId
	tx, err = s.CreateUpdateTransaction(ctx, requests.UpdateBalanceRequest{
		OwnerId:     "1234-5678-9",
		Amount:      1000,
		Description: "monthly subscription",
	})
	if assert.Error(t, err) {
		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 0, count2-count)
		}
	}

	// fail database error
	tx, err = s.CreateUpdateTransaction(ctx, requests.UpdateBalanceRequest{
		OwnerId:     "11111111-1111-1111-1111-111111111111",
		Amount:      1000,
		Description: "monthly subscription",
	})
	if assert.Error(t, err) {
		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 0, count2-count)
		}
	}
}

func TestService_CreateTransferTransaction(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	s := NewService(&mockTransactionRepository{}, logger)

	// initial count
	count, err := s.Count(ctx)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 0, count)
	}

	// success
	tx, err := s.CreateTransferTransaction(ctx, requests.TransferRequest{
		SenderId:    id1.String(),
		RecipientId: id2.String(),
		Amount:      1000,
		Description: "thanks for dinner!",
	})
	if assert.NoError(t, err) {
		assert.Equal(t, id1, tx.SenderId)
		assert.Equal(t, id2, tx.RecipientId)
		assert.EqualValues(t, 1000, tx.Amount)
		assert.Equal(t, "thanks for dinner!", tx.Description)

		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1, count2-count)
			count++
		}
	}

	// success no description
	tx, err = s.CreateTransferTransaction(ctx, requests.TransferRequest{
		SenderId:    id1.String(),
		RecipientId: id2.String(),
		Amount:      1000,
		Description: "",
	})
	if assert.NoError(t, err) {
		assert.Equal(t, id1, tx.SenderId)
		assert.Equal(t, id2, tx.RecipientId)
		assert.EqualValues(t, 1000, tx.Amount)
		assert.Equal(t, "", tx.Description)

		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 1, count2-count)
			count++
		}
	}

	// fail negative amount
	tx, err = s.CreateTransferTransaction(ctx, requests.TransferRequest{
		SenderId:    id1.String(),
		RecipientId: id2.String(),
		Amount:      -1000,
		Description: "hacker attack!",
	})
	if assert.Error(t, err) {
		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 0, count2-count)
		}
	}

	// fail missing SenderId
	tx, err = s.CreateTransferTransaction(ctx, requests.TransferRequest{
		SenderId:    "",
		RecipientId: id2.String(),
		Amount:      1000,
		Description: "hacker attack!",
	})
	if assert.Error(t, err) {
		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 0, count2-count)
		}
	}

	// fail missing RecipientId
	tx, err = s.CreateTransferTransaction(ctx, requests.TransferRequest{
		SenderId:    id1.String(),
		RecipientId: "",
		Amount:      1000,
		Description: "hacker attack!",
	})
	if assert.Error(t, err) {
		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 0, count2-count)
		}
	}

	// fail too long description
	tx, err = s.CreateTransferTransaction(ctx, requests.TransferRequest{
		SenderId:    id1.String(),
		RecipientId: id2.String(),
		Amount:      1000,
		Description: strings.Repeat("test", 100),
	})
	if assert.Error(t, err) {
		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 0, count2-count)
		}
	}

	// fail database error
	tx, err = s.CreateTransferTransaction(ctx, requests.TransferRequest{
		SenderId:    "11111111-1111-1111-1111-111111111111",
		RecipientId: id2.String(),
		Amount:      1000,
		Description: "",
	})
	if assert.Error(t, err) {
		count2, err := s.Count(ctx)
		if assert.NoError(t, err) {
			assert.EqualValues(t, 0, count2-count)
		}
	}
}

func TestService_GetHistory(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	txsList := []entity.Transaction{
		{Id: 0, SenderId: id1, RecipientId: id2, Amount: 1000, Description: "transfer1"},
		{Id: 1, SenderId: id2, RecipientId: id1, Amount: 2000, Description: "transfer2"},
		{Id: 2, SenderId: id1, RecipientId: id2, Amount: 3000, Description: "transfer3"},
		{Id: 3, SenderId: uuid.Nil, RecipientId: id1, Amount: 4000, Description: "top-up"},
		{Id: 4, SenderId: id1, RecipientId: uuid.Nil, Amount: 5000, Description: "withdrawal"},
	}
	s := NewService(&mockTransactionRepository{items: txsList}, logger)

	// success id1's transactions
	txs, err := s.GetHistory(ctx, requests.GetHistoryRequest{OwnerId: id1.String()})
	if assert.NoError(t, err) {
		assert.Equal(t, txsList, txs)
	}

	// success id2's transactions
	txs, err = s.GetHistory(ctx, requests.GetHistoryRequest{OwnerId: id2.String()})
	if assert.NoError(t, err) {
		assert.Equal(t, txsList[:3], txs)
	}

	// success id2's transactions
	txs, err = s.GetHistory(ctx, requests.GetHistoryRequest{OwnerId: id2.String()})
	if assert.NoError(t, err) {
		assert.Equal(t, txsList[:3], txs)
	}

	// fail invalid OwnerId
	txs, err = s.GetHistory(ctx, requests.GetHistoryRequest{OwnerId: "123-456-789"})
	assert.Error(t, err)

	// fail invalid limit and offset
	txs, err = s.GetHistory(ctx, requests.GetHistoryRequest{OwnerId: id1.String(), Limit: -1, Offset: -1})
	assert.Error(t, err)

	// fail database error
	txs, err = s.GetHistory(ctx, requests.GetHistoryRequest{OwnerId: "11111111-1111-1111-1111-111111111111"})
	assert.Error(t, err)
}

type mockTransactionRepository struct {
	items          []entity.Transaction
	lastInsertedId int64
}

func (m *mockTransactionRepository) Create(ctx context.Context, tx *entity.Transaction) error {
	if tx.Amount < 0 {
		return databaseError
	}
	// simulate database error
	if tx.SenderId.String() == "11111111-1111-1111-1111-111111111111" || tx.RecipientId.String() == "11111111-1111-1111-1111-111111111111" {
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

	// simulate database error
	if ownerId.String() == "11111111-1111-1111-1111-111111111111" {
		return result, databaseError
	}

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
