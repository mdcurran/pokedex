# Pokédex: Go SDK for PokéAPI

## Installation

```
go get github.com/mdcurran/pokedex
```

## Usage

Fetching a single resource (Pokémon):

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mdcurran/pokedex"
)

func main() {
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
```

Fetching a paginated list of a resource (Natures):

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mdcurran/pokedex"
	"github.com/mdcurran/pokedex/internal/iterator"
)

func main() {
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
```

## Supported Operations

- `GetNature` - Get a Nature by ID or Name.
- `GetPokemon` - Get a Pokemon by ID or Name.
- `GetStat` - Get a Stat by ID or Name.
- `ListNatures` - Receive a paginator for all Natures.
- `ListPokemon` - Receive a paginator for all Pokemon.
- `ListStats` - Receive a paginator for all Stats.

## Design

There are a few key principles I wanted to get across with this SDK:

### Code Generation

I didn't want to manually type out all the structs for the required types. This
isn't scalable anyway if trying to build a generalised SDK tool. 

### Caching

The PokéAPI developers mandate that responses should be cached wherever
possible. This SDK comes with client-side caching built-in by default.

### Ease of Use

The "gnarly" parts of the API should be hidden from users. Specifically how
this API handles pagination, where you make 1 request to get URLs for the
resources, and then you need to make follow-up requests to get the data
you care about.

I roughly copied how the BigQuery Go SDK does iteration for query results.
A pagination object that can be iterated until its end.

### Extendibility

There are a bunch of things you can enrich an SDK with. Caching, security, rate
limiting, error handling, etc. I wanted to design the SDK client so it would be
relatively easy to add, say, authentication as a feature. Token-based
authentication seems trivial to add. Certificate-based auth a little more
involved but relatively "plug-and-play".

## Development

To get started with development run: `make init`. This installs any tools
required for local development (i.e. `go-jsonschema`), and setups the
necessary directory structure if needed.

### Testing

To run the test-suite, run: `make test`.

In `examples/` there are test scripts to validate the API worked during
development. They can be run with `go run examples/[gets|lists]/main.go`.
I used these for debugging, but are a good place to mess around with if
you're interested in how the SDK works.

`go run examples/main.go` will run the examples in the README.

### Generating Models

The PokéAPI structs in `models/` were generated from the schemas defined in
the `PokeAPI/api-data` repository ([link](https://github.com/PokeAPI/api-data/tree/master/data/schema/v2)).
I copied the relevant models over (`nature`, `pokemon`, `stat`), as well as
the generic models that are used throughtout the API. There were a few tweaks
to ensure the `$ref` paths are valid.

To generate Go structs from the JSON schemas using the `go-jsonschema` tool
([link](https://github.com/omissis/go-jsonschema)), run:

```
make generate-models
```

The generated structs will require a bit of manual adjustment. Some of the
common models are generated multiple times (e.g. `NamedApiResource`) and
`go-jsonschema` doesn't hand `oneOf` types particularly well. Anything that
is optionally null has been typed with `interface{}`.

## Potential Improvements

A few things that came to mind, but I didn't want to address due to the
time constraints:

- The cache is a bit dumb. If a client calls `/pokemon/1` and then
  `/pokemon/bulbasaur` thereafter, even though the response is already cached
  we'll make a new API request as the cache key is based on the URL.
  Updating the cache key to something like `/pokemon/1:bulbasaur` and then
  doing a string search in `cache.Get()` would fix this.
- When paginating over a list of resources, if one of the goroutines returns
  an error we only inform the client of one error. It might be nice to return
  multiple errors for certain operations.
- The code generation is a bit rough around the edges.
