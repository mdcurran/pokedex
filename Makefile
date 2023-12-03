.PHONY: test
test:
	go test ./...

.PHONY: init
init: install-tools
	mkdir -p models

.PHONY: generate-models
generate-models: \
	gen/api_resource \
	gen/named_api_resource_list \
	gen/named_api_resource \
	gen/nature \
	gen/pokemon \
	gen/stat \

gen/%: FORCE
	go-jsonschema api/$*.json \
		-p models \
		--capitalization ID \
		--only-models \
		--resolve-extension json \
		--tags json \
	> models/$*.go

.PHONY: install-tools
install-tools:
	go install github.com/atombender/go-jsonschema@latest

FORCE:
