# Pokédex: Go SDK for PokéAPI

## Development

To get started with development run: `make init`. This installs any tools
required for local development (i.e. `go-jsonschema`), and setups the
necessary directory structure if needed.

### Testing

To run the test-suite, run: `make test`.

In `examples/` there are test scripts to validate the API worked during
development. They can be run with `go run examples/<DIR>/main.go`, but they're
mostly for my own debugging.

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
