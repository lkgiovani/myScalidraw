# MyScalidraw - Docker Deployment Guide

Este guia explica como fazer build e deploy das imagens Docker do MyScalidraw.

## 📋 Pré-requisitos

- Docker e Docker Compose instalados
- Conta no Docker Hub (ou outro registry)
- Make (opcional, mas recomendado)

## 🚀 Quick Start

### 1. Configurar suas credenciais

Edite o arquivo `env.example` e renomeie para `.env`:

```bash
cp env.example .env
# Edite o arquivo .env com suas configurações
```

### 2. Fazer login no Docker Hub

```bash
make login
# ou
docker login docker.io
```

### 3. Build e Push das imagens

```bash
# Usando Make (recomendado)
make push USERNAME=seu-usuario VERSION=v1.0.0

# Ou usando o script bash
chmod +x build-and-push.sh
./build-and-push.sh v1.0.0
```

### 4. Deploy

```bash
# Desenvolvimento (build local)
make deploy-dev

# Produção (imagens do registry)
make deploy-prod USERNAME=seu-usuario VERSION=v1.0.0
```

## 📦 Comandos Disponíveis

### Make Commands

```bash
make help          # Mostra todos os comandos disponíveis
make build          # Build das imagens localmente
make push           # Build e push para o registry
make deploy-dev     # Deploy ambiente de desenvolvimento
make deploy-prod    # Deploy ambiente de produção
make stop           # Para todos os containers
make clean          # Limpa imagens e containers
make logs           # Mostra logs dos containers
make status         # Status dos containers
```

### Comandos Manuais

```bash
# Build das imagens
docker build -t seu-usuario/myscalidraw-backend:latest ./back-end
docker build -t seu-usuario/myscalidraw-frontend:latest ./front-end

# Push para o registry
docker push seu-usuario/myscalidraw-backend:latest
docker push seu-usuario/myscalidraw-frontend:latest

# Deploy com docker-compose
docker-compose up -d                           # Desenvolvimento
docker-compose -f docker-compose.prod.yml up -d  # Produção
```

## 🌐 Acessos

Após o deploy, os serviços estarão disponíveis em:

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **MinIO Console**: http://localhost:9001 (admin/admin123)
- **PostgreSQL**: localhost:5432

## 🔧 Configurações

### Variáveis de Ambiente

| Variável   | Descrição                | Padrão        |
| ---------- | ------------------------ | ------------- |
| `REGISTRY` | Registry Docker          | `docker.io`   |
| `USERNAME` | Seu username no registry | `seu-usuario` |
| `VERSION`  | Versão da imagem         | `latest`      |

### Portas

| Serviço       | Porta Host | Porta Container |
| ------------- | ---------- | --------------- |
| Frontend      | 3000       | 80              |
| Backend       | 8080       | 8080            |
| PostgreSQL    | 5432       | 5432            |
| MinIO API     | 9000       | 9000            |
| MinIO Console | 9001       | 9001            |

## 📁 Estrutura dos Arquivos

```
├── back-end/
│   └── Dockerfile              # Dockerfile do backend Go
├── front-end/
│   ├── Dockerfile              # Dockerfile do frontend React
│   └── nginx.conf              # Configuração do Nginx
├── docker-compose.yml          # Desenvolvimento (build local)
├── docker-compose.prod.yml     # Produção (imagens do registry)
├── build-and-push.sh           # Script para build e push
├── Makefile                    # Comandos Make
└── env.example                 # Exemplo de variáveis de ambiente
```

## 🔍 Troubleshooting

### Problemas Comuns

1. **Erro de permissão no Docker**:

   ```bash
   sudo usermod -aG docker $USER
   # Faça logout e login novamente
   ```

2. **Erro de autenticação no registry**:

   ```bash
   docker login docker.io
   ```

3. **Containers não sobem**:
   ```bash
   make logs  # Verificar logs
   make clean # Limpar e tentar novamente
   ```

### Logs e Debug

```bash
# Ver logs de todos os serviços
make logs

# Ver logs de um serviço específico
docker-compose logs -f backend
docker-compose logs -f frontend

# Status dos containers
make status
```

## 🚀 Deploy em Produção

Para deploy em produção, recomenda-se:

1. Usar um registry privado
2. Configurar secrets para credenciais
3. Usar volumes persistentes para dados
4. Configurar backup do PostgreSQL
5. Usar HTTPS com certificados SSL
6. Configurar monitoring e logs centralizados

### Exemplo com Docker Swarm

```bash
# Inicializar swarm
docker swarm init

# Deploy da stack
docker stack deploy -c docker-compose.prod.yml myscalidraw
```
