# go-prime

Biblioteca modular em Go para construção de aplicações web e integrações de alta performance, com foco em produtividade, extensibilidade e boas práticas.

## O que é

O `go-prime` é um conjunto de pacotes Go para facilitar a criação de servidores HTTP, clientes HTTP robustos e integração com Redis, fornecendo middlewares prontos, interfaces simples e documentação clara.

## O que a lib entrega

- **server/**: Servidor HTTP baseado em Fiber, com middlewares para forwarding de headers, controle de cache, healthcheck e fácil extensibilidade.
- **clients/httpclient/**: Cliente HTTP extensível, com suporte a middlewares (logging, headers, cache, circuit breaker), base URL, timeout e todos os métodos HTTP.
- **clients/redisclient/**: Cliente Redis pronto para uso em cache, filas e integrações, com suporte a Standalone, Cluster e Sentinel.

## Documentação dos módulos

- [server/README.md](server/README.md): Como criar e configurar servidores HTTP, middlewares e healthcheck.
- [clients/httpclient/README.md](clients/httpclient/README.md): Como usar o cliente HTTP, middlewares, exemplos de requisições e dicas de integração.
- [clients/redisclient/README.md](clients/redisclient/README.md): Como configurar e usar o cliente Redis em diferentes modos.

## Instalação

Adicione o módulo ao seu projeto:

```bash
go get gitlab.globoi.com/globoplay/go-prime
```

## Exemplos rápidos

Veja exemplos completos e detalhados nos READMEs de cada módulo.

## Contribuindo

Contribuições são bem-vindas! Abra issues ou pull requests.

## Licença

MIT
