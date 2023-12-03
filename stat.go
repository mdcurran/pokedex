package pokedex

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mdcurran/pokedex/models"
)

type GetStatResponse struct {
	Stat *models.Stat
}

func (c *Client) GetStat(ctx context.Context, r GetRequest) (*GetStatResponse, error) {
	resource, err := r.GetResource()
	if err != nil {
		return nil, err
	}
	stat, err := c.getStat(ctx, resource)
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
