# Fin News API

## Versao - 1.0.0

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

## Criar img

```bash
docker buildx build --platform=linux/arm64/v8 -t maykealisson/fin-news:{{VERSION}} .
```

```bash
docker push maykealisson/fin-news:{{VERSION}}
```
