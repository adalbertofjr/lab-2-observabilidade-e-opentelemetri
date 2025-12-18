# 🌡️ Sistema de Temperatura por CEP - Go + OTEL + Zipkin

[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/)
[![OpenTelemetry](https://img.shields.io/badge/OpenTelemetry-Enabled-blue)](https://opentelemetry.io/)
[![Zipkin](https://img.shields.io/badge/Zipkin-Tracing-orange)](https://zipkin.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

Sistema distribuído em Go que recebe um CEP, identifica a cidade e retorna o clima atual (temperatura em Celsius, Fahrenheit e Kelvin) com **tracing distribuído** implementado usando **OpenTelemetry** e **Zipkin**.

---

## 📋 Índice

1. [Quick Start](#-quick-start)
   - [Docker Compose (Recomendado)](#1-docker-compose-recomendado)
   - [Execução Local](#2-execução-local)
2. [Sobre o Projeto](#-sobre-o-projeto)
3. [Arquitetura](#-arquitetura)
4. [Tecnologias](#-tecnologias)
5. [Pré-requisitos](#-pré-requisitos)
6. [Configuração](#-configuração)
7. [API Endpoints](#-api-endpoints)
8. [Tracing Distribuído](#-tracing-distribuído)
9. [Estrutura do Projeto](#-estrutura-do-projeto)
10. [Testes](#-testes)
11. [Troubleshooting](#-troubleshooting)

---

## ⚡ Quick Start

### 1. Docker Compose (Recomendado)

A forma mais rápida de executar todo o sistema com OTEL e Zipkin:

```bash
# 1. Clone o repositório
git clone https://github.com/adalbertofjr/lab-2-observabilidade-e-opentelemetri.git
cd lab-2-observabilidade-e-opentelemetri

# 2. Configure a API Key do WeatherAPI (OBRIGATÓRIO)
# Obtenha sua chave gratuita em: https://www.weatherapi.com/signup.aspx
# Edite o docker-compose.yaml e substitua 'your_api_key_here' pela sua chave:
nano docker-compose.yaml  # ou vim, code, etc.
# Linha 36: WEATHERAPI_KEY=sua_chave_aqui

# 3. Inicie todos os serviços
docker-compose up -d

# 4. Aguarde os serviços iniciarem (~30 segundos)
docker-compose ps

# 5. Teste o sistema
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"cep": "01001000"}'

# 6. Acesse o Zipkin para visualizar os traces
open http://localhost:9411
```

**Serviços disponíveis:**
- 🔵 **Serviço A** (Input): http://localhost:8080
- 🟢 **Serviço B** (Orquestração): http://localhost:8000
- 🟠 **Zipkin UI**: http://localhost:9411
- 🔴 **OTEL Collector**: http://localhost:4317 (gRPC)

---

### 2. Execução Local

Para desenvolvimento local, execute cada serviço manualmente:

#### Passo 1: Clone o repositório

```bash
git clone https://github.com/adalbertofjr/lab-2-observabilidade-e-opentelemetri.git
cd lab-2-observabilidade-e-opentelemetri
```

#### Passo 2: Inicie a infraestrutura (OTEL Collector e Zipkin)

```bash
# Na raiz do projeto, inicie apenas Zipkin e OTEL Collector
docker-compose up zipkin-all-in-one otel-collector -d
```

#### Passo 3: Configure o Serviço B

```bash
# Entre no diretório do Serviço B
cd serviceB/cmd/server

# Crie o arquivo .env a partir do exemplo
cp .env.example .env

# Edite o arquivo .env
nano .env  # ou vim, code, etc.
```

**Adicione sua WeatherAPI Key no arquivo .env:**
```env
WEATHERAPI_KEY=sua_chave_aqui_sem_aspas
WEB_SERVER_PORT=:8000
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
```

> 💡 Obtenha sua chave gratuita em: https://www.weatherapi.com/signup.aspx

```bash
# Baixe as dependências do Go
go mod download
```

#### Passo 4: Inicie o Serviço B (Terminal 1)

```bash
# A partir do diretório serviceB/cmd/server
go run main.go
```

O Serviço B estará rodando em: http://localhost:8000

#### Passo 5: Inicie o Serviço A (Terminal 2)

```bash
# Volte para a raiz do projeto
cd ../../..

# Entre no diretório do Serviço A
cd serviceA/cmd/server

# Baixe as dependências do Serviço A
go mod download

# Execute o serviço (não precisa configurar .env)
go run main.go
```

> 💡 **Nota**: O Serviço A já possui valores padrão (`localhost:4317` e `http://localhost:8000`).  
> Só exporte variáveis se precisar customizar:
> ```bash
> export OTEL_EXPORTER_OTLP_ENDPOINT=outro_endpoint:4317
> export SERVICE_B_URL=http://outro_host:8000
> ```

O Serviço A estará rodando em: http://localhost:8080

#### Passo 6: Teste o sistema

```bash
# Teste o fluxo completo através do Serviço A
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"cep": "01001000"}'

# Ou teste o Serviço B diretamente
curl "http://localhost:8000/?cep=01001000"

# Acesse o Zipkin para visualizar os traces
open http://localhost:9411
```

---

## 📖 Sobre o Projeto

Este projeto foi desenvolvido como parte de um laboratório sobre **Observabilidade e OpenTelemetry** em sistemas distribuídos. O objetivo é demonstrar a implementação de **tracing distribuído** em uma arquitetura de microserviços usando Go.

### 🎯 Objetivos

- ✅ Implementar comunicação entre microserviços
- ✅ Validar entrada de dados (CEP)
- ✅ Integrar com APIs externas (ViaCEP e WeatherAPI)
- ✅ Converter temperaturas (Celsius → Fahrenheit, Kelvin)
- ✅ Implementar observabilidade com OpenTelemetry
- ✅ Visualizar traces distribuídos com Zipkin

### 🔍 O que o sistema faz?

1. **Serviço A** recebe um CEP via POST
2. Valida o formato do CEP (8 dígitos numéricos)
3. Encaminha para o **Serviço B**
4. **Serviço B** busca a localização (ViaCEP)
5. **Serviço B** busca a temperatura atual (WeatherAPI)
6. Converte temperatura para Celsius, Fahrenheit e Kelvin
7. Retorna os dados formatados
8. Todo o fluxo é rastreado com **spans distribuídos**

---

## 🏗️ Arquitetura

### Diagrama de Comunicação

```
┌──────────────┐     POST /cep      ┌──────────────┐
│   Cliente    │ ─────────────────> │  Serviço A   │
│              │     (JSON)         │   :8080      │
└──────────────┘                    └───────┬──────┘
                                            │
                                            │ GET /?cep=xxx
                                            │
                                    ┌───────▼──────┐
                                    │  Serviço B   │
                                    │   :8000      │
                                    └───────┬──────┘
                                            │
                        ┌───────────────────┼───────────────────┐
                        │                   │                   │
                  ┌─────▼──────┐     ┌─────▼──────┐     ┌─────▼──────┐
                  │  ViaCEP    │     │ WeatherAPI │     │   Entity   │
                  │ (Location) │     │  (Temp)    │     │ (Conversão)│
                  └────────────┘     └────────────┘     └────────────┘
                        │                   │                   │
                        └───────────────────┴───────────────────┘
                                            │
                        ┌───────────────────▼───────────────────┐
                        │      OTEL Collector (gRPC)            │
                        │            :4317                      │
                        └───────────────────┬───────────────────┘
                                            │
                        ┌───────────────────▼───────────────────┐
                        │         Zipkin (UI + Storage)         │
                        │            :9411                      │
                        └───────────────────────────────────────┘
```

### Serviço A - Input e Validação

**Responsabilidades:**
- Receber requisições POST com CEP
- Validar formato do CEP (8 dígitos)
- Encaminhar para Serviço B via HTTP
- Propagar contexto de tracing

**Stack Técnico:**
- **Chi Router** - HTTP routing
- **OpenTelemetry** - Instrumentação
- **Clean Architecture** - Organização de código

**Estrutura:**
```
serviceA/
├── cmd/server/main.go          # Inicialização + OTEL setup
├── internal/
│   ├── domain/
│   │   ├── entity/             # Weather entity
│   │   └── gateway/            # Interface para Serviço B
│   ├── usecase/weather/        # Validação + chamada ao Serviço B
│   └── infra/
│       ├── api/                # HTTP handler
│       └── gateway/            # Cliente HTTP para Serviço B
└── pkg/utility/                # Validador de CEP
```

### Serviço B - Orquestração e APIs Externas

**Responsabilidades:**
- Receber CEP do Serviço A
- Buscar localização no ViaCEP
- Buscar temperatura no WeatherAPI
- Converter temperaturas (F, K)
- Retornar dados formatados

**Stack Técnico:**
- **Chi Router** - HTTP routing
- **Viper** - Configuração
- **OpenTelemetry** - Instrumentação
- **Clean Architecture** - Organização de código

**Estrutura:**
```
serviceB/
├── cmd/
│   ├── configs/                # Viper configuration
│   └── server/main.go          # Inicialização + OTEL setup
├── internal/
│   ├── domain/
│   │   ├── entity/             # Weather (com conversões)
│   │   └── gateway/            # Interface para APIs
│   ├── usecase/weather/        # Orquestração da busca
│   └── infra/
│       ├── api/                # HTTP handlers
│       ├── gateway/            # Clientes ViaCEP e WeatherAPI
│       ├── internal_error/     # Erros customizados (422, 404)
│       └── web/                # WebServer
└── pkg/utility/                # Validador de CEP
```

### Clean Architecture

Ambos os serviços seguem **Clean Architecture**:

```
┌─────────────────────────────────────────────────┐
│              Handlers (HTTP)                    │  ← Camada Externa
├─────────────────────────────────────────────────┤
│              UseCases                           │  ← Lógica de Aplicação
├─────────────────────────────────────────────────┤
│              Gateways (Implementações)          │  ← Adaptadores
├─────────────────────────────────────────────────┤
│         Domain (Entities + Interfaces)          │  ← Núcleo do Negócio
└─────────────────────────────────────────────────┘
```

**Benefícios:**
- ✅ Baixo acoplamento
- ✅ Testabilidade alta
- ✅ Facilidade de manutenção
- ✅ Independência de frameworks

---

## 🛠️ Tecnologias

### Backend
- **[Go 1.23](https://go.dev/)** - Linguagem de programação
- **[Chi Router](https://github.com/go-chi/chi)** - HTTP router leve e rápido
- **[Viper](https://github.com/spf13/viper)** - Gerenciamento de configurações

### Observabilidade
- **[OpenTelemetry](https://opentelemetry.io/)** - Instrumentação de tracing
- **[OTEL Collector](https://opentelemetry.io/docs/collector/)** - Agregação de traces
- **[Zipkin](https://zipkin.io/)** - Visualização de traces distribuídos

### APIs Externas
- **[ViaCEP](https://viacep.com.br/)** - Consulta de CEP (gratuita)
- **[WeatherAPI](https://www.weatherapi.com/)** - Dados meteorológicos (gratuita)

### Infraestrutura
- **[Docker](https://www.docker.com/)** - Containerização
- **[Docker Compose](https://docs.docker.com/compose/)** - Orquestração local

---

## ✅ Pré-requisitos

### Para Docker Compose (Recomendado)
- [Docker](https://www.docker.com/get-started) 20.10+
- [Docker Compose](https://docs.docker.com/compose/install/) 2.0+
- **Chave API do [WeatherAPI](https://www.weatherapi.com/signup.aspx)** (gratuita - OBRIGATÓRIA)

### Para Execução Local
- [Go 1.23+](https://go.dev/dl/)
- [Docker](https://www.docker.com/get-started) (para OTEL Collector e Zipkin)
- **Chave API do [WeatherAPI](https://www.weatherapi.com/signup.aspx)** (gratuita - OBRIGATÓRIA)

**Verificar instalação do Go:**
```bash
go version  # Deve retornar go1.23 ou superior
```

---

## ⚙️ Configuração

### 1. Clone o repositório

```bash
git clone https://github.com/adalbertofjr/lab-2-observabilidade-e-opentelemetri.git
cd lab-2-observabilidade-e-opentelemetri
```

### 2. Configure a WeatherAPI Key (OBRIGATÓRIO)

⚠️ **O projeto NÃO possui chave padrão.** Você precisa criar sua própria chave gratuita:

1. **Acesse:** https://www.weatherapi.com/signup.aspx
2. **Crie uma conta gratuita** (não precisa cartão de crédito)
3. **Copie sua API key** do dashboard

**Para Docker Compose:**

**Opção A - Editar docker-compose.yaml** (recomendado - mais confiável):
```bash
# Edite o arquivo docker-compose.yaml
nano docker-compose.yaml  # ou vim, code, etc.

# Encontre a linha 36 e substitua:
# - WEATHERAPI_KEY=${WEATHERAPI_KEY:-your_api_key_here}
# Por:
# - WEATHERAPI_KEY=sua_chave_aqui

# Salve e inicie os serviços
docker-compose up -d
```

**Opção B - Exportar variável de ambiente** (pode não funcionar em alguns ambientes):
```bash
export WEATHERAPI_KEY=sua_chave_aqui
docker-compose up -d
```

> ⚠️ **Importante**: Se usar Opção A, não faça commit do arquivo com a chave. Adicione ao .gitignore ou use git update-index --skip-worktree docker-compose.yaml

**Para execução local do Serviço B:**
```bash
cd serviceB/cmd/server
cp .env.example .env

# Edite o arquivo .env:
nano .env  # ou use vim, code, etc.
```

**Conteúdo do arquivo .env:**
```env
WEATHERAPI_KEY=sua_chave_aqui_sem_aspas
WEB_SERVER_PORT=:8000
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
```

> ✅ **Serviço A**: Não precisa de .env (usa valores padrão hardcoded)  
> ⚠️ **Serviço B**: Requer .env com WEATHERAPI_KEY obrigatória

### 3. Variáveis de Ambiente

#### Serviço A (porta 8080)

**Não requer arquivo .env** - todas as variáveis têm valores padrão adequados.

| Variável | Padrão | Descrição |
|----------|--------|-----------||
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `localhost:4317` | Endpoint do OTEL Collector |
| `SERVICE_B_URL` | `http://localhost:8000` | URL do Serviço B |

> 💡 **Dica**: O serviço funciona sem configuração adicional. Só exporte variáveis se precisar alterar os endpoints padrão.

#### Serviço B (porta 8000)

**Requer arquivo .env** com a WeatherAPI Key.

| Variável | Padrão | Descrição |
|----------|--------|-----------||
| `WEATHERAPI_KEY` | *(nenhum - obrigatório)* | **API key do WeatherAPI** - [Obtenha aqui](https://www.weatherapi.com/signup.aspx) |
| `WEB_SERVER_PORT` | `:8000` | Porta do servidor HTTP |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `otel-collector:4317` (Docker)<br>`localhost:4317` (local) | Endpoint do OTEL Collector |

---

## 📡 API Endpoints

### Serviço A - Input (Porta 8080)

#### `POST /`
Recebe um CEP e retorna a temperatura.

**Request:**
```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{"cep": "01001000"}'
```

**Request Body:**
```json
{
  "cep": "01001000"
}
```

**Response 200 - Sucesso:**
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

**Response 422 - CEP Inválido:**
```json
invalid zipcode
```

### Serviço B - Orquestração (Porta 8000)

#### `GET /?cep={cep}`
Processa o CEP e retorna temperatura.

**Request:**
```bash
curl "http://localhost:8000/?cep=01001000"
```

**Response 200 - Sucesso:**
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

**Response 422 - CEP Inválido:**
```
Invalid zipcode
```

**Response 404 - CEP Não Encontrado:**
```
Can not find zipcode
```

#### `GET /health`
Health check do serviço.

**Response 200:**
```json
{
  "status": "OK"
}
```

---

## 🔍 Tracing Distribuído

### OpenTelemetry + Zipkin

O projeto implementa **tracing distribuído** completo:

#### Spans Criados

**Serviço A:**
1. `POST /cep` - Handler principal
2. `validate_cep` - Validação do formato
3. `call_service_b` - Chamada HTTP ao Serviço B

**Serviço B:**
4. `Get /?cep=xxx` - Handler principal
5. `validate_cep` - Validação do formato
6. `fetch_weather_data` - Orquestração completa
7. `fetch_cep_location` - Chamada ao ViaCEP
8. `fetch_current_weather` - Chamada ao WeatherAPI

### Visualizando Traces no Zipkin

1. **Acesse o Zipkin UI**: http://localhost:9411

2. **Busque traces**:
   - Clique em "RUN QUERY"
   - Ou filtre por serviço: `ServiceA` ou `ServiceB`

3. **Analise o trace**:
   - Veja o tempo total da requisição
   - Identifique gargalos (qual span demorou mais)
   - Visualize a propagação do contexto entre serviços

### Exemplo de Trace

```
ServiceA: POST /cep (total: 245ms)
├─ validate_cep (2ms)
└─ call_service_b (243ms)
   └─ ServiceB: Get /?cep=xxx (240ms)
      ├─ validate_cep (1ms)
      └─ fetch_weather_data (239ms)
         ├─ fetch_cep_location (120ms) ← ViaCEP
         └─ fetch_current_weather (119ms) ← WeatherAPI
```

### Propagação de Contexto

O projeto usa **W3C Trace Context** para propagar o contexto:

**Serviço A injeta headers:**
```go
otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
```

**Serviço B extrai headers:**
```go
ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
```

---

## 📂 Estrutura do Projeto

### Árvore Completa

```
.
├── .docker/
│   └── otel-collector-config.yaml  # Configuração do OTEL Collector
├── serviceA/                       # Serviço de Input
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── entity/
│   │   │   └── gateway/
│   │   ├── usecase/weather/
│   │   └── infra/
│   │       ├── api/
│   │       └── gateway/
│   ├── pkg/utility/
│   ├── Dockerfile
│   └── go.mod
├── serviceB/                       # Serviço de Orquestração
│   ├── cmd/
│   │   ├── configs/
│   │   └── server/main.go
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── entity/
│   │   │   └── gateway/
│   │   ├── usecase/weather/
│   │   └── infra/
│   │       ├── api/
│   │       ├── gateway/
│   │       ├── internal_error/
│   │       └── web/
│   ├── pkg/utility/
│   ├── Dockerfile
│   └── go.mod
├── docker-compose.yaml             # Orquestração completa
├── REQUIREMENTS.md                 # Requisitos do projeto
├── AVALIACAO.md                    # Avaliação técnica
└── README.md                       # Este arquivo
```

---

## 🧪 Testes

### Serviço A

```bash
cd serviceA
go test -v ./...
```

**Cobertura:**
- Testes de validação de CEP
- Testes de integração com Serviço B
- Testes de propagação de contexto

### Serviço B

```bash
cd serviceB
go test -v ./...
```

**Cobertura:**
- Testes unitários (Entity, UseCase)
- Testes de integração (Handlers)
- Testes de erros customizados
- Testes de conversão de temperatura

### Executar todos os testes

```bash
# Serviço A
(cd serviceA && go test -v ./...)

# Serviço B
(cd serviceB && go test -v ./...)
```

---

## 🐛 Troubleshooting

### Problema: "Cannot connect to OTEL Collector"

**Solução:**
```bash
# Verifique se o OTEL Collector está rodando
docker-compose ps otel-collector

# Verifique os logs
docker-compose logs otel-collector

# Reinicie o serviço
docker-compose restart otel-collector
```

### Problema: "Zipkin não mostra traces"

**Soluções:**
1. Aguarde alguns segundos após fazer requisições
2. Verifique se os serviços estão enviando dados:
```bash
docker-compose logs otel-collector | grep -i trace
```
3. Acesse http://localhost:13133 (health check do collector)
4. Limpe o cache do Zipkin: http://localhost:9411

### Problema: "invalid zipcode" para CEP válido

**Causa:** CEP deve ter exatamente 8 dígitos numéricos.

**Soluções:**
```bash
# ✅ Correto
{"cep": "01001000"}
{"cep": "01001-000"}  # Aceita hífen

# ❌ Errado
{"cep": "1001000"}     # 7 dígitos
{"cep": "010010000"}   # 9 dígitos
{"cep": "ABC01000"}    # Letras
```

### Problema: "can not find zipcode"

**Causa:** CEP não existe na base do ViaCEP.

**Solução:** Use CEPs reais brasileiros. Exemplos:
- `01001000` - Praça da Sé, São Paulo - SP
- `20040020` - Centro, Rio de Janeiro - RJ
- `30130100` - Centro, Belo Horizonte - MG
- `40010000` - Comércio, Salvador - BA
- `80010000` - Centro, Curitiba - PR

### Problema: Serviços não conseguem se comunicar (Docker)

**Solução:**
```bash
# Verifique a rede Docker
docker network inspect lab-2-observabilidade-e-opentelemetri_default

# Reinicie os serviços com rebuild
docker-compose down
docker-compose up --build -d
```

### Problema: WeatherAPI retorna erro

**Causas comuns:**
1. **Chave não configurada** - A chave é obrigatória, não há chave padrão
2. **Chave inválida** - Verifique se copiou corretamente (sem espaços)
3. **Limite excedido** - Plano gratuito tem 1M requisições/mês
4. **Arquivo .env não encontrado** - Certifique-se de criar o .env no ServiceB

**Soluções:**

1. **Crie sua chave gratuita**: https://www.weatherapi.com/signup.aspx

2. **Para Docker Compose**:
   ```bash
   export WEATHERAPI_KEY=sua_chave_aqui
   docker-compose up -d
   ```

3. **Para execução local** (apenas ServiceB precisa):
   ```bash
   cd serviceB/cmd/server
   cp .env.example .env
   # Edite o .env e adicione: WEATHERAPI_KEY=sua_chave_aqui
   ```

4. **Teste a chave diretamente**:
   ```bash
   curl "https://api.weatherapi.com/v1/current.json?key=SUA_CHAVE&q=Sao%20Paulo&aqi=no"
   ```

5. **Verifique se a variável foi carregada**:
   ```bash
   docker-compose exec serviceB env | grep WEATHERAPI_KEY
   ```

### Logs úteis

```bash
# Ver logs de todos os serviços
docker-compose logs -f

# Ver logs de um serviço específico
docker-compose logs -f serviceA
docker-compose logs -f serviceB
docker-compose logs -f otel-collector
docker-compose logs -f zipkin-all-in-one

# Ver últimas 100 linhas
docker-compose logs --tail=100 serviceA
```

---

## 📚 Referências

### Documentação Oficial
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/)
- [OpenTelemetry Instrumentation](https://opentelemetry.io/docs/languages/go/instrumentation/)
- [OTEL Collector](https://opentelemetry.io/docs/collector/)
- [Zipkin Documentation](https://zipkin.io/)
- [Go Documentation](https://go.dev/doc/)

### APIs Utilizadas
- [ViaCEP API](https://viacep.com.br/)
- [WeatherAPI Docs](https://www.weatherapi.com/docs/)

### Arquitetura e Design
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

---

## 📝 Licença

Este projeto foi desenvolvido para fins educacionais como parte do laboratório de **Observabilidade e OpenTelemetry**.

---

## 👨‍💻 Autor
**Autor:** Adalberto F. Jr.  
**Repositório:** https://github.com/adalbertofjr/lab-2-observabilidade-e-opentelemetri
