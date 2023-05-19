export VAULT_ADDR ?= http://127.0.0.1:8200
export VAULT_ROLE_ID ?= xxx

.DEFAULT_GOAL := help

.PHONY: vault vault-destroy test test-env coverage help html html-gen html-wsl

vault: ## Run vault in docker
	docker run -d --name vault -p 8200:8200 vault

vault-destroy: ## Destroy vault container
	docker rm -f vault

vault-login: TOKEN ?= root
vault-login: ## Login to vault
	vault login $(TOKEN)

.ONESHELL:
vault-role:
	vault auth enable approle
	cat <<-EOF | vault policy write secret-read -
		path "secret/*" {
		capabilities = ["read", "list"]
		}
	EOF
	vault write auth/approle/role/my-role bind_secret_id=false secret_id_bound_cidrs="0.0.0.0/0" policies="default,secret-read"
	vault read -field=role_id auth/approle/role/my-role/role-id

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
