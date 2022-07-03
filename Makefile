.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

.PHONY: install-staticcheck
install-staticcheck: ## Install staticcheck.
ifeq ($(shell command -v staticcheck 2> /dev/null),)
	go install honnef.co/go/tools/cmd/staticcheck@latest
endif

.PHONY: gofmt
gofmt: ## Run gofmt.
	test -z "$(gofmt -s -l . | tee /dev/stderr)"

.PHONY: vet
vet: ## Run vet.
	go vet -v ./...

.PHONY: staticcheck
staticcheck: ## Run staticcheck.
	staticcheck ./...

.PHONY: test
test: ## Run tests.
	go test -v -race ./...

.PHONY: test-cover
test-cover: ## Run tests with cover options. ex. make test-cover OUT="c.out"
	go test -v -race -cover -coverprofile=$(OUT) -covermode=atomic ./...

.PHONY: test-benchmark
test-benchmark: ## Run benchmark tests.
	go test -bench=. -cpu=4 -benchmem