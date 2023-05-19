# Vault

Vault loader, read/write key-value pairs from vault.

```sh
go get github.com/rytsh/liz/vault
```

## Usage

| Environment variable    | Description          |
| ----------------------- | -------------------- |
| VAULT_ROLE_ID           | Role ID              |
| VAULT_ROLE_SECRET       | Role secret          |
| VAULT_ADDR              | Vault address        |
| VAULT_TOKEN             | Vault token          |
| VAULT_AGENT_ADDR        | Vault agent address  |
| VAULT_MAX_RETRIES       | Max retries          |
| VAULT_CACERT            | CA certificate       |
| VAULT_CACERT_BYTES      | CA certificate bytes |
| VAULT_CAPATH            | CA path              |
| VAULT_CLIENT_CERT       | Client certificate   |
| VAULT_CLIENT_KEY        | Client key           |
| VAULT_RATE_LIMIT        | Rate limit           |
| VAULT_CLIENT_TIMEOUT    | Client timeout       |
| VAULT_SKIP_VERIFY       | Skip verify          |
| VAULT_SRV_LOOKUP        | SRV lookup           |
| VAULT_TLS_SERVER_NAME   | TLS server name      |
| VAULT_HTTP_PROXY        | HTTP proxy           |
| VAULT_PROXY_ADDR        | Proxy address        |
| VAULT_DISABLE_REDIRECTS | Disable redirects    |
| VAULT_APPROLE_BASE_PATH | /auth/approle/login/ |

## Test

First initialize vault:

```sh
make vault
# make vault-destroy
```

Login with token:

```sh
make vault-login TOKEN=...
```

Set approle:

```sh
make vault-role
```

Then run tests:

```sh
VAULT_ROLE_ID=... make test
```
