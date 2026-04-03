# Dev Tools

This directory contains small helper tools used during local development. Right now it includes a simple Python HTTP server that can act as a fake webhook receiver.

## Project Structure

```text
dev-tools/
├── Dockerfile              # container image for the helper server
└── server.py               # simple HTTP server for testing requests
```

## How To Run

### With Docker Compose

From the repo root:

```bash
docker compose --profile dev up --build dev-server
```

The server is exposed at `http://localhost:3002`.

### Local

```bash
cd dev-tools
python3 server.py
```

## What It Does

- logs request method, headers, query params, and body
- returns `200 OK` by default
- returns `400` on `/4xx`
- returns `500` on `/5xx`
