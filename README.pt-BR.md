# myScalidraw

Uma aplicação de desenho auto-hospedada inspirada no Excalidraw, permitindo que usuários criem, editem e salvem desenhos com um sistema completo de gerenciamento de arquivos.

## 🎨 Funcionalidades

- **Canvas de Desenho**: Powered by Excalidraw para criar diagramas e desenhos bonitos
- **Gerenciamento de Arquivos**: Crie, renomeie, exclua e organize seus desenhos em uma estrutura de árvore
- **Salvamento em Tempo Real**: Salvamento automático do seu trabalho para prevenir perda de dados
- **Design Responsivo**: UI moderna construída com React e Tailwind CSS
- **Auto-hospedado**: Controle completo sobre seus dados com sua própria instância
- **Armazenamento de Arquivos**: Armazenamento seguro de arquivos usando MinIO object storage
- **Integração com Banco de Dados**: PostgreSQL para metadados e organização de arquivos

## 🏗️ Arquitetura

O projeto consiste em:

- **Frontend**: Aplicação React + TypeScript + Vite com componentes UI modernos
- **Backend**: API REST em Go (Golang) usando framework Fiber
- **Banco de Dados**: PostgreSQL para armazenar metadados e organização de arquivos
- **Armazenamento**: MinIO para armazenamento seguro de arquivos
- **Infraestrutura**: Docker Compose para deploy fácil

## 🚀 Início Rápido

### Pré-requisitos

- Docker e Docker Compose
- Git

### Instalação

1. **Clone o repositório**

   ```bash
   git clone <repository-url>
   cd myScalidraw
   ```

2. **Configure as variáveis de ambiente**

   ```bash
   cp env.example .env
   ```

   Edite o arquivo `.env` com sua configuração desejada:

   ```env
   # Configuração HTTP
   URL=http://localhost
   PORT=8181

   # Configuração do Banco de Dados
   URL_DB=postgres://myuser:mypassword@postgres:5432/myscalidraw?sslmode=disable

   # Configuração do Frontend
   FRONTEND_URL=http://localhost:5173
   VITE_API_BASE_URL=http://localhost:8181

   # Configuração JWT
   JWT_SECRET=sua-chave-jwt-super-secreta-aqui
   URL_SHORTENED_PREFIX=http://localhost:8181

   # PostgreSQL
   POSTGRES_USER=myuser
   POSTGRES_PASSWORD=mypassword
   POSTGRES_DB=myscalidraw
   POSTGRES_PORT=5432

   # Configuração MinIO
   MINIO_ENDPOINT=minio:9000
   MINIO_ACCESS_KEY=minioadmin
   MINIO_SECRET_KEY=minioadmin123
   MINIO_BUCKET=myscalidraw
   MINIO_USE_SSL=false
   MINIO_ROOT_USER=minioadmin
   MINIO_ROOT_PASSWORD=minioadmin123
   MINIO_PORT=9000
   MINIO_CONSOLE_PORT=9001
   ```

3. **Inicie a aplicação**

   ```bash
   docker-compose up -d
   ```

4. **Acesse a aplicação**
   - Frontend: http://localhost:5173
   - API Backend: http://localhost:8181
   - Console MinIO: http://localhost:9001

## 🛠️ Desenvolvimento

### Configuração para Desenvolvimento Local

#### Desenvolvimento Backend

1. **Pré-requisitos**

   - Go 1.25.0+
   - PostgreSQL
   - MinIO

2. **Configuração**
   ```bash
   cd back-end
   go mod download
   go run cmd/main.go
   ```

#### Desenvolvimento Frontend

1. **Pré-requisitos**

   - Node.js 18+
   - npm ou yarn

2. **Configuração**
   ```bash
   cd front-end
   npm install
   npm run dev
   ```

### Stack Tecnológica

#### Frontend

- **React 18** - Framework UI
- **TypeScript** - Segurança de tipos
- **Vite** - Ferramenta de build
- **Tailwind CSS** - Estilização
- **Excalidraw** - Canvas de desenho
- **Radix UI** - Componentes UI acessíveis
- **TanStack Query** - Busca de dados
- **Zustand** - Gerenciamento de estado
- **ky** - Cliente HTTP

#### Backend

- **Go (Golang)** - Linguagem de programação
- **Fiber** - Framework web
- **GORM** - ORM para operações de banco de dados
- **PostgreSQL** - Banco de dados principal
- **MinIO** - Armazenamento de objetos
- **Uber FX** - Injeção de dependência

## 📁 Estrutura do Projeto

```
myScalidraw/
├── back-end/                 # Aplicação backend em Go
│   ├── cmd/                  # Pontos de entrada da aplicação
│   ├── internal/             # Código privado da aplicação
│   │   ├── delivery/         # Handlers HTTP
│   │   ├── domain/           # Lógica de negócio e modelos
│   │   └── useCase/          # Casos de uso da aplicação
│   ├── infra/                # Código de infraestrutura
│   └── pkg/                  # Pacotes reutilizáveis
├── front-end/                # Aplicação frontend em React
│   ├── src/
│   │   ├── components/       # Componentes UI reutilizáveis
│   │   ├── pages/           # Páginas da aplicação
│   │   ├── hooks/           # Hooks personalizados do React
│   │   ├── stores/          # Gerenciamento de estado
│   │   └── lib/             # Funções utilitárias
│   └── public/              # Assets estáticos
├── docker-compose.yml        # Composição Docker
└── env.example              # Template de variáveis de ambiente
```

## 🔧 Configuração

A aplicação usa variáveis de ambiente para configuração. Configurações principais incluem:

- **Banco de Dados**: String de conexão PostgreSQL
- **Armazenamento**: Endpoint e credenciais do MinIO
- **Segurança**: Segredo JWT para autenticação
- **CORS**: URL do frontend para requisições cross-origin

## 🚀 Deploy

### Usando Docker (Recomendado)

A maneira mais fácil de fazer deploy do myScalidraw é usando a configuração Docker Compose fornecida:

```bash
docker-compose up -d
```

### Deploy Manual

1. Faça deploy das instâncias PostgreSQL e MinIO
2. Construa e faça deploy do backend Go
3. Construa e faça deploy do frontend React
4. Configure proxy reverso (nginx/traefik) se necessário

## 📝 Documentação da API

O backend fornece uma API REST para gerenciamento de arquivos:

- `POST /api/files` - Criar um novo arquivo
- `GET /api/files` - Listar todos os arquivos
- `GET /api/files/:id` - Buscar arquivo por ID
- `PUT /api/files/:id` - Atualizar arquivo
- `DELETE /api/files/:id` - Excluir arquivo

## 🤝 Contribuindo

1. Faça um fork do repositório
2. Crie sua branch de feature (`git checkout -b feature/funcionalidade-incrivel`)
3. Faça commit das suas mudanças (`git commit -m 'feat: adiciona funcionalidade incrível'`)
4. Faça push para a branch (`git push origin feature/funcionalidade-incrivel`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo LICENSE para detalhes.

```
Licença MIT

Copyright (c) 2024 myScalidraw

É concedida permissão, gratuitamente, a qualquer pessoa que obtenha uma cópia
deste software e arquivos de documentação associados (o "Software"), para lidar
no Software sem restrição, incluindo, sem limitação, os direitos
de usar, copiar, modificar, mesclar, publicar, distribuir, sublicenciar e/ou vender
cópias do Software, e permitir às pessoas a quem o Software é
fornecido a fazê-lo, sujeito às seguintes condições:

O aviso de copyright acima e este aviso de permissão devem ser incluídos em todas
as cópias ou partes substanciais do Software.

O SOFTWARE É FORNECIDO "COMO ESTÁ", SEM GARANTIA DE QUALQUER TIPO, EXPRESSA OU
IMPLÍCITA, INCLUINDO MAS NÃO SE LIMITANDO ÀS GARANTIAS DE COMERCIALIZAÇÃO,
ADEQUAÇÃO A UM PROPÓSITO ESPECÍFICO E NÃO VIOLAÇÃO. EM NENHUM CASO OS
AUTORES OU DETENTORES DE DIREITOS AUTORAIS SERÃO RESPONSÁVEIS POR QUALQUER REIVINDICAÇÃO, DANOS OU OUTRA
RESPONSABILIDADE, SEJA EM AÇÃO DE CONTRATO, DELITO OU DE OUTRA FORMA, DECORRENTE DE,
FORA DE OU EM CONEXÃO COM O SOFTWARE OU O USO OU OUTRAS NEGOCIAÇÕES NO
SOFTWARE.
```

## 🙏 Agradecimentos

- [Excalidraw](https://excalidraw.com/) - Pelo incrível canvas de desenho
- [Fiber](https://gofiber.io/) - Pelo framework web Go rápido
- [React](https://reactjs.org/) - Pelo framework UI
- [MinIO](https://min.io/) - Pela solução de armazenamento de objetos

---

**myScalidraw** - Aplicação de desenho auto-hospedada para criar e gerenciar seus diagramas e desenhos.

## Idiomas

- [English](README.md)
- [Português (Brasil)](README.pt-BR.md)
