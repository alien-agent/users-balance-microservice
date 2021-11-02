package deposit

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"users-balance-microservice/internal/errors"
	"users-balance-microservice/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	r.Post("/deposits/balance", res.getBalance)
	r.Post("/deposits/update", res.updateBalance)
	r.Post("/deposits/transfer", res.transfer)
	r.Post("/deposits/history", res.history)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getBalance(c *routing.Context) error {
	var input GetBalanceRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	balance, err := r.service.GetBalance(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.Write(balance)
}

func (r resource) updateBalance(c *routing.Context) error {
	var input UpdateBalanceRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	transaction, err := r.service.Update(c.Request.Context(), input)
	if err != nil {
		return err
	}
	return c.Write(transaction)
}

func (r resource) transfer(c *routing.Context) error {
	var input TransferRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	transaction, err := r.service.Transfer(c.Request.Context(), input)
	if err != nil {
		return err
	}
	return c.Write(transaction)
}

func (r resource) history(c *routing.Context) error {
	var input GetHistoryRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	transactions, err := r.service.GetHistory(c.Request.Context(), input)
	if err != nil {
		return err
	}
	return c.Write(transactions)
}
