package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mdcurran/pokedex"
	"github.com/mdcurran/pokedex/internal/iterator"
)

// These are the examples from the README.
func main() {
	get()
	list()
}

func get() {
	ctx := context.Background()

	sdk, err := pokedex.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sdk.Close()

	// Metapod good.
	res, err := sdk.GetPokemon(ctx, pokedex.GetRequest{Name: "metapod"})
	if err != nil {
		log.Fatal(err)
	}
	pokemon := res.Pokemon

	fmt.Printf("%s's ID is %d\n", pokemon.Name, pokemon.ID)
	if len(pokemon.Abilities) > 0 {
		first := pokemon.Abilities[0]
		fmt.Printf("%s's first ability is called: %q\n", pokemon.Name, first.Ability.Name)
	}

	// ...
}

func list() {
	ctx := context.Background()

	sdk, err := pokedex.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sdk.Close()

	res, err := sdk.ListNatures(ctx, pokedex.ListRequest{PageSize: 10})
	if err != nil {
		log.Fatal(err)
	}
	it := res.Iterator

	for {
		natures, err := it.Next(ctx)
		if err == iterator.EndOfIterator {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, n := range natures {
			fmt.Printf("id: %d name: %s\n", n.ID, n.Name)
		}
	}

	// ...
}
