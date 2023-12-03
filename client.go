package pokedex

import (
	"net/http"
	"net/url"
	"time"
)

const defaultBaseURL = "https://pokeapi.co/api/v2"

type Client struct {
	http    *http.Client
	baseURL *url.URL
}

type Options struct {
	BaseURL string
	Timeout time.Duration
}

func defaultOptions() Options {
	return Options{
		BaseURL: defaultBaseURL,
		Timeout: 5 * time.Second,
	}
}

func New() (*Client, error) {
	return NewWithOptions(defaultOptions())
}

func NewWithOptions(options Options) (*Client, error) {
	u, err := url.Parse(options.BaseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		http:    &http.Client{Timeout: options.Timeout},
		baseURL: u,
	}, nil
}

func (c *Client) do(r *http.Request) (*http.Response, error) {
	return c.http.Do(r)
}
