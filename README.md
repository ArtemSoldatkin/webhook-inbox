# Webhook Inbox

A minimal webhook inbox service built with Go and PostgreSQL.

This project allows developers to create unique HTTP endpoints that capture and store incoming webhook requests. It’s useful for inspecting, debugging, or replaying events from third-party services like Stripe, GitHub, or Twilio.

## ✨ Features

- Create endpoints with unique ingest keys
- Receive and store full webhook payloads + headers
- Stable event ordering by timestamp + ID
- Cursor-based pagination for event browsing
- Raw JSON body and headers stored in Postgres (JSONB)
- Fast, non-blocking ingestion (returns `202 Accepted`)

## 🛠 Stack

- **Go** (API)
- **PostgreSQL** (event storage)
- **Chi** (router)
- **SQLC** (DB layer)
- (Optional) Svetle + Vite for frontend

## 📦 Project Structure

```
webhook-inbox/
├── cmd/
│   └── webhook-inbox/       # main.go (entrypoint)
├── internal/                # all app logic, not exported
│   ├── api/                 # HTTP handlers (versioned if needed)
│   ├── config/              # app config
│   ├── db/                  # SQLC or data access logic
│   ├── ingest/              # validation + event persistence
│   ├── models/              # types (Event, Endpoint, etc.)
│   └── service/             # business logic of the app
├── migrations/              # SQL files for schema
├── web/                     # optional UI (React, HTMX, etc.)
├── .dockerignore
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── sqlc.yaml
```

## 🚧 Status

This is a portfolio/study project. It's not production-ready, but demonstrates how to build a simple webhook recorder from scratch.

## 📄 License

[MIT](./LICENSE)
