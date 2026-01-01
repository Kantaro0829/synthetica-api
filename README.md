# Synthetica API

Welcome to the **Synthetica API**! This project is a backend REST API service built with **Go (Golang)**, following the **Clean Architecture** principles. It uses **PostgreSQL** as the database and runs in a **Docker** container.

This guide is designed to help you get started, understand the structure, and run the project, even if you are new to Go or backend development!

## 🚀 Tech Stack

- **Language**: [Go](https://go.dev/) (1.20+)
- **Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin) (HTTP Router)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **ORM**: [GORM](https://gorm.io/) (Object Relational Mapper)
- **Logger**: [Zap](https://github.com/uber-go/zap) (Blazing fast, structured logging)
- **Containerization**: [Docker](https://www.docker.com/) & MySQL

## 🏗️ Architecture (Clean Architecture)

This project follows the **Clean Architecture** design pattern to separate concerns and make the code testable and maintainable.

```text
cmd/api/        --> Main entry point of the application.
internal/       --> Private application code.
  ├── domain/   --> Core business models (Entities) & Interface definitions. (No external dependencies here!)
  ├── repository/  --> Database logic (Implementation of domain interfaces).
  ├── usecase/  --> Business logic (Connects Delivery to Repository).
  └── delivery/ --> HTTP Handlers (Gin) to handle requests.
pkg/            --> Public library code (Logger, Database connection helper).
```

### Flow of Request
1. **HTTP Request** hits the **Delivery Layer** (Handler).
2. Handler calls the **Usecase Layer**.
3. Usecase calls the **Repository Layer** to get/save data.
4. Repository interacts with the **Database**.
5. Data flows back up to the user.

## 🛠️ Prerequisites

Before you start, make sure you have installed:
1. **Go**: [Install Go](https://go.dev/doc/install)
2. **Docker Desktop**: [Install Docker](https://docs.docker.com/get-docker/)

## 🏃‍♂️ How to Run

### 1. Clone the Project
```bash
git clone <repository_url>
cd synthetica
```

### 2. Start the Database
We use Docker to run PostgreSQL. This saves you from installing Postgres manually on your machine.
```bash
docker compose up -d
```
*This command downloads the Postgres image and starts it in the background.*

### 3. Run the API Server
Now, start the Go application.
```bash
go run cmd/api/main.go
```
You should see logs indicating the server is running on `:8080`.

## 🔌 API Endpoints

You can test these using `curl` or a tool like [Postman](https://www.postman.com/).

### 1. Create a User
**POST** `/users`

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Kantaro", "email":"kan@example.com", "password":"password123"}'
```

**Response:**
```json
{
  "id": 1,
  "name": "Kantaro",
  "email": "kan@example.com",
  "created_at": "...",
  "updated_at": "..."
}
```

### 2. Get All Users
**GET** `/users`

```bash
curl http://localhost:8080/users
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Kantaro",
    "email": "kan@example.com",
    ...
  }
]
```

### 3. Get User by ID
**GET** `/users/:id`

```bash
curl http://localhost:8080/users/1
```

## 🧪 How to Test

We use `go test` and Docker to run integration tests for the repository layer.

### 1. Start Test Database
Navigate to the `test` directory and start the test database:
```bash
cd test
docker-compose up -d
```

### 2. Run Tests
Run all repository tests:
```bash
go test -v ./internal/repository/...
```

## ❓ Troubleshooting

- **"Connection Refused"**: Make sure your Docker container is running (`docker compose ps`).
- **"bind: address already in use"**: Something else is running on port 8080. Kill it or change the port in `main.go`.

Happy Coding! 🚀
