package pokedex

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mdcurran/pokedex/internal/store"
)

var ErrClientClosed = errors.New("sdk client closed")

const defaultBaseURL = "https://pokeapi.co/api/v2"

type Client struct {
	http    *http.Client
	baseURL *url.URL
	cache   *store.Cache
	// closed indicates if the SDK client has been previously closed.
	// If closed is true the response cache has been shutdown. Therefore we
	// want to prevent requests using a closed client, as no responses would
	// be cached.
	closed bool
}

type Options struct {
	BaseURL string
	Timeout time.Duration
	// A maximum size of ~16MB seems sufficient to hold many responses
	// from the PokéAPI, without claiming an unreasonable amount of user
	// memory.
	CacheMaximumSize int64
	// PokéAPI data is updated very infrequently, so the cache TTL could be a
	// much larger value. However many "real-world" APIs will have much more
	// frequent updates. 10 minutes seems like a reasonable compromise.
	CacheTTL time.Duration
}

func defaultOptions() Options {
	return Options{
		BaseURL:          defaultBaseURL,
		Timeout:          5 * time.Second,
		CacheMaximumSize: 1 << 27, // 16.4MB
		CacheTTL:         10 * time.Minute,
	}
}

// New instantiates a PokéAPI SDK client with a set of sensible client
// defaults.
func New() (*Client, error) {
	return NewWithOptions(defaultOptions())
}

// NewWithOptions instantiates a PokéAPI SDK client with the provided client
// settings.
func NewWithOptions(options Options) (*Client, error) {
	u, err := url.Parse(options.BaseURL)
	if err != nil {
		return nil, err
	}

	cache, err := store.NewCache(store.CacheOptions{
		MaximumSize: options.CacheMaximumSize,
		TTL:         options.CacheTTL,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		http:    &http.Client{Timeout: options.Timeout},
		baseURL: u,
		cache:   cache,
	}, nil
}

func (c *Client) Close() {
	c.cache.Close()
	c.closed = true
}

func (c *Client) fetch(ctx context.Context, url string) ([]byte, *http.Response, error) {
	if c.closed {
		return nil, nil, NewError(ErrClientClosed.Error(), CodeClientClosed, nil)
	}

	b, ok := c.cache.Get(url)
	if ok {
		return b, nil, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, nil, NewError(err.Error(), CodeInternal, nil)
	}

	res, err := c.do(req)
	if err != nil {
		return nil, nil, NewError(err.Error(), CodeInternal, nil)
	}
	if res != nil && res.StatusCode == http.StatusNotFound {
		return nil, res, NewError("Not Found", res.StatusCode, res)
	}
	if res != nil && res.StatusCode != http.StatusOK {
		return nil, res, NewError("Unexpected Error", res.StatusCode, res)
	}
	defer res.Body.Close()

	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, res, NewError(err.Error(), http.StatusUnprocessableEntity, res)
	}

	return b, res, nil
}

func (c *Client) do(r *http.Request) (*http.Response, error) {
	if c.closed {
		return nil, ErrClientClosed
	}
	return c.http.Do(r)
}
