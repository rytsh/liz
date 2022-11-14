# Consul

Consul loader, read/write/delete key-value pairs from consul.

## Usage

| Environment variable    | Description          |
| ----------------------- | -------------------- |
| CONSUL_HTTP_ADDR |  Consul address |
| CONSUL_HTTP_TOKEN_FILE | Consul token file |
| CONSUL_HTTP_TOKEN | Consul token |
| CONSUL_HTTP_AUTH | Consul auth |
| CONSUL_HTTP_SSL | Consul ssl |
| CONSUL_TLS_SERVER_NAME | Consul tls server name |
| CONSUL_CACERT | Consul cacert |
| CONSUL_CAPATH | Consul capath |
| CONSUL_CLIENT_CERT | Consul client cert |
| CONSUL_CLIENT_KEY | Consul client key |
| CONSUL_HTTP_SSL_VERIFY | Consul ssl verify |
| CONSUL_NAMESPACE | Consul namespace |
| CONSUL_PARTITION | Consul partition |

## Test

```sh
make test-env
# make test-env html-wsl TEST=coverage
```
