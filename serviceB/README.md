# Weather API - Go + Cloud Run

[![Tests](https://github.com/adalbertofjr/lab-1-go-weather-cloud-run/actions/workflows/go-weather-cloud-run-tests.yml/badge.svg)](https://github.com/adalbertofjr/lab-1-go-weather-cloud-run/actions/workflows/go-weather-cloud-run-tests.yml)
[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/)
[![Coverage](https://img.shields.io/badge/coverage-90%25-brightgreen)](./coverage)

API REST em Go para consulta de temperatura por CEP, integrando ViaCEP e WeatherAPI. Desenvolvida com Clean Architecture e pronta para deploy no Google Cloud Run.

## 📋 Índice

1. [Quick Start](#1-quick-start)
2. [Tecnologias](#2-tecnologias)
3. [Arquitetura](#3-arquitetura)
4. [Pré-requisitos](#4-pré-requisitos)
5. [Configuração](#5-configuração)
6. [Executando o Projeto](#6-executando-o-projeto)
7. [Executando os Testes](#7-executando-os-testes)
8. [API Endpoints](#8-api-endpoints)
9. [Estrutura do Projeto](#9-estrutura-do-projeto)
10. [CI/CD](#10-cicd)
11. [Docker](#11-docker)
12. [Desenvolvimento](#12-desenvolvimento)

## 1. ⚡ Quick Start

### 🌐 Usar API em Produção (Google Cloud Run)

A API já está deployada e disponível para uso imediato:

```bash
# Testar CEP válido
curl "https://lab-1-go-weather-cloud-run-1080779949140.us-central1.run.app/?cep=01001000"

# Resposta esperada:
# {"localidade":"Sao Paulo","temp_c":20.2,"temp_f":68.36,"temp_k":293.2}

# Health check
curl "https://lab-1-go-weather-cloud-run-1080779949140.us-central1.run.app/health"
```

### 💻 Executar Localmente

```bash
# 1. Clonar repositório
git clone https://github.com/adalbertofjr/lab-1-go-weather-cloud-run.git
cd lab-1-go-weather-cloud-run

# 2. Configurar variáveis de ambiente
cd cmd/server
cp .env.example .env
# Edite .env e adicione sua WEATHERAPI_KEY

# 3. Executar aplicação
go run main.go
# Acesse: http://localhost:8000

# 4. Testar (em outro terminal)
curl "http://localhost:8000/?cep=01001000"

# 5. Executar testes localmente
go test -v ./...

# 6. Executar testes via Docker
make test-docker
# ou: docker compose -f docker-compose.test.yml run --rm test
```

## 2. 🚀 Tecnologias

- **Go 1.23** - Linguagem de programação
- **Chi Router** - HTTP router leve e rápido
- **Viper** - Gerenciamento de configurações
- **Docker** - Containerização
- **GitHub Actions** - CI/CD
- **Google Cloud Run** - Deploy (serverless)

### APIs Externas

- [ViaCEP](https://viacep.com.br/) - Consulta de CEP
- [WeatherAPI](https://www.weatherapi.com/) - Dados meteorológicos

## 3. 🏗️ Arquitetura

Projeto estruturado seguindo **Clean Architecture**:

```
┌─────────────┐
│   Handler   │  HTTP (Chi Router)
└──────┬──────┘
       │
┌──────▼──────┐
│   UseCase   │  Regras de Negócio
└──────┬──────┘
       │
┌──────▼──────┐
│   Gateway   │  Integrações Externas (ViaCEP, WeatherAPI)
└──────┬──────┘
       │
┌──────▼──────┐
│   Entity    │  Modelos de Domínio
└─────────────┘
```

**Camadas:**
- **Domain** - Entidades e interfaces de negócio
- **UseCase** - Lógica de aplicação e orquestração
- **Infrastructure** - Handlers HTTP, Gateways, Web Server
- **Main** - Configuração e inicialização

## 4. ✅ Pré-requisitos

- [Go 1.23+](https://go.dev/dl/)
- [Docker](https://www.docker.com/get-started) (opcional, para testes)
- Chave API do [WeatherAPI](https://www.weatherapi.com/signup.aspx) (gratuita)

## 5. ⚙️ Configuração

```bash
git clone https://github.com/adalbertofjr/lab-1-go-weather-cloud-run.git
cd lab-1-go-weather-cloud-run
```

### 2. Configure as variáveis de ambiente

```bash
cd cmd/server
cp .env.example .env
```

Edite o arquivo `.env` e adicione sua chave da WeatherAPI:

```env
WEATHERAPI_KEY=sua_chave_aqui
WEB_SERVER_PORT=:8000
```

> 💡 **Obtenha sua chave gratuita:** https://www.weatherapi.com/signup.aspx

### 3. Instale as dependências

```bash
# Na raiz do projeto
go mod download
```

## 5. ⚙️ Configuração

### Opção 1: Execução Local

```bash
cd cmd/server
go run main.go
```

O servidor estará disponível em: **http://localhost:8000**

### Opção 2: Com Docker (recomendado para produção)

```bash
# Build da imagem (multi-stage build otimizado)
docker build -t weather-api .

# Executar container
docker run --rm -p 8080:8080 \
  -e WEATHERAPI_KEY=sua_chave_aqui \
  -e WEB_SERVER_PORT=:8080 \
  weather-api
```

**Características do Dockerfile:**
- 🏗️ **Multi-stage build** (builder + runtime)
- 📦 **Imagem final ~15MB** (Alpine + binário estático)
- 🔒 **CGO_ENABLED=0** - binário 100% estático
- 🔐 **Certificados SSL** incluídos (ca-certificates)
- ⚡ **Otimizado para Cloud Run**

### Testando a API

```bash
# Health check
curl http://localhost:8080/health

# Consultar temperatura por CEP
curl "http://localhost:8080/?cep=01001000"
```

## 6. 🎯 Executando o Projeto

O projeto possui **37 testes** com **90%+ de cobertura** nas camadas críticas.

### Opção 1: Testes Locais (rápido)

```bash
# Executar todos os testes
go test -v ./...

# Com cobertura
go test -v -cover ./...

# Gerar relatório HTML de cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # macOS
```

### Opção 2: Testes com Docker (ambiente isolado)

#### Usando script automatizado (mais fácil):

```bash
./test-docker.sh test      # Executar testes
./test-docker.sh coverage  # Gerar relatório HTML
./test-docker.sh shell     # Shell interativo
./test-docker.sh clean     # Limpar containers
```

#### Usando Makefile:

```bash
make help                  # Ver todos os comandos
make test-local            # Testes locais
make test-docker           # Testes no Docker
make test-docker-coverage  # Cobertura no Docker
make docker-clean          # Limpar containers
```

#### Usando Docker Compose:

```bash
docker compose -f docker-compose.test.yml run --rm test
docker compose -f docker-compose.test.yml run --rm test-coverage
```

### Cobertura por Camada

| Camada | Cobertura | Testes |
|--------|-----------|--------|
| Entity | 100% | 11 |
| UseCase | 100% | 8 |
| InternalError | 100% | 5 |
| Handler | 63.2% | 11 |
| Utility | 90% | 2 |

### Executar testes específicos

```bash
# Testar apenas UseCase
go test -v ./internal/usecase/weather/

# Testar apenas Entity
go test -v ./internal/domain/entity/

# Executar teste específico
go test -v -run TestGetCurrentWeather_Success ./...
```

## 7. 🧪 Executando os Testes

### `GET /health`
Verifica se a API está rodando.

**Resposta:**
```json
{
  "status": "OK"
}
```

### `GET /?cep={cep}`
Retorna a temperatura atual para o CEP informado.

**Parâmetros:**
- `cep` (string, obrigatório) - CEP com ou sem hífen (ex: `01001000` ou `01001-000`)

**Exemplo de Requisição:**
```bash
curl "http://localhost:8000/?cep=01001000"
```

**Resposta de Sucesso (200):**
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

**Erros:**

- **422 - CEP inválido:**
```json
{
  "message": "invalid zipcode"
}
```

- **404 - CEP não encontrado:**
```json
{
  "message": "can not find zipcode"
}
```

## 8. 📡 API Endpoints

```
.
├── cmd/
│   ├── configs/          # Configurações (Viper)
│   └── server/
│       ├── main.go       # Entry point
│       └── .env.example  # Exemplo de variáveis de ambiente
├── internal/
│   ├── domain/
│   │   ├── entity/       # Weather (conversões de temperatura)
│   │   └── gateway/      # Interface WeatherGateway
│   ├── usecase/
│   │   └── weather/      # GetCurrentWeather (lógica de negócio)
│   └── infra/
│       ├── api/          # Handlers HTTP + DTOs
│       ├── gateway/      # WeatherAPI (implementação)
│       ├── internal_error/ # Erros customizados
│       └── web/          # WebServer (Chi)
├── pkg/
│   ├── net/              # HTTP connection utilities
│   └── utility/          # CEP validator/formatter
├── .github/
│   └── workflows/
│       └── test.yml      # CI/CD (GitHub Actions)
├── Dockerfile.test       # Imagem Docker para testes
├── docker-compose.test.yml # Orquestração de testes
├── Makefile              # Comandos simplificados
├── test-docker.sh        # Script automatizado de testes
└── go.mod                # Dependências
```

## 9. 📁 Estrutura do Projeto

O projeto usa **GitHub Actions** para executar testes automaticamente em cada push/PR.

### Workflow: `.github/workflows/go-weather-cloud-run-tests.yml`

```yaml
name: Go Weather Cloud Run - Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: docker compose -f docker-compose.test.yml run --rm test
```

**Vantagens:**
- ✅ Ambiente isolado e reproduzível
- ✅ Mesma versão Go (1.23) em qualquer lugar
- ✅ Sem necessidade de configurar Go no runner
- ✅ Cache automático de dependências

Ver status dos testes: [Actions](https://github.com/adalbertofjr/lab-1-go-weather-cloud-run/actions)

## 10. 🔄 CI/CD

### Arquivos Docker

- **`Dockerfile`** - Imagem de produção (multi-stage build, ~15MB)
- **`Dockerfile.test`** - Imagem Alpine otimizada para testes
- **`docker-compose.test.yml`** - 3 serviços (test, test-coverage, test-watch)
- **`.dockerignore`** - Otimização de build

### Dockerfile de Produção

O `Dockerfile` usa **multi-stage build** para criar imagem extremamente otimizada:

**Stage 1 - Builder:**
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather-api ./cmd/server
```

**Stage 2 - Runtime:**
```dockerfile
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/weather-api .
ENTRYPOINT ["/app/weather-api"]
```

**Resultado:**
- 📦 Imagem final: **~15MB** (vs ~300MB sem otimização)
- 🔒 Binário estático (CGO_ENABLED=0)
- 🔐 Certificados SSL para APIs externas
- ⚡ Cold start rápido no Cloud Run

### Comandos Docker

```bash
# Produção
docker build -t weather-api .
docker run --rm -p 8080:8080 \
  -e WEATHERAPI_KEY=sua_chave \
  -e WEB_SERVER_PORT=:8080 \
  weather-api

# Testes
docker compose -f docker-compose.test.yml run --rm test

# Shell interativo
docker compose -f docker-compose.test.yml run --rm test-watch

# Limpar
docker compose -f docker-compose.test.yml down --rmi all
```

Documentação completa: [DOCKER_TESTS.md](./DOCKER_TESTS.md)

## 12. 🛠️ Desenvolvimento

- [DOCKER_TESTS.md](./DOCKER_TESTS.md) - Guia completo de testes com Docker
- [QUICK_START_DOCKER.md](./QUICK_START_DOCKER.md) - Quick start para testes

## 11. 🐳 Docker

### Adicionar novos testes

```bash
# Criar arquivo de teste
touch internal/domain/entity/new_test.go

# Executar apenas esse teste
go test -v ./internal/domain/entity/ -run TestNew
```

### Validar antes de commit

```bash
# Executar todos os testes
make test-local

# Verificar cobertura
make test-coverage
```

## 📝 Licença

Este projeto é parte de um laboratório de estudos de Pós-Graduação em Go.

---

**Autor:** Adalberto F. Jr.  
**Repositório:** https://github.com/adalbertofjr/lab-1-go-weather-cloud-run
