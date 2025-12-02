.PHONY: help test test-local test-docker test-coverage test-docker-coverage docker-build docker-clean clean

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

test-local: ## Executa testes localmente
	go test -v -cover ./...

test: test-local ## Alias para test-local

test-coverage: ## Gera relatório de cobertura local
	@mkdir -p coverage
	go test -v -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Relatório de cobertura gerado em: coverage/coverage.html"

docker-build: ## Constrói a imagem Docker para testes
	docker compose -f docker-compose.test.yml build test

test-docker: ## Executa testes no container Docker
	docker compose -f docker-compose.test.yml run --rm test

test-docker-coverage: ## Executa testes com cobertura no container Docker
	@mkdir -p coverage
	docker compose -f docker-compose.test.yml run --rm test-coverage
	@echo "Relatório de cobertura gerado em: coverage/coverage.html"

test-docker-shell: ## Abre shell no container para testes interativos
	docker compose -f docker-compose.test.yml run --rm test-watch

docker-clean: ## Remove containers e imagens Docker de teste
	docker compose -f docker-compose.test.yml down --rmi all --volumes
	@echo "Containers e imagens removidos"

clean: ## Remove arquivos de cobertura
	rm -rf coverage/
	find . -name "*.out" -type f -delete
	find . -name "*.test" -type f -delete
