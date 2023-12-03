package pokedex

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mdcurran/pokedex/models"
)

type GetNatureResponse struct {
	Nature *models.Nature
}

func (c *Client) GetNature(ctx context.Context, r GetRequest) (*GetNatureResponse, error) {
	resource, err := r.GetResource()
	if err != nil {
		return nil, err
	}
	nature, err := c.getNature(ctx, resource)
	if err != nil {
		return nil, err
	}
	return &GetNatureResponse{Nature: nature}, nil
}

func (c *Client) getNature(ctx context.Context, resource string) (*models.Nature, error) {
	u := c.baseURL.JoinPath("nature", resource)

	b, res, err := c.fetch(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var nature *models.Nature
	err = json.Unmarshal(b, &nature)
	if err != nil {
		return nil, NewError(err.Error(), http.StatusUnprocessableEntity, res)
	}
	c.cache.Set(u.String(), b)

	return nature, nil
}
