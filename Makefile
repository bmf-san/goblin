.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

.PHONY: test
test: ## Run tests.
	go test -v -race ./...

.PHONY: test-cover
test-cover: ## Run tests with cover options.
	go test -v -race -cover -coverprofile=c.out -covermode=atomic ./...

.PHONY: test-benchmark
test-benchmark: ## Run benchmark tests.
	go test -bench=. -cpu=4 -benchmem