# myScalidraw

A self-hosted drawing application inspired by Excalidraw, allowing users to create, edit, and save drawings with a complete file management system.

## 🎨 Features

- **Drawing Canvas**: Powered by Excalidraw for creating beautiful diagrams and drawings
- **File Management**: Create, rename, delete, and organize your drawings in a tree structure
- **Real-time Saving**: Automatic saving of your work to prevent data loss
- **Responsive Design**: Modern UI built with React and Tailwind CSS
- **Self-hosted**: Complete control over your data with your own instance
- **File Storage**: Secure file storage using MinIO object storage
- **Database Integration**: PostgreSQL for metadata and file organization

## 🏗️ Architecture

The project consists of:

- **Frontend**: React + TypeScript + Vite application with modern UI components
- **Backend**: Go (Golang) REST API using Fiber framework
- **Database**: PostgreSQL for storing file metadata and organization
- **Storage**: MinIO for secure file storage
- **Infrastructure**: Docker Compose for easy deployment

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Git

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd myScalidraw
   ```

2. **Configure environment variables**

   ```bash
   cp env.example .env
   ```

   Edit the `.env` file with your desired configuration:

   ```env
   # HTTP Configuration
   URL=http://localhost
   PORT=8181

   # Database Configuration
   URL_DB=postgres://myuser:mypassword@postgres:5432/myscalidraw?sslmode=disable

   # Frontend Configuration
   FRONTEND_URL=http://localhost:5173
   VITE_API_BASE_URL=http://localhost:8181

   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key-here
   URL_SHORTENED_PREFIX=http://localhost:8181

   # PostgreSQL
   POSTGRES_USER=myuser
   POSTGRES_PASSWORD=mypassword
   POSTGRES_DB=myscalidraw
   POSTGRES_PORT=5432

   # MinIO Configuration
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

3. **Start the application**

   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8181
   - MinIO Console: http://localhost:9001

## 🛠️ Development

### Local Development Setup

#### Backend Development

1. **Prerequisites**

   - Go 1.25.0+
   - PostgreSQL
   - MinIO

2. **Setup**
   ```bash
   cd back-end
   go mod download
   go run cmd/main.go
   ```

#### Frontend Development

1. **Prerequisites**

   - Node.js 18+
   - npm or yarn

2. **Setup**
   ```bash
   cd front-end
   npm install
   npm run dev
   ```

### Tech Stack

#### Frontend

- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **Excalidraw** - Drawing canvas
- **Radix UI** - Accessible UI components
- **TanStack Query** - Data fetching
- **Zustand** - State management
- **ky** - HTTP client

#### Backend

- **Go (Golang)** - Programming language
- **Fiber** - Web framework
- **GORM** - ORM for database operations
- **PostgreSQL** - Primary database
- **MinIO** - Object storage
- **Uber FX** - Dependency injection

## 📁 Project Structure

```
myScalidraw/
├── back-end/                 # Go backend application
│   ├── cmd/                  # Application entry points
│   ├── internal/             # Private application code
│   │   ├── delivery/         # HTTP handlers
│   │   ├── domain/           # Business logic and models
│   │   └── useCase/          # Application use cases
│   ├── infra/                # Infrastructure code
│   └── pkg/                  # Reusable packages
├── front-end/                # React frontend application
│   ├── src/
│   │   ├── components/       # Reusable UI components
│   │   ├── pages/           # Application pages
│   │   ├── hooks/           # Custom React hooks
│   │   ├── stores/          # State management
│   │   └── lib/             # Utility functions
│   └── public/              # Static assets
├── docker-compose.yml        # Docker composition
└── env.example              # Environment variables template
```

## 🔧 Configuration

The application uses environment variables for configuration. Key settings include:

- **Database**: PostgreSQL connection string
- **Storage**: MinIO endpoint and credentials
- **Security**: JWT secret for authentication
- **CORS**: Frontend URL for cross-origin requests

## 🚀 Deployment

### Using Docker (Recommended)

The easiest way to deploy myScalidraw is using the provided Docker Compose configuration:

```bash
docker-compose up -d
```

### Manual Deployment

1. Deploy PostgreSQL and MinIO instances
2. Build and deploy the Go backend
3. Build and deploy the React frontend
4. Configure reverse proxy (nginx/traefik) if needed

## 📝 API Documentation

The backend provides a REST API for file management:

- `POST /api/files` - Create a new file
- `GET /api/files` - List all files
- `GET /api/files/:id` - Get file by ID
- `PUT /api/files/:id` - Update file
- `DELETE /api/files/:id` - Delete file

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🙏 Acknowledgments

- [Excalidraw](https://excalidraw.com/) - For the amazing drawing canvas
- [Fiber](https://gofiber.io/) - For the fast Go web framework
- [React](https://reactjs.org/) - For the UI framework
- [MinIO](https://min.io/) - For object storage solution

---

**myScalidraw** - Self-hosted drawing application for creating and managing your diagrams and drawings.

## Languages

- [English](README.md)
- [Português (Brasil)](README.pt-BR.md)
