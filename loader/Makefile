.DEFAULT_GOAL := help

.PHONY: vault consul test test-env coverage help html html-gen html-wsl

vault: ## Run vault in docker -it
	docker run -it --name vault -p 8200:8200 vault

consul: ## Run consul in docker -it
	docker run -it --name consul -p 8500:8500 consul

test: export CONSUL_HTTP_ADDR ?= http://localhost:8500
test: export VAULT_ADDR ?= http://localhost:8200
# test: export VAULT_ROLE_ID ?= 88eda05e-b98e-dda4-7251-e97a0638adc9
test: ## Run tests
	 go test -v -race ./...

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
