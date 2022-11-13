.DEFAULT_GOAL := help

.PHONY: consul consul-destroy test test-env coverage help html html-gen html-wsl

consul: ## Run consul in docker
	docker run -d --name consul -p 8500:8500 consul

consul-destroy: ## Destroy consul container
	docker rm -f consul

test-env: TEST ?= test
test-env: consul $(TEST) consul-destroy ## Up consul and run tests

test: ## Run tests
	go test -race ./...

coverage: ## Run unit tests with coverage
	@go test -v -race -cover -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

html: ## Show html coverage result
	@go tool cover -html=./coverage.out

html-gen: ## Export html coverage result
	@go tool cover -html=./coverage.out -o ./coverage.html

html-wsl: html-gen ## Open html coverage result in wsl
	@explorer.exe `wslpath -w ./coverage.html` || true

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
