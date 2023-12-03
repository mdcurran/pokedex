package pokedex

import (
	"context"
	"encoding/json"

	"github.com/mdcurran/pokedex/models"
)

type GetPokemonRequest struct {
	ID   int
	Name string
}

type GetPokemonResponse struct {
	Pokemon *models.Pokemon
}

func (c *Client) GetPokemon(ctx context.Context, r GetPokemonRequest) (*GetPokemonResponse, error) {
	pokemon, err := c.getPokemon(ctx, "1")
	if err != nil {
		return nil, err
	}
	return &GetPokemonResponse{Pokemon: pokemon}, nil
}

func (c *Client) getPokemon(ctx context.Context, resource string) (*models.Pokemon, error) {
	u := c.baseURL.JoinPath("pokemon", resource)

	b, err := c.fetch(ctx, u.String())
	if err != nil {
		return nil, err
	}

	var pokemon *models.Pokemon
	err = json.Unmarshal(b, &pokemon)
	if err != nil {
		return nil, err
	}
	c.cache.Set(u.String(), b)

	return pokemon, nil
}
