# Concurrent CSV Processor (Go)

## Overview
A concurrent CSV processing tool written in Go that:
- Reads multiple CSV files in parallel
- Uses a bounded worker pool
- Streams data to stay memory efficient
- Aggregates results and errors
- Displays real-time progress
- Shuts down cleanly and safely

## Run Locally
```bash
go mod init csvproc
go run ./cmd/csvproc --dir ./data --workers 4

## Testing

Run all tests with race detection:

```bash
go test ./... -race
```
