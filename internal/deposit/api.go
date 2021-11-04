package deposit

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"users-balance-microservice/internal/errors"
	"users-balance-microservice/internal/requests"
	"users-balance-microservice/internal/transaction"
	"users-balance-microservice/pkg/dbcontext"
	"users-balance-microservice/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(
	r *routing.RouteGroup,
	depositService Service,
	transactionService transaction.Service,
	logger log.Logger,
	db *dbcontext.DB,
) {
	res := resource{depositService, transactionService, logger}

	r.Post("/deposits/balance", res.getBalance)
	r.Post("/deposits/update", db.TransactionHandler(), res.updateBalance)
	r.Post("/deposits/transfer", db.TransactionHandler(), res.transfer)
	r.Post("/deposits/history", res.history)
}

type resource struct {
	depositService     Service
	transactionService transaction.Service
	logger             log.Logger
}

func (r resource) getBalance(c *routing.Context) error {
	var input requests.GetBalanceRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	balance, err := r.depositService.GetBalance(c.Request.Context(), input)
	if err != nil {
		return err
	}
	return c.Write(balance)
}

func (r resource) updateBalance(c *routing.Context) error {
	var input requests.UpdateBalanceRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	err := r.depositService.Update(c.Request.Context(), input)
	if err != nil {
		return err
	}
	tx, err := r.transactionService.CreateUpdateTransaction(c.Request.Context(), input)
	if err != nil {
		return err
	}
	return c.Write(tx)
}

func (r resource) transfer(c *routing.Context) error {
	var input requests.TransferRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	err := r.depositService.Transfer(c.Request.Context(), input)
	if err != nil {
		return err
	}
	tx, err := r.transactionService.CreateTransferTransaction(c.Request.Context(), input)
	if err != nil {
		return err
	}
	return c.Write(tx)
}

func (r resource) history(c *routing.Context) error {
	var input requests.GetHistoryRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	transactions, err := r.transactionService.GetHistory(c.Request.Context(), input)
	if err != nil {
		return err
	}
	return c.Write(transactions)
}
