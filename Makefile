# Configurações
REGISTRY ?= docker.io
USERNAME ?= seu-usuario
PROJECT_NAME = myscalidraw
VERSION ?= latest

# Cores para output
GREEN = \033[0;32m
YELLOW = \033[1;33m
RED = \033[0;31m
NC = \033[0m # No Color

.PHONY: help build push deploy clean login

help: ## Mostra esta ajuda
	@echo "$(GREEN)MyScalidraw Docker Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

login: ## Faz login no Docker registry
	@echo "$(GREEN)🔐 Fazendo login no Docker registry...$(NC)"
	docker login $(REGISTRY)

build: ## Faz build das imagens localmente
	@echo "$(GREEN)📦 Building images...$(NC)"
	docker build -t $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-backend:$(VERSION) ./back-end
	docker build -t $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-backend:latest ./back-end
	docker build -t $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-frontend:$(VERSION) ./front-end
	docker build -t $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-frontend:latest ./front-end
	@echo "$(GREEN)✅ Build completed!$(NC)"

push: build ## Faz build e push das imagens para o registry
	@echo "$(GREEN)📤 Pushing images...$(NC)"
	docker push $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-backend:$(VERSION)
	docker push $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-backend:latest
	docker push $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-frontend:$(VERSION)
	docker push $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-frontend:latest
	@echo "$(GREEN)✅ Push completed!$(NC)"

deploy-dev: ## Sobe o ambiente de desenvolvimento
	@echo "$(GREEN)🚀 Starting development environment...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)✅ Development environment started!$(NC)"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend: http://localhost:8080"
	@echo "MinIO Console: http://localhost:9001"

deploy-prod: ## Sobe o ambiente de produção com imagens do registry
	@echo "$(GREEN)🚀 Starting production environment...$(NC)"
	REGISTRY=$(REGISTRY) USERNAME=$(USERNAME) VERSION=$(VERSION) docker-compose -f docker-compose.prod.yml up -d
	@echo "$(GREEN)✅ Production environment started!$(NC)"

stop: ## Para todos os containers
	@echo "$(YELLOW)⏹️  Stopping containers...$(NC)"
	docker-compose down
	docker-compose -f docker-compose.prod.yml down
	@echo "$(GREEN)✅ Containers stopped!$(NC)"

clean: stop ## Para containers e remove imagens locais
	@echo "$(YELLOW)🧹 Cleaning up...$(NC)"
	docker image rm $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-backend:$(VERSION) 2>/dev/null || true
	docker image rm $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-backend:latest 2>/dev/null || true
	docker image rm $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-frontend:$(VERSION) 2>/dev/null || true
	docker image rm $(REGISTRY)/$(USERNAME)/$(PROJECT_NAME)-frontend:latest 2>/dev/null || true
	docker system prune -f
	@echo "$(GREEN)✅ Cleanup completed!$(NC)"

logs: ## Mostra logs dos containers
	docker-compose logs -f

status: ## Mostra status dos containers
	docker-compose ps
