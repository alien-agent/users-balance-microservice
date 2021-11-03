package exchangerates

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"users-balance-microservice/pkg/log"
)

const (
	baseCurrency = "RUB"
	apiPath      = "https://api.exchangerate.host/latest"
)

var currencyUnavailableError = errors.New("currency is not present in either cache or API response")

// RatesService handles exchange rates.
type RatesService interface {
	// Get returns the exchange ratio for specific currency code against baseCurrency(RUB).
	Get(code string) (float32, error)
}

type service struct {
	cache  *CacheService
	logger log.Logger
}

// NewService creates a new exchange rates service.
func NewService(expiry time.Duration, logger log.Logger) RatesService {
	store := cache.New(expiry, 5*time.Minute)
	cacheService := NewCacheService(store)
	return service{cache: cacheService, logger: logger}
}

// RatesResponse holds an API response with a list of RUB\CURRENCY ratios for all currencies.
type RatesResponse struct {
	Rates map[string]float32 `json:"rates"`
}

// Get will fetch a single rate for a given currency either from the cache or the API.
func (s service) Get(code string) (float32, error) {
	if code == baseCurrency {
		return 1, nil
	}

	// If we have cached results, use them.
	if result, ok := s.cache.Get(code); ok {
		return result, nil
	}

	// No cached results, go and fetch them.
	if err := s.fetch(); err != nil {
		s.logger.Error("failed to fetch currency rates: ", err)
		return 0, err
	}

	// Currency should be in cache by now. If failed, then particular currency is unavailable in service right now.
	if result, ok := s.cache.Get(code); ok {
		return result, nil
	} else {
		s.logger.Info(fmt.Sprintf("client requested rate for \"%s\", which was not found in API response", code))
		return 0, currencyUnavailableError
	}
}

// Fetch all RUB/CURRENCY rates from API.
func (s service) fetch() error {
	fullUrl := fmt.Sprintf("%s?base=%s", apiPath, baseCurrency)
	response, err := http.Get(fullUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	latest := RatesResponse{}
	err = json.NewDecoder(response.Body).Decode(&latest)
	if err != nil {
		return err
	}

	// Store our results.
	s.cache.Store(&latest)

	return nil
}
