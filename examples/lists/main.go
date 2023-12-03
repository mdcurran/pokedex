package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mdcurran/pokedex"
	"github.com/mdcurran/pokedex/internal/iterator"
)

func Run() error {
	ctx := context.Background()

	sdk, err := pokedex.New()
	if err != nil {
		return err
	}

	err = natures(ctx, sdk)
	if err != nil {
		return err
	}

	return nil
}

func natures(ctx context.Context, sdk *pokedex.Client) error {
	res, err := sdk.ListNatures(ctx, pokedex.ListRequest{PageSize: 10})
	if err != nil {
		return err
	}
	it := res.Iterator

	for {
		natures, err := it.Next(ctx)
		if err == iterator.EndOfIterator {
			break
		}
		if err != nil {
			return err
		}
		for _, n := range natures {
			fmt.Printf("id: %d name: %s\n", n.ID, n.Name)
		}
	}

	return nil
}

func main() {
	err := Run()
	if err != nil {
		log.Fatalln(err)
	}
}
