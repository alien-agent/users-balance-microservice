package transaction

import (
	"net/http"

	"github.com/go-ozzo/ozzo-routing/v2"
	"users-balance-microservice/internal/errors"
	"users-balance-microservice/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger log.Logger) {
	res := resource{service, logger}

	r.Post("/deposits/change", res.createOneway)
	r.Post("/deposits/transfer", res.createTwoway)
	r.Post("/deposits/history", res.getForUser)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) createOneway(c *routing.Context) error {
	var input CreateOnewayTransactionRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	transaction, err := r.service.CreateOneway(c.Request.Context(), input)
	if err != nil{
		return err
	}
	return c.WriteWithStatus(transaction, http.StatusCreated)
}

func (r resource) createTwoway(c *routing.Context) error {
	var input CreateTwowayTransactionRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	transaction, err := r.service.CreateTwoway(c.Request.Context(), input)
	if err != nil {
		return err
	}
	return c.WriteWithStatus(transaction, http.StatusCreated)
}

func (r resource) getForUser(c *routing.Context) error {
	var input GetHistoryRequest
	if err := c.Read(&input); err != nil{
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	transactions, err := r.service.GetForUser(c.Request.Context(), input)
	if err != nil{
		return err
	}
	return c.Write(transactions)
}
