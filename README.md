# Go CLI App

A small Go HTTP API that exposes CRUD endpoints for an in-memory todo list.

## Prerequisites

- Go 1.22 or newer
- Git
- Docker, optional for container runs

## Run

```powershell
go run ./cmd/go-language-app
```

The API starts on `http://localhost:8080`.

To use another port:

```powershell
go run ./cmd/go-language-app -port 9090
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

Build the Docker image:

```powershell
docker build -t go-cli-app .
```

Run the container:

```powershell
docker run --rm -p 8080:8080 go-cli-app
```

Run on another host port:

```powershell
docker run --rm -p 9090:8080 go-cli-app
```

Then call the API:

```powershell
curl.exe http://localhost:8080/health
```
