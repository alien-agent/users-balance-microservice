package rates

import (
	"github.com/patrickmn/go-cache"
)

// CacheService handles in-memory caching of exchange rates.
type CacheService struct {
	store *cache.Cache
}

// NewCacheService creates a new handler for this service.
func NewCacheService(store *cache.Cache) *CacheService {
	return &CacheService{store}
}

// Get will return our in-memory stored currency/rates.
func (s *CacheService) Get(code string) (float32, bool) {
	if x, found := s.store.Get(code); found {
		return x.(float32), found
	}
	return 0, false
}

// Store unpacks currencies and corresponding rates from ratesResponse and saves them to cache.
func (s *CacheService) Store(rsp *ratesResponse) {
	for code, rate := range rsp.Rates {
		s.store.Set(code, rate, cache.DefaultExpiration)
	}
}

// IsExpired checks whether the rate stored is expired.
func (s *CacheService) IsExpired(code string) bool {
	if _, found := s.store.Get(code); found {
		return false
	}
	return true
}

// Expire will expire the cache for a given currency code.
func (s *CacheService) Expire(code string) {
	s.store.Delete(code)
}
