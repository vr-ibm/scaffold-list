MM ?= ~/Documents/github/magic-modules
OUTPUT ?= /tmp

.PHONY: build test generate clean

build:
	go build ./...

test:
	go test ./...

generate:
	go run cmd/scaffold-list/main.go \
		--magic-modules $(MM) \
		--resource $(RESOURCE) \
		--output $(OUTPUT)

# Regenerate the checked-in example
examples/data_source_managed_zone_list.go:
	go run cmd/scaffold-list/main.go \
		--magic-modules $(MM) \
		--resource dns/ManagedZone.yaml \
		--output examples

clean:
	rm -f /tmp/data_source_*_list.go
