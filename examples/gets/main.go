package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mdcurran/pokedex"
)

func Run() error {
	ctx := context.Background()

	sdk, err := pokedex.New()
	if err != nil {
		return err
	}
	defer sdk.Close()

	err = nature(ctx, sdk)
	if err != nil {
		return err
	}
	err = pokemon(ctx, sdk)
	if err != nil {
		return err
	}
	err = stat(ctx, sdk)
	if err != nil {
		return err
	}

	return nil
}

func nature(ctx context.Context, sdk *pokedex.Client) error {
	res, err := sdk.GetNature(ctx, pokedex.GetRequest{ID: 2})
	if err != nil {
		return err
	}

	// Hit cache.
	for i := 1; i <= 10; i++ {
		res, err = sdk.GetNature(ctx, pokedex.GetRequest{ID: 2})
		if err != nil {
			return err
		}
		fmt.Printf("cached %d name: %s\n", i, res.Nature.Name)
	}

	return nil
}

func pokemon(ctx context.Context, sdk *pokedex.Client) error {
	res, err := sdk.GetPokemon(ctx, pokedex.GetRequest{Name: "metapod"})
	if err != nil {
		return err
	}

	// Hit cache.
	for i := 1; i <= 10; i++ {
		res, err = sdk.GetPokemon(ctx, pokedex.GetRequest{Name: "metapod"})
		if err != nil {
			return err
		}
		fmt.Printf("cached %d name: %s\n", i, res.Pokemon.Name)
	}

	return nil
}

func stat(ctx context.Context, sdk *pokedex.Client) error {
	res, err := sdk.GetStat(ctx, pokedex.GetRequest{ID: 4})
	if err != nil {
		return err
	}

	// Hit cache.
	for i := 1; i <= 10; i++ {
		res, err = sdk.GetStat(ctx, pokedex.GetRequest{ID: 4})
		if err != nil {
			return err
		}
		fmt.Printf("cached %d name: %s\n", i, res.Stat.Name)
	}

	return nil
}

func main() {
	err := Run()
	if err != nil {
		log.Fatalln(err)
	}
}
