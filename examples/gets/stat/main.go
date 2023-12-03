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

	res, err := sdk.GetStat(ctx, pokedex.GetStatRequest{})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", res.Stat)

	// Hit cache.
	for i := 1; i <= 10; i++ {
		res, err = sdk.GetStat(ctx, pokedex.GetStatRequest{})
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
