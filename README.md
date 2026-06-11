# task-manager-go

A RESTful API for managing personal tasks, built as a portfolio project to demonstrate backend development with Go, PostgreSQL, Docker, and CI/CD using GitHub Actions.

## Features

- User authentication with JWT (register & login)
- Create tasks with title and description
- View all tasks belonging to the authenticated user
- View task detail by ID
- Update task status (`todo` → `in_progress` → `done`)

## Tech Stack

- **Language:** Go 1.26
- **Router:** [chi](https://github.com/go-chi/chi)
- **Database:** PostgreSQL 16
- **Auth:** JWT (golang-jwt)
- **Containerization:** Docker & Docker Compose
- **CI/CD:** GitHub Actions → Docker Hub
- **Deployment:** Render (API) + Supabase (Database)

## Project Structure

```
task-manager/
├── cmd/
│   └── api/
│       └── main.go           # Entry point, router setup
├── internal/
│   ├── handler/
│   │   ├── auth.go           # Register & login handlers
│   │   ├── task.go           # Task handlers
│   │   └── helper.go         # JSON response helpers
│   ├── middleware/
│   │   └── auth.go           # JWT middleware
│   ├── model/
│   │   ├── user.go
│   │   └── task.go
│   └── repository/
│       ├── user.go           # User database queries
│       ├── task.go           # Task database queries
│       └── errors.go
├── db/
│   └── migrations/
│       └── 001_init.sql
├── .github/
│   └── workflows/
│       └── ci.yml            # GitHub Actions workflow
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## API Endpoints

### Auth

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/auth/register` | Register a new account | ❌ |
| `POST` | `/auth/login` | Login and receive JWT token | ❌ |

### Tasks

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/tasks` | Get all tasks for the logged-in user | ✅ |
| `POST` | `/tasks` | Create a new task | ✅ |
| `GET` | `/tasks/{id}` | Get task detail by ID | ✅ |
| `PATCH` | `/tasks/{id}/status` | Update task status | ✅ |

### Task Status Flow

```
todo  ──→  in_progress  ──→  done
```

## Getting Started (Local)

### Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)

### 1. Clone the repository

```bash
git clone https://github.com/arsyad-id-99/task-manager-go.git
cd task-manager-go
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Edit `.env` and fill in the values:

```env
DATABASE_URL=postgres://appuser:apppassword@localhost:5432/taskdb?sslmode=disable
JWT_SECRET=your-random-secret-here
PORT=8080
```

Generate a secure JWT secret:

```bash
openssl rand -hex 32
```

### 3. Run with Docker Compose

```bash
docker compose up --build
```

The API will be available at `http://localhost:8080`.

### 4. Verify it's running

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## Usage Examples

### Register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Budi","email":"budi@mail.com","password":"rahasia123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"budi@mail.com","password":"rahasia123"}'
```

Save the token from the response, then use it in subsequent requests:

```bash
TOKEN=<your-jwt-token-here>
```

### Create a Task

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Belajar Go","description":"Selesaikan tutorial backend"}'
```

### List Tasks

```bash
curl http://localhost:8080/tasks \
  -H "Authorization: Bearer $TOKEN"
```

### Update Task Status

```bash
curl -X PATCH http://localhost:8080/tasks/{id}/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'
```

## CI/CD Pipeline

Every push to `main` triggers the GitHub Actions workflow:

1. **Lint & Test** — runs `go vet` and `go test`
2. **Build & Push** — builds Docker image and pushes to Docker Hub with `latest` and commit SHA tags

## Live Demo

Base URL: `https://task-manager-go-9xct.onrender.com`

> Note: The service may take ~30 seconds to wake up on the first request as it runs on Render's free tier.

## License

MIT