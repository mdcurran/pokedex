package pokedex

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mdcurran/pokedex/models"
)

type GetStatRequest struct {
	ID   int
	Name string
}

type GetStatResponse struct {
	Stat *models.Stat
}

func (c *Client) GetStat(ctx context.Context, r GetStatRequest) (*GetStatResponse, error) {
	stat, err := c.getStat(ctx, "1")
	if err != nil {
		return nil, err
	}
	return &GetStatResponse{Stat: stat}, nil
}

func (c *Client) getStat(ctx context.Context, resource string) (*models.Stat, error) {
	u := c.baseURL.JoinPath("stat", resource)

	b, res, err := c.fetch(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var stat *models.Stat
	err = json.Unmarshal(b, &stat)
	if err != nil {
		return nil, NewError(err.Error(), http.StatusUnprocessableEntity, res)
	}
	c.cache.Set(u.String(), b)

	return stat, nil
}
