# Web UI

This is the SvelteKit frontend for Webhook Inbox. It provides pages for listing sources, viewing captured events, inspecting payloads, and checking delivery attempts.

## Project Structure

```text
web/
├── src/
│   ├── lib/                 # shared UI helpers, API utilities, components
│   ├── routes/              # SvelteKit pages and route components
│   └── test/                # frontend test setup and mocks
├── static/                  # static assets
├── package.json
├── vite.config.ts
└── svelte.config.js
```

## How To Run

### Local

```bash
npm ci
npm run dev
```

By default the UI runs on `http://localhost:5173`.

### With Docker Compose

From the repo root:

```bash
cp .env.example .env
```

```bash
docker compose --profile dev up --build web
```

## Useful Commands

Run tests:

```bash
npm test
```

Run type checks:

```bash
npm run check
```
