# GoNectar

GoNectar is a small honeypot-style HTTP trap written in Go. It listens for incoming HTTP requests, records basic request metadata and bodies into a newline-delimited JSON file (`events.jsonl`), and exposes a tiny HTML login page for interaction testing.

This repository contains a minimal collector, configuration, and an HTTP trap module that ingests events for later analysis.

## Features

- Lightweight HTTP honeypot listening on a configurable address (default `:8080`).
- Records request type, timestamp, metadata and request body to `events.jsonl`.
- Simple built-in login page at `/login` to attract credential submission attempts.
- Graceful shutdown handling for both the HTTP server and the collector file.

## Repository layout

- `cmd/honeypot/` - application entry point that wires together components and starts the trap.
- `internal/trap/` - HTTP trap implementation.
- `internal/collector/` - event collector that writes JSONL events to `events.jsonl`.
- `internal/config/` - configuration types.
- `internal/collector/collector.go` - collector implementation and `Event` definition.
- `web/` - static assets (e.g. `index.html`) used by or alongside the honeypot.

## Event format

Events are stored in `events.jsonl` (one JSON object per line). Each event has the shape:

{
"Type": "http.request",
"Time": "2025-...T...Z",
"Date": {
"method": "GET",
"path": "/",
"remote": "1.2.3.4:56789",
"ua": "User-Agent string",
"headers": { ... },
"body": "..."
}
}

Note: the collector uses the `Date` field name (not `Data`) for arbitrary event data.

## Quick start (build & run)

Requirements: Go 1.20+ (use the version specified in your environment if different).

1. Clone the repository (if you haven't already) and change to the repository root:

```powershell
git clone https://github.com/DeveloperAromal/GoNectar.git
cd GoNectar
```

2. Build the honeypot binary:

```powershell
go build -o bin/gonectar ./cmd/honeypot
```

3. Run the binary (this will create/append to `events.jsonl` in the working directory):

```powershell
.\bin\gonectar.exe
```

or

```powershell
.\bin\gonectar
```

4. Visit `http://localhost:8080/` or `http://localhost:8080/login` to generate events. Observe `events.jsonl` updating.

Alternatively you can run directly with `go run`:

```powershell
go run ./cmd/honeypot
```

## Configuration

Configuration is provided by constructing a `config.Config` value in `cmd/honeypot/main.go`. By default the code sets `HTTPAddr: ":8080"`. You can change this before building or modify the code to load configuration from a file or environment variables.

## Development notes

- The HTTP trap creates events by calling `collector.IngestEvent`. The `Event` struct uses `Date map[string]interface{}` to carry arbitrary payload data — ensure you use `Date` (not `Data`) when creating events.
- The collector writes events to `events.jsonl` using `json.Marshal` and appends a newline. The collector logs a short message after ingestion.
- Shutdown: the application listens for SIGINT/SIGTERM and attempts a graceful shutdown: it stops the HTTP server and then stops the collector (which closes the file).

## Common tasks

- Running tests: there are no tests currently included. Add Go test files under packages and run `go test ./...`.
- Add fields to events: update `collector.Event` and ensure all locations that construct events use the `Date` field.

## Troubleshooting

- If `go build ./...` fails with `undefined: config.Config`, check that `internal/config/nectar.config.go` exports `Config` (capitalized).
- If you see `unknown field Data in struct literal of type collector.Event`, change `Data:` to `Date:` in the code creating the event.
- If the program can't bind to `:8080`, ensure nothing else is listening on that port or change `HTTPAddr` to an available port.

## License

This project is licensed under the MIT License — see the `LICENSE` file for details.
