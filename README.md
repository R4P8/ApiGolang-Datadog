# ApiGolangv2 🚀

A simple RESTful API written in Go, focused on task management (CRUD operations). The project is structured with clean architecture principles and includes Docker support for easy development and deployment.

## 🗂️ Project Structure

```
cmd/                - application entrypoint (main.go)
internal/           - core application code (config, entity, handler, repository, service, routes)
pkg/                - library code for external use (if any)
Dockerfile          - definition for building the service container
Docker-compose.yaml - compose file to run the app with dependencies (e.g. database, redis)
go.mod              - Go module definition
```

## ⚙️ Features

- Task entity with create, read, update, delete operations
- Configuration management (database, Redis)
- Clean separation of concerns (handlers, services, repositories)
- Dockerized for containerized development and deployment

## 🛠️ Prerequisites

- Go 1.18+ installed
- Docker & Docker Compose (optional, for containerized setup)

## 🚀 Getting Started

### Running locally (without Docker)

```bash
# clone repository
git clone <repo-url>
cd ApiGolangv2

# download dependencies
go mod tidy

# run application
go run ./cmd
```

The API will be available on `http://localhost:8080` by default (depending on configuration).

### Using Docker

```bash
docker-compose up --build
```

This will start the service along with any configured dependencies (e.g. database, redis).

## 📝 API Endpoints

Example endpoints (adjust according to `routes/routes.go`):

- `GET /tasks`             - list all tasks
- `GET /tasks/{id}`        - get a single task
- `POST /tasks`            - create a task
- `PUT /tasks/{id}`        - update a task
- `DELETE /tasks/{id}`     - delete a task


> **Note:** Update this section as the routes in your project evolve.


## 🧩 Configuration

Configuration values are managed in `internal/config`. Adjust environment variables or config files as needed.

Common fields:

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`
- `REDIS_ADDR`, `REDIS_PASS`
- `SERVER_PORT`

## 📦 Build/Run with Make (if applicable)

You can add a `Makefile` for convenience or run the build manually:

```bash
go build -o bin/api cmd/main.go
./bin/api
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Open a pull request

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

