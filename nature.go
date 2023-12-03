package pokedex

import (
	"context"
	"encoding/json"

	"github.com/mdcurran/pokedex/models"
)

type GetNatureRequest struct {
	ID   int
	Name string
}

type GetNatureResponse struct {
	Nature *models.Nature
}

func (c *Client) GetNature(ctx context.Context, r GetNatureRequest) (*GetNatureResponse, error) {
	nature, err := c.getNature(ctx, "1")
	if err != nil {
		return nil, err
	}
	return &GetNatureResponse{Nature: nature}, nil
}

func (c *Client) getNature(ctx context.Context, resource string) (*models.Nature, error) {
	u := c.baseURL.JoinPath("nature", resource)

	b, err := c.fetch(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var nature *models.Nature
	err = json.Unmarshal(b, &nature)
	if err != nil {
		return nil, err
	}
	c.cache.Set(u.String(), b)

	return nature, nil
}
