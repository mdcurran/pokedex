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

	res, err := sdk.GetPokemon(ctx, pokedex.GetRequest{Name: "metapod"})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", res.Pokemon)

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

func main() {
	err := Run()
	if err != nil {
		log.Fatalln(err)
	}
}
