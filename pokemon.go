package pokedex

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mdcurran/pokedex/models"
)

type GetPokemonResponse struct {
	Pokemon *models.Pokemon
}

func (c *Client) GetPokemon(ctx context.Context, r GetRequest) (*GetPokemonResponse, error) {
	resource, err := r.GetResource()
	if err != nil {
		return nil, err
	}
	pokemon, err := c.getPokemon(ctx, resource)
	if err != nil {
		return nil, err
	}
	return &GetPokemonResponse{Pokemon: pokemon}, nil
}

func (c *Client) getPokemon(ctx context.Context, resource string) (*models.Pokemon, error) {
	u := c.baseURL.JoinPath("pokemon", resource)

	b, res, err := c.fetch(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var pokemon *models.Pokemon
	err = json.Unmarshal(b, &pokemon)
	if err != nil {
		return nil, NewError(err.Error(), http.StatusUnprocessableEntity, res)
	}
	c.cache.Set(u.String(), b)

	return pokemon, nil
}
