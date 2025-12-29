# Find My Phone API

Production-ready REST API built with Gin and Clean Architecture for tracking devices and their last known locations.

## Features
- Users linked to devices via a public `device_id`
- Device registry with IMEI uniqueness and lost flag
- Location logs with latest lookup by device
- Clean Architecture layers and structured logging
- PostgreSQL (default) with AutoMigrate, SQLite fallback
- CORS, request ID, graceful shutdown

## Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose (for Postgres)

### Environment Variables
| Name | Default | Description |
| --- | --- | --- |
| `SERVER_ADDRESS` | `:8081` | Address the HTTP server listens on |
| `DATABASE_TYPE` | `postgres` | `postgres` or `sqlite` |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/findmyphone?sslmode=disable` | Connection string |
| `SHUTDOWN_GRACE` | `15s` | Graceful shutdown timeout |

### Run with Docker Compose (Postgres)
```bash
docker-compose up -d
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/findmyphone?sslmode=disable"
go run ./cmd/server
```

**Troubleshooting:** If Docker cleaned up old images and `docker-compose up -d` fails with a
`No such image` message, rebuild the app image before starting containers:
```bash
docker-compose down --remove-orphans
docker-compose build
docker-compose up -d
```

### Run with SQLite (quick start)
```bash
export DATABASE_TYPE=sqlite
export DATABASE_URL=./findmyphone.sqlite
go run ./cmd/server
```

### Testing
```bash
go test ./...
```

## API
Base path: `/api/v1`

### Create User
`POST /api/v1/users`
```bash
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","device_id":"device-123","phone":"+12025550123","backup_phone":"+12025550124"}'
```

### Get User by device_id
`GET /api/v1/users/by-device/:device_id`
```bash
curl http://localhost:8081/api/v1/users/by-device/device-123
```

### Update User by device_id
`POST /api/v1/users/by-device/:device_id`
```bash
curl -X POST http://localhost:8081/api/v1/users/by-device/device-123 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","phone":"+12025550123","backup_phone":"+12025550124"}'
```

### Create Device
`POST /api/v1/devices`
```bash
curl -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{"device_id":"device-123","imei":"123456789012345","generation":"gen1","name":"Pixel","lost":false}'
```

### Create Log
`POST /api/v1/logs`
```bash
curl -X POST http://localhost:8081/api/v1/logs \
  -H "Content-Type: application/json" \
  -d '{"device_id":"device-123","longitude":-122.1,"latitude":37.4}'
```

### Get Last Log by device_id
`GET /api/v1/logs/last/:device_id`
```bash
curl http://localhost:8081/api/v1/logs/last/device-123
```

## Project Structure
```
cmd/server        # application entrypoint
internal/domain   # entities and domain contracts
internal/usecase  # business logic
internal/infrastructure # db, repositories, config, logger
internal/interface/http # transport layer (Gin)
```
