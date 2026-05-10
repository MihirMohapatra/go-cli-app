# Go CLI App

A small Go HTTP API that exposes CRUD endpoints for a PostgreSQL-backed todo list.

## Prerequisites

- Go 1.22 or newer
- Git
- PostgreSQL, for local runs without Docker
- Docker, optional for container runs

## Configuration

The app reads the PostgreSQL connection string from `DATABASE_URL`.

Default value:

```text
postgres://postgres:postgres@localhost:5432/go_cli_app?sslmode=disable
```

The app creates the `todos` table automatically on startup.

## Run

```powershell
go run ./cmd/go-language-app
```

The API starts on `http://localhost:8080`.

Make sure PostgreSQL is running first and that the `go_cli_app` database exists.

To use another port:

```powershell
go run ./cmd/go-language-app -port 9090
```

To use a custom database URL:

```powershell
$env:DATABASE_URL = "postgres://postgres:postgres@localhost:5432/go_cli_app?sslmode=disable"
go run ./cmd/go-language-app
```

## Endpoints

| Method | Path | Description |
| --- | --- | --- |
| GET | `/health` | Check API health |
| GET | `/todos` | List all todos |
| POST | `/todos` | Create a todo |
| GET | `/todos/{id}` | Get one todo |
| PUT | `/todos/{id}` | Update one todo |
| DELETE | `/todos/{id}` | Delete one todo |

## Examples

Create a todo:

```powershell
curl.exe -X POST http://localhost:8080/todos `
  -H "Content-Type: application/json" `
  -d "{\"title\":\"Learn Go\"}"
```

List todos:

```powershell
curl.exe http://localhost:8080/todos
```

Get one todo:

```powershell
curl.exe http://localhost:8080/todos/1
```

Update a todo:

```powershell
curl.exe -X PUT http://localhost:8080/todos/1 `
  -H "Content-Type: application/json" `
  -d "{\"title\":\"Learn Go\",\"completed\":true}"
```

Delete a todo:

```powershell
curl.exe -X DELETE http://localhost:8080/todos/1
```

## Test

```powershell
go test ./...
```

## Build

```powershell
go build -o bin/go-language-app.exe ./cmd/go-language-app
```

## Docker

Run the API and PostgreSQL together:

```powershell
docker compose up --build
```

Stop the stack:

```powershell
docker compose down
```

Stop the stack and delete the database volume:

```powershell
docker compose down -v
```

Build the Docker image:

```powershell
docker build -t go-cli-app .
```

Run only the API container against a PostgreSQL database running on your machine:

```powershell
docker run --rm -p 8080:8080 `
  -e DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/go_cli_app?sslmode=disable" `
  go-cli-app
```

Run on another host port:

```powershell
docker run --rm -p 9090:8080 `
  -e DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/go_cli_app?sslmode=disable" `
  go-cli-app
```

Then call the API:

```powershell
curl.exe http://localhost:8080/health
```
