# coco-task-manager 
*A lightweight cron expression parser and task scheduler written in Go*

`coco-task-manager` is a simple, fast, and self-contained **cron expression parser** and **task scheduler** that expands a standard cron string into its individual components and manages scheduled tasks via CLI or HTTP API.  

## Features

- Parses standard cron expressions with 5 fields:
  - minute
  - hour
  - day of month
  - month
  - day of week
- Supports:
  - Wildcards (`*`)
  - Lists (e.g. `1,15,30`)
  - Ranges (e.g. `1-5`)
  - Step values (e.g. `*/15`)
  - Singular values (e.g. `4`)
- **Task scheduling** via CLI and HTTP API
- **Message queue integration** (RabbitMQ) for task execution

## Dependencies
- Built with Go 1.25+
- **Docker & Docker Compose** - For containerized deployment
- **Task** (taskfile.dev) - For running project tasks
- **golangci-lint** - For linting and formatting

## Installation
### Clone the repository
```bash
git clone <repository-url>
cd coco-cron-parser
```

### Install Go dependencies
```bash
go mod download
```

## Configuration

The application can be configured via environment variables or a `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `CRONTAB_FILE` | Path to the crontab file | `./e2e/storage/crontab` |
| `RABBITMQ_HOST` | RabbitMQ connection URL | `amqp://localhost:5672` |
| `RABBITMQ_USER` | RabbitMQ username | `guest` |
| `RABBITMQ_PASS` | RabbitMQ password | `guest` |

Example `.env` file:
```env
CRONTAB_FILE="/path/to/crontab"
RABBITMQ_USER="guest"
RABBITMQ_PASS="guest"
RABBITMQ_HOST="amqp://localhost:5672"
```

## Usage

### Running the CLI

The CLI provides commands for scheduling tasks and managing the message queue:

```bash
# Schedule a task
go run ./cmd/cli schedule-task "*/15 * * * *" "start-game 123"

# Start a game (sends message to dealer API)
go run ./cmd/cli start-game <room_id>

# Pull messages from a topic (for debugging)
go run ./cmd/cli pull-messages <topic>
```

### Running the API Server

The HTTP API provides endpoints for task management:

```bash
# Start the server
go run ./cmd/api
```

The server will start on port `:3000`.

#### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/livez` | Health check endpoint |
| GET | `/api/v1/tasks/` | List all tasks |
| GET | `/api/v1/tasks/scheduled` | List scheduled tasks |
| POST | `/api/v1/tasks/` | Schedule a new task |
| DELETE | `/api/v1/tasks/{uuid}` | Remove a task |

### Running with Docker

#### Build the Docker image
```bash
docker build -t coco-task-manager:latest -f ./build/docker/Dockerfile .
```

Or use Task:
```bash
task build-image
```

#### Run with Docker Compose (includes RabbitMQ)
```bash
# Start RabbitMQ
docker-compose up -d

# Run the application
docker run -it --rm coco-task-manager:latest
```

## Development

### Project Tasks (Taskfile)

The project uses [Task](https://taskfile.dev/) for common development tasks. List them with `task --list`

### Air Configuration

The project includes `.air.toml` for live reloading during development:

```bash
# Install air
go install github.com/air-serve/air@latest

# Run with live reload
air
```

