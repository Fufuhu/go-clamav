# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

go-clamav is a Go service that scans S3-uploaded files for viruses via ClamAV. It long-polls an SQS queue that receives S3 event notifications, fetches the referenced object, streams it to clamd over TCP (`INSTREAM`), and persists the result to DynamoDB. Local development is fully containerized with S3-, SQS-, and DynamoDB-compatible substitutes.

## Common commands

Go version is pinned in `.go-version` (currently 1.22.3) — match it via goenv/asdf before building.

```bash
# Build the binary (output: ./go-clamav)
go build

# Bring up local infra: clamd + freshclam + minio + elasticmq + dynamodb (+ admin UIs)
# Re-run if DynamoDB tables aren't created on first boot — the aws-cli sidecars race the dynamodb container
docker compose up -d

# Run the poller locally against the docker compose stack
./exec_local.sh        # wraps `./go-clamav poll` with env vars for the local stack

# Run all tests
go test ./...

# Run a single package's tests
go test ./internal/virus_scan/scanner/...

# Run a single test by name
go test ./config -run TestGetConfig
```

Local UIs: minio `http://localhost:9001`, ElasticMQ `http://localhost:9325`, DynamoDB Admin `http://localhost:8001`.

To exercise the scan pipeline locally, upload `eicar.txt` and `test.txt` to the `test` bucket in MinIO, then send an S3-event-shaped SQS message to ElasticMQ — see README.md for the exact `aws sqs send-message` payloads (MinIO→ElasticMQ event wiring is not configured, so SQS must be poked manually).

## Architecture

### Runtime flow (single command: `go-clamav poll`)
1. `main.go` → `cmd.Execute()` (cobra) → `pollCmd` → `internal/cmd/poll.CommandPoll.Run`.
2. `Run` constructs four clients (SQS, S3, DynamoDB, clamd) from the `config.Configuration` and wires them into a `scanner.Scanner`.
3. `sqs.Client.Poll` loops forever: `ReceiveMessages` → for each, `process(...)` (= `scanner.Process`) → `DeleteMessage` on success. Errors are logged and the message is left for redelivery (no DLQ wiring in app code).
4. `scanner.Process`: filter by `SCANNING_TARGET_FILE_PATTERNS` regex (`message.IsTargetFile`) → `s3.GetObject` → `clamav.Scan` (streams chunks to clamd) → `dynamodb.PutScanResult`. Infected results are written to **both** `DYNAMODB_TABLE` and `DYNAMODB_TABLE_INFECTED`.

### Layering
- `cmd/` — cobra entrypoints. Thin; delegates to `internal/cmd/<subcommand>`. Add new subcommands by implementing `internal/cmd.CommandInterface` and registering in `cmd/`.
- `config/` — single `Configuration` struct populated from env vars via `kelseyhightower/envconfig`. `GetConfig()` memoizes; `Initialize()` clears the cache (used in tests).
- `internal/queue/clients/` — `QueueMessageInterface` is the abstraction the scanner consumes. `QueueMessage` is the concrete struct shared across packages (it's also embedded in `model.ScanResult`). The SQS `Client.Poll` parses the S3 event JSON inline and emits one `QueueMessage` per `Records[]` entry.
- `internal/virus_scan/clients/clamav/` — raw TCP `INSTREAM` protocol implementation (no third-party clamav lib). Chunk size is 1024 bytes; result is the raw clamd response string, compared to `ResultOK` (`"stream: OK\n"`).
- `internal/objects/clients/s3/`, `internal/db/clients/dynamodb/` — thin AWS SDK v2 wrappers. All three AWS clients honor a `*BaseUrl` config to redirect to local emulators (MinIO/ElasticMQ/DynamoDB Local).
- `internal/model/ScanResult` embeds `clients.QueueMessage`, so DynamoDB writes use the same Bucket/Key/ObjectPath methods as queue handling. The DynamoDB hash key is `ObjectPath` (`s3://bucket/key`), so re-scans of the same object **overwrite** the prior row.

### Sidecar / infra components (not Go code)
- `clamav/` and `freshclam/` — Dockerfiles for the clamd daemon and its signature updater. They share the `signatures` named volume in `compose.yaml`.
- `signature_downloader/` — separate ECS-targeted tool that mirrors ClamAV signatures to a private S3 bucket; used when `freshclam` is configured to pull from a private mirror instead of `database.clamav.net`. Has its own `docker-compose.yaml` and README.
