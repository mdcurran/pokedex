package pokedex

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/mdcurran/pokedex/internal/iterator"
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

type ListPokemonResponse struct {
	Iterator *iterator.Paginator[*models.Pokemon]
}

func (c *Client) ListPokemon(ctx context.Context, r ListRequest) (*ListPokemonResponse, error) {
	it := iterator.NewPaginator(ctx, r.PageSize, func(ctx context.Context, start, end uint) ([]*models.Pokemon, error) {
		resourceList, err := c.fetchResourceList(ctx, "pokemon", start, end-start)
		if err != nil {
			return nil, err
		}

		var (
			wg      sync.WaitGroup
			pokemon = make([]*models.Pokemon, len(resourceList.Results))
			errors  = make(chan error)
		)
		for i, item := range resourceList.Results {
			wg.Add(1)
			go func(i int, ctx context.Context, resource string) {
				p, err := c.getPokemon(ctx, resource)
				if err != nil {
					errors <- err
					return
				}
				pokemon[i] = p
				wg.Done()
			}(i, ctx, item.Name)
		}
		wg.Wait()

		select {
		case err := <-errors:
			return nil, err
		default:
			return pokemon, nil
		}
	})

	return &ListPokemonResponse{Iterator: it}, nil
}
