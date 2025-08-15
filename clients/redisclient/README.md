# redisclient

[![Go Reference](https://pkg.go.dev/badge/gitlab.globoi.com/globoplay/go-prime/clients/redisclient.svg)](https://pkg.go.dev/gitlab.globoi.com/globoplay/go-prime/clients/redisclient)

Cliente Redis extensível para Go, baseado em [go-redis](https://github.com/redis/go-redis), com suporte a conexões Standalone, Cluster e Sentinel, interface simples e integração fácil com middlewares e cache.

## Instalação

```bash
go get gitlab.globoi.com/globoplay/go-prime/clients/redisclient
```

## Visão Geral

- Suporte a Redis Standalone, Cluster e Sentinel
- Interface simples: `Get`, `Set`
- Ideal para cache, filas, locks e integrações
- Pronto para uso em middlewares HTTP

## Modos de Conexão

### Standalone

Conexão simples com um único nó Redis:

```go
import "gitlab.globoi.com/globoplay/go-prime/clients/redisclient"

client := redisclient.NewRedisClientFromURL("redis://localhost:6379")
```

### Cluster

Conexão com Redis Cluster:

```go
client := redisclient.NewRedisClientFromURL("redis+cluster://host1:6379,host2:6379,host3:6379")
```

### Sentinel

Conexão com Redis Sentinel:

```go
client := redisclient.NewRedisClientFromURL("redis+sentinel://sentinel1:26379,sentinel2:26379/service_name:mymaster")
```

## Interface e Métodos

```go
type IRedisClient interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value any, expiration time.Duration) error
}
```

- **Get**: Busca o valor de uma chave
- **Set**: Define o valor de uma chave, com expiração opcional

## Exemplos de Uso

```go
import (
    "context"
    "log"
    "time"
    "gitlab.globoi.com/globoplay/go-prime/clients/redisclient"
)

func main() {
    client, err := redisclient.NewRedisClientFromURL("redis://localhost:6379")
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()

    // SET com expiração de 1 hora
    err = client.Set(ctx, "key", "value", time.Hour)
    if err != nil {
        log.Fatalf("Erro ao definir valor: %v", err)
    }

    // GET
    val, err := client.Get(ctx, "key")
    if err != nil {
        log.Fatalf("Erro ao obter valor: %v", err)
    }
    log.Printf("Valor obtido: %s", val)
}
```

## Dicas e Integração

- Use a interface `IRedisClient` para facilitar testes e mocks.
- Integre facilmente com middlewares de cache HTTP.
- Suporta múltiplos modos de conexão sem alterar código de uso.
- Para ambientes com autenticação, inclua usuário/senha na URL: `redis://:senha@host:6379`.

## API

Veja a documentação completa em [pkg.go.dev](https://pkg.go.dev/gitlab.globoi.com/globoplay/go-prime/clients/redisclient).

## Contribuindo

Contribuições são bem-vindas! Abra issues ou pull requests.
