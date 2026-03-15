# Boilerplate Golang

API base em Go usando Fiber, GORM, PostgreSQL e migrations com `golang-migrate`.

## Tecnologias

- Go
- Fiber
- GORM
- PostgreSQL (Docker)
- golang-migrate

## Pré-requisitos

- Go instalado (projeto usa Go `1.26.1` no `go.mod`)
- Docker e Docker Compose instalados
- `make` instalado

## Passo a Passo Para Rodar o Projeto

### 1. Configurar variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto com o conteúdo listado no arquivo `.env.example.org`.

Observação: o app atualmente sobe na porta `3000` (definida no código), então mantenha a porta livre.

### 2. Subir o banco de dados

```bash
docker compose up -d
```

Para verificar se o container subiu:

```bash
docker ps
```

### 3. Rodar a aplicação

```bash
make run
```

Na inicialização, as migrations são aplicadas automaticamente.

### 4. Testar se está funcionando

Health check:

```bash
curl http://localhost:3000/health
```

Resposta esperada:

```json
{ "status": "ok" }
```

Documentação Swagger (UI):

- `http://localhost:3000/docs`

## Comandos Úteis

Rodar testes:

```bash
make test
```

Build local:

```bash
make build
```

Parar e remover containers:

```bash
docker compose down
```

Parar e remover containers + volume do banco (reset completo):

```bash
docker compose down -v
```

## Troubleshooting Rápido

Erro de migration com banco em estado `dirty`:

1. Derrube e remova volume:

```bash
docker compose down -v
```

2. Suba novamente:

```bash
docker compose up -d
```

3. Rode a API de novo:

```bash
make run
```
