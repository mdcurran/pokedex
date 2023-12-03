package pokedex

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/mdcurran/pokedex/internal/iterator"
	"github.com/mdcurran/pokedex/models"
)

type GetNatureResponse struct {
	Nature *models.Nature
}

// GetNature returns a single Nature according to an ID or name.
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

type ListNaturesResponse struct {
	Iterator *iterator.Paginator[*models.Nature]
}

// ListNatures returns an iterator with a user-provided page size over all
// Natures.
func (c *Client) ListNatures(ctx context.Context, r ListRequest) (*ListNaturesResponse, error) {
	it := iterator.NewPaginator(ctx, r.PageSize, func(ctx context.Context, start, end uint) ([]*models.Nature, error) {
		resourceList, err := c.fetchResourceList(ctx, "nature", start, end-start)
		if err != nil {
			return nil, err
		}

		var (
			wg      sync.WaitGroup
			natures = make([]*models.Nature, len(resourceList.Results))
			errors  = make(chan error)
		)
		// As we know the number of results from the NamedApiResourceList
		// we can create a slice that size and each goroutine updates its
		// own memory based on the index i.
		for i, item := range resourceList.Results {
			wg.Add(1)
			go func(i int, ctx context.Context, resource string) {
				nature, err := c.getNature(ctx, resource)
				if err != nil {
					errors <- err
					return
				}
				natures[i] = nature
				wg.Done()
			}(i, ctx, item.Name)
		}
		wg.Wait()

		select {
		case err := <-errors:
			return nil, err
		default:
			return natures, nil
		}
	})

	return &ListNaturesResponse{Iterator: it}, nil
}
