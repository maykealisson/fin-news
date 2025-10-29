# Fin News API

## Versao - 0.0.1

## Configuração

```bash
# Instalar dependências
go mod download

# Configurar modo de ambiente (development/production)
export GIN_MODE=release  # Para produção
export GIN_MODE=debug   # Para desenvolvimento
```

## Rodando o projeto

```bash
# Modo desenvolvimento
go run main.go

# Modo produção
GIN_MODE=release go run main.go
```

## Endpoints

### GET /noticias

Busca notícias relacionadas a um ativo financeiro.

Query Parameters:

- ativo: Código do ativo (ex: PETR4)

```bash
curl "http://localhost:3001/noticias?ativo=PETR4"
```
