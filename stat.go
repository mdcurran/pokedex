package pokedex

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/mdcurran/pokedex/iterator"
	"github.com/mdcurran/pokedex/models"
)

type GetStatResponse struct {
	Stat *models.Stat
}

// GetStat returns a single Stat according to an ID or name.
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

type ListStatsResponse struct {
	Iterator *iterator.Paginator[*models.Stat]
}

// ListStats returns an iterator with a user-provided page size over all Stats.
func (c *Client) ListStats(ctx context.Context, r ListRequest) (*ListStatsResponse, error) {
	it := iterator.NewPaginator(ctx, r.PageSize, func(ctx context.Context, start, end uint) ([]*models.Stat, error) {
		resourceList, err := c.fetchResourceList(ctx, "stat", start, end-start)
		if err != nil {
			return nil, err
		}

		var (
			wg     sync.WaitGroup
			stats  = make([]*models.Stat, len(resourceList.Results))
			errors = make(chan error)
		)
		for i, item := range resourceList.Results {
			wg.Add(1)
			go func(i int, ctx context.Context, resource string) {
				nature, err := c.getStat(ctx, resource)
				if err != nil {
					errors <- err
					return
				}
				stats[i] = nature
				wg.Done()
			}(i, ctx, item.Name)
		}
		wg.Wait()

		select {
		case err := <-errors:
			return nil, err
		default:
			return stats, nil
		}
	})

	return &ListStatsResponse{Iterator: it}, nil
}
