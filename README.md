# Webhook Inbox

Webhook Inbox is a Go + PostgreSQL application for capturing incoming webhooks, storing events, and reviewing them through a small Svelte web UI.

## What It Includes

- Go API under `/api/v1`
- PostgreSQL storage
- Svelte web UI for browsing sources, events, and delivery attempts
- Docker Compose setup for local development

## Project Structure

```text
webhook-inbox/
├── cmd/webhook-inbox/       # application entrypoint
├── internal/
│   ├── api/                 # HTTP routes, DTO mapping, request parsing
│   ├── config/              # environment-driven config
│   ├── db/                  # generated database layer and SQL
│   ├── delivery_engine/     # delivery and retry workers
│   ├── service/             # business logic
│   ├── struct_parser/       # request/query parsing helpers
│   └── utils/               # shared helpers
├── migrations/              # database migrations
├── web/                     # Svelte frontend
├── dev-tools/               # small local helper services
├── docker-compose.yml       # local containers
└── Dockerfile               # API image
```

## How To Run

### Docker Compose

```bash
docker compose --profile dev up --build
```

Services:

- API: `http://localhost:3001`
- Web: `http://localhost:5173`
- Dev test server: `http://localhost:3002`
- Postgres: `localhost:5432`

### Run Without Docker

Backend:

```bash
go run ./cmd/webhook-inbox
```

Frontend:

```bash
cd web
npm ci
npm run dev
```

## Tests

Backend:

```bash
go test ./...
```

Frontend:

```bash
cd web
npm test -- --run
```
