# httpclient

[![Go Reference](https://pkg.go.dev/badge/gitlab.globoi.com/globoplay/go-prime/clients/httpclient.svg)](https://pkg.go.dev/gitlab.globoi.com/globoplay/go-prime/clients/httpclient)

Cliente HTTP extensível para Go, com suporte a middlewares (logging, headers, cache, circuit breaker), base URL, timeout e métodos HTTP completos.

## Instalação

```bash
go get gitlab.globoi.com/globoplay/go-prime/clients/httpclient
```

## Visão Geral

- Base URL e timeout configuráveis
- Headers globais
- Middlewares plugáveis
- Métodos HTTP: GET, POST, PUT, PATCH, DELETE, HEAD

## Exemplo Rápido

```go
import (
    "context"
    "time"
    "gitlab.globoi.com/globoplay/go-prime/clients/httpclient"
)

client := httpclient.NewHTTPClient(
    "https://api.example.com",
    5*time.Second,
)

resp, err := client.Get(context.Background(), "/users/1")
```

## Middlewares

O httpclient suporta middlewares para customizar o comportamento das requisições. Você pode combiná-los conforme a necessidade do seu projeto.

### Logging Middleware

Registra todas as requisições e respostas, incluindo método, URL, status, duração e cache. Útil para auditoria e troubleshooting.

**Configuração:**

```go
client := httpclient.NewHTTPClient(
    baseURL,
    5*time.Second,
    httpclient.NewLoggingMiddleware("my-service"),
)
```

### Header Middleware

Adiciona ou sobrescreve headers em todas as requisições. Ideal para autenticação, rastreamento e customização de chamadas.

**Configuração:**

```go
client := httpclient.NewHTTPClient(
    baseURL,
    5*time.Second,
    httpclient.NewHeaderMiddleware(map[string]string{
        "Authorization": "Bearer token",
        "X-Request-ID": "123",
    }),
)
```

### Cache Middleware

Cacheia respostas de GET usando Redis. Reduz latência e carga em serviços externos.

**Configuração:**

```go
import "gitlab.globoi.com/globoplay/go-prime/clients/redisclient"
redis := redisclient.NewRedisClient()
cfg := &httpclient.CacheConfig{
    RedisClient: redis,
    TTL:         30 * time.Second,
    OverrideTTL: true,
    Headers:     []string{"Authorization"},
}
client := httpclient.NewHTTPClient(baseURL, 5*time.Second, httpclient.CacheMiddleware(cfg))
```

### Circuit Breaker Middleware

Protege contra falhas em serviços externos, abrindo o circuito após muitos erros. Evita sobrecarga e melhora a resiliência.

**Configuração:**

```go
client := httpclient.NewHTTPClient(
    baseURL,
    5*time.Second,
    httpclient.NewCircuitBreakerMiddleware("my-service"),
)
```

### Ordem recomendada dos middlewares

1. Logging
2. Headers
3. Cache
4. Circuit Breaker

```go
client := httpclient.NewHTTPClient(
    baseURL,
    5*time.Second,
    httpclient.NewLoggingMiddleware("my-service"),
    httpclient.NewHeaderMiddleware(map[string]string{"Authorization": "Bearer token"}),
    httpclient.CacheMiddleware(cfg),
    httpclient.NewCircuitBreakerMiddleware("my-service"),
)
```

## Métodos Disponíveis

| Método   | Descrição                                 |
|----------|-------------------------------------------|
| Get      | Requisição GET                            |
| Post     | Requisição POST com body                  |
| Put      | Requisição PUT com body                   |
| Patch    | Requisição PATCH com body                 |
| Delete   | Requisição DELETE                         |
| Head     | Requisição HEAD                           |

Todos os métodos recebem `context.Context` e retornam `*HTTPResponse` e `error`.

## Exemplos de Uso

```go
import (
    "context"
    "strings"
    "log"
)

ctx := context.Background()

// GET
resp, err := client.Get(ctx, "/users/1")
if err != nil { log.Fatal(err) }
log.Println("GET:", resp.StatusCode, resp.Body)

// POST
body := strings.NewReader(`{"name":"Luiz"}`)
resp, err = client.Post(ctx, "/users", body)
if err != nil { log.Fatal(err) }
log.Println("POST:", resp.StatusCode, resp.Body)

// PUT
body = strings.NewReader(`{"name":"Novo Nome"}`)
resp, err = client.Put(ctx, "/users/1", body)
if err != nil { log.Fatal(err) }
log.Println("PUT:", resp.StatusCode, resp.Body)

// PATCH
body = strings.NewReader(`{"name":"Parcial"}`)
resp, err = client.Patch(ctx, "/users/1", body)
if err != nil { log.Fatal(err) }
log.Println("PATCH:", resp.StatusCode, resp.Body)

// DELETE
resp, err = client.Delete(ctx, "/users/1")
if err != nil { log.Fatal(err) }
log.Println("DELETE:", resp.StatusCode)

// HEAD
resp, err = client.Head(ctx, "/users/1")
if err != nil { log.Fatal(err) }
log.Println("HEAD:", resp.StatusCode)
```

## Sequência de execução dos middlewares

1. Logging Middleware
2. Header Middleware
3. Cache Middleware
4. Circuit Breaker Middleware

## API

Veja a documentação completa em [pkg.go.dev](https://pkg.go.dev/gitlab.globoi.com/globoplay/go-prime/clients/httpclient).

## Contribuindo

Contribuições são bem-vindas! Abra issues ou pull requests.

## Licença

MIT
