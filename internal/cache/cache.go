package cache

import (
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

// Cache stores data fetched from the PokéAPI in-memory. PokéAPI mandates that
// users should cache values wherever possible. Each record in the cache is the
// response from a single API call, for example:
//
// - https://pokeapi.co/api/v2/pokemon/1
// - https://pokeapi.co/api/v2/nature/hardy
// - https://pokeapi.co/api/v2/stat/3
//
// Cache can be used by multiple goroutines simultaneously.
type Cache struct {
	mu    sync.RWMutex
	cache *ristretto.Cache
	// ttl is how long a cache entry will be available after being set.
	// By default the TTL for all API responses is set to 10 minutes.
	ttl time.Duration
}

type Options struct {
	MaximumSize int64
	TTL         time.Duration
	// debug determines whether cache statistics are kept during the cache's
	// lifecycle. This is useful for testing, as we can assert the result of
	// specific operations against the cache.
	//
	// Cache metrics could be potentially exposed to clients if there's a
	// suitable use-case. For example, users may want to observe the cache
	// hit rate & validate that against some hypothetical API limits, etc.
	debug bool
}

// New instantiates an in-memory cache for PokéAPI responses.
func New(options Options) (*Cache, error) {
	rst, err := ristretto.NewCache(&ristretto.Config{
		// MaxCost is the maximum size (in bytes) of the cache.
		MaxCost: options.MaximumSize,
		// Ristretto recommend setting NumCounters to 10x the MaxCost value.
		NumCounters: options.MaximumSize * 10,
		// Ristretto recommend setting BufferItems to 64 for "generally good
		// performance".
		BufferItems: 64,
		Metrics:     options.debug,
	})
	if err != nil {
		return nil, err
	}
	return &Cache{cache: rst, ttl: options.TTL}, nil
}

// Set adds a PokeAPI response to the cache. There is no guarantee a call to
// Set() will successfully add a record to the cache. However since each HTTP
// response is waited equally with regards to the admission policy, the only
// plausible scenario a record isn't added to the cache is if the cache is
// full.
func (c *Cache) Set(url string, body []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.SetWithTTL(url, body, 0, c.ttl)
	// Wait for the value we've just added to the cache to pass through any
	// internal ristretto buffers. This ensures once we release the write lock
	// the new value is immediately available to subsequent readers.
	c.cache.Wait()
}

// Get fetches a PokéAPI response from the cache. If the endpoints is
// found, the value and a boolean (true) are returned. If the cache does not
// contain a record a given response, value is nil and the boolean false.
func (c *Cache) Get(url string) ([]byte, bool) {
	body, ok := c.cache.Get(url)
	if !ok {
		return nil, false
	}
	return body.([]byte), ok
}

// Close gracefully shutsdown the cache.
func (c *Cache) Close() {
	c.cache.Close()
}
