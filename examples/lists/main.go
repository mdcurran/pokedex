package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mdcurran/pokedex"
	"github.com/mdcurran/pokedex/iterator"
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
	err = pokemon(ctx, sdk)
	if err != nil {
		return err
	}
	err = stats(ctx, sdk)
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

func pokemon(ctx context.Context, sdk *pokedex.Client) error {
	res, err := sdk.ListPokemon(ctx, pokedex.ListRequest{PageSize: 20})
	if err != nil {
		return err
	}
	it := res.Iterator

	// There are lots of pokemon, for this test let's just exit after
	// a few pages so we don't flood the API.
	for i := 0; i < 10; i++ {
		pokemon, err := it.Next(ctx)
		if err == iterator.EndOfIterator {
			return nil
		}
		if err != nil {
			return err
		}
		for _, p := range pokemon {
			fmt.Printf("id: %d name: %s\n", p.ID, p.Name)
			if len(p.Abilities) > 0 {
				first := p.Abilities[0]
				fmt.Printf("first %s ability: %s\n", p.Name, first.Ability.Name)
			}
		}
	}

	return nil
}

func stats(ctx context.Context, sdk *pokedex.Client) error {
	res, err := sdk.ListStats(ctx, pokedex.ListRequest{PageSize: 10})
	if err != nil {
		return err
	}
	it := res.Iterator

	for {
		stats, err := it.Next(ctx)
		if err == iterator.EndOfIterator {
			break
		}
		if err != nil {
			return err
		}
		for _, s := range stats {
			fmt.Printf("id: %d name: %s\n", s.ID, s.Name)
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
