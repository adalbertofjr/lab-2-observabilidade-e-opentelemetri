# Executando Testes com Docker

Este guia mostra como executar os testes da aplicação usando Docker e Docker Compose.

## Pré-requisitos

- Docker instalado
- Docker Compose instalado

## Comandos Disponíveis

### Usando Make (Recomendado)

```bash
# Ver todos os comandos disponíveis
make help

# Executar testes localmente (sem Docker)
make test-local

# Executar testes no Docker
make test-docker

# Gerar relatório de cobertura no Docker
make test-docker-coverage

# Abrir shell interativo no container
make test-docker-shell

# Limpar containers e imagens
make docker-clean
```

### Usando Docker Compose Diretamente

```bash
# Executar testes simples
docker-compose -f docker-compose.test.yml run --rm test

# Executar com relatório de cobertura
docker-compose -f docker-compose.test.yml run --rm test-coverage

# Shell interativo para testes
docker-compose -f docker-compose.test.yml run --rm test-watch
```

### Usando Docker Diretamente

```bash
# Build da imagem
docker build -f Dockerfile.test -t weather-api-test .

# Executar testes
docker run --rm weather-api-test

# Executar testes específicos
docker run --rm weather-api-test go test -v ./internal/usecase/weather/...
```

## Estrutura dos Arquivos

- **Dockerfile.test**: Imagem Docker otimizada para execução de testes
- **docker-compose.test.yml**: Orquestração com diferentes serviços de teste
- **.dockerignore**: Arquivos excluídos do build Docker
- **Makefile**: Comandos simplificados para execução

## Serviços Disponíveis

### `test`
Executa todos os testes com output verbose e cobertura básica.

```bash
docker-compose -f docker-compose.test.yml run --rm test
```

### `test-coverage`
Gera relatório HTML de cobertura na pasta `coverage/`.

```bash
docker-compose -f docker-compose.test.yml run --rm test-coverage
# Abre coverage/coverage.html no navegador
```

### `test-watch`
Shell interativo para executar testes customizados.

```bash
docker-compose -f docker-compose.test.yml run --rm test-watch
# Dentro do container:
go test -v ./internal/domain/entity/...
go test -run TestGetWeather_Success ./...
```

## Exemplos de Uso

### Executar todos os testes
```bash
make test-docker
```

### Gerar relatório de cobertura
```bash
make test-docker-coverage
open coverage/coverage.html  # macOS
```

### Executar teste específico
```bash
docker-compose -f docker-compose.test.yml run --rm test \
  go test -v -run TestGetCurrentWeather_Success ./internal/usecase/weather/
```

### Executar testes de um pacote específico
```bash
docker-compose -f docker-compose.test.yml run --rm test \
  go test -v ./internal/domain/entity/
```

### Verificar cobertura de pacote específico
```bash
docker-compose -f docker-compose.test.yml run --rm test \
  go test -cover ./internal/usecase/weather/
```

## Volumes e Persistência

O `docker-compose.test.yml` monta o código fonte como volume, permitindo:
- Testar alterações sem rebuild da imagem
- Persistir relatórios de cobertura no host

Para executar sem montar volumes (imagem isolada):
```bash
docker build -f Dockerfile.test -t weather-api-test . && docker run --rm weather-api-test
```

## Otimizações

### Cache de Dependências
O `Dockerfile.test` usa multi-stage com `go.mod` e `go.sum` copiados primeiro, otimizando o cache do Docker quando apenas o código muda (não as dependências).

### Build Rápido
```bash
# Build apenas se necessário
make docker-build

# Executar testes sem rebuild
make test-docker
```

## Limpeza

```bash
# Remover containers, imagens e volumes
make docker-clean

# Remover apenas arquivos de cobertura local
make clean
```

## CI/CD

Exemplo de uso em pipeline CI/CD:

```yaml
# .github/workflows/test.yml
test:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v3
    - run: docker-compose -f docker-compose.test.yml run --rm test
```

## Troubleshooting

### Erro de permissão em `coverage/`
```bash
# Linux: ajustar permissões
sudo chown -R $USER:$USER coverage/
```

### Container não encontra módulos
```bash
# Forçar rebuild sem cache
docker-compose -f docker-compose.test.yml build --no-cache test
```

### Testes lentos
```bash
# Executar apenas testes modificados
docker-compose -f docker-compose.test.yml run --rm test \
  go test -v ./internal/usecase/weather/ ./internal/infra/api/
```
