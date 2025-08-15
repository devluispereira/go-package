# server

[![Go Reference](https://pkg.go.dev/badge/gitlab.globoi.com/globoplay/go-prime/server.svg)](https://pkg.go.dev/gitlab.globoi.com/globoplay/go-prime/server)

Servidor HTTP modular baseado em [Fiber](https://gofiber.io/), com middlewares para forwarding de headers, controle de cache, cabeçalhos customizados e healthcheck.

## Instalação

```bash
go get gitlab.globoi.com/globoplay/go-prime/server
```

## Visão Geral

- Remove cabeçalhos padrão de identificação do servidor
- Adiciona cabeçalho `X-Origin-App` para rastreio de origem
- Middleware para forwarding de headers customizáveis
- Middleware para controle de cache HTTP
- Endpoint `/healthcheck` pronto para uso

## Exemplo Rápido

```go
import (
    "log"
    "gitlab.globoi.com/globoplay/go-prime/server"
)

func main() {
    srv := server.NewServer("my-app", []string{"x-request-id", "x-client-user-agent"})
    log.Fatal(srv.App.Listen(":8080"))
}
```

## Middlewares

### ForwardHeadersMiddleware

Coleta e encaminha headers HTTP de interesse para handlers e serviços downstream. Útil para rastreamento, autenticação e contexto de requisição.

- Por padrão, encaminha headers comuns de tracing e identificação.
- Permite customizar a lista de headers.
- Adiciona sempre o header `x-origin-app`.

**Configuração:**

```go
app.Use(server.ForwardHeadersMiddleware("my-app", []string{"x-request-id", "x-client-user-agent"}))
```

**Acessando headers encaminhados:**

```go
headers := c.Locals("forwardedHeaders").(map[string]string)
```

### SetCacheControlMiddleware

Define o header `Cache-Control` para rotas ou grupos, facilitando o controle de cache HTTP.

- Suporta os tipos: `public`, `private`, `no-store`, `no-cache`.
- Permite definir TTL (max-age) em segundos.

**Configuração:**

```go
app.Get("/public", server.SetCacheControlMiddleware(server.CachePublic, 60), handler)
app.Get("/private", server.SetCacheControlMiddleware(server.CachePrivate, 0), handler)
```

## Endpoint de Healthcheck

O servidor já expõe o endpoint `/healthcheck` para monitoramento:

```bash
curl http://localhost:8080/healthcheck
# Resposta: OK
```

## Estrutura e Extensibilidade

- O tipo `Server` expõe o campo `App` para customização avançada com Fiber.
- Adicione middlewares, rotas e handlers conforme sua necessidade.
- Integre facilmente com middlewares de autenticação, logging, cache, etc.

## API

Veja a documentação completa em [pkg.go.dev](https://pkg.go.dev/gitlab.globoi.com/globoplay/go-prime/server).

## Contribuindo

Contribuições são bem-vindas! Abra issues ou pull requests.

## Licença

MIT
