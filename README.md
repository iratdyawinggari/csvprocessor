# CSV Processor (Concurrent Data Processing in Go)

This project is a **concurrent CSV file processor** written in Go, designed as a take-home technical assessment. It demonstrates idiomatic Go concurrency, worker pool patterns, structured error handling, progress tracking, and testability.

---

## 1. Setup & Running Instructions

### Prerequisites
- Go **1.21+**
- Linux / macOS / WSL (recommended)

Verify Go installation:
```bash
go version
```

---

### Project Structure
```
Csv_Processor/
├── cmd/
│   └── csvproc/
│       └── main.go          # Application entry point
├── internal/
│   ├── reader/              # CSV reading logic
│   ├── worker/              # Worker pool implementation
│   ├── result/              # Result aggregation
│   └── progress/            # Progress tracking
├── data/                    # Sample CSV files
├── go.mod
├── go.sum
└── README.md
```

---

### Running the Application

```bash
go run ./cmd/csvproc \
  -dir ./data \
  -workers 4
```

Arguments:
- `-dir` : directory containing CSV files
- `-workers` : number of concurrent workers

---

### Running Tests

Run all tests with race detection:
```bash
go test ./... -race
```

---

## 2. Architecture Explanation

### High-Level Flow

```
CSV Files
   │
   ▼
[Reader Goroutines]
   │ (jobs channel)
   ▼
[Worker Pool]
   │ (results / errors)
   ▼
[Collector]
   │
   ▼
Summary Output
```

---

### Core Components

#### 1. Reader (`internal/reader`)
- Reads CSV files line-by-line
- Emits rows as `Row` structs into a **jobs channel**
- Increments progress tracker per row
- **Never closes shared channels** (ownership-safe)

#### 2. Worker Pool (`internal/worker`)
- Implements classic **fan-out / fan-in** pattern
- Each worker:
  - Listens on `jobs`
  - Processes data
  - Emits `result.Item` or `error`
  - Marks progress as done
- Uses `sync.WaitGroup` to signal worker completion

#### 3. Result Collector (`internal/result`)
- Consumes `results` and `errs`
- Aggregates summary:
  - Success count
  - Failure count
  - Error list
- Terminates only after both channels are closed

#### 4. Progress Tracker (`internal/progress`)
- Tracks total vs completed tasks
- Uses atomic counters + `sync.WaitGroup`
- Guarantees:
  - No negative counters
  - No deadlocks
  - Safe concurrent updates

---

### Channel Ownership Rules (Critical Design)

| Channel  | Owner        | Closed By |
|--------|-------------|-----------|
| jobs   | main         | main      |
| results| main         | main      |
| errs   | main         | main      |

Workers **never close channels**.
This avoids:
- `panic: close of closed channel`
- Hanging tests
- Race conditions

---

## 3. Technology Choices Justification

### Go Language
- Native concurrency (goroutines, channels)
- Strong tooling (`go test`, `-race`)
- Simple deployment (static binary)

### Worker Pool Pattern
- Prevents unbounded goroutines
- Controls memory usage
- Scales with CPU cores

### Channels over Mutexes
- Clear ownership semantics
- Easier reasoning for pipelines
- Idiomatic Go design

### Internal Package Layout
- Enforces encapsulation
- Prevents misuse from outside modules
- Clean separation of concerns

---

## 4. Known Limitations

1. **CSV Schema Validation**
   - Assumes consistent column counts
   - No schema inference

2. **Backpressure Control**
   - Channels are unbuffered by default
   - Large files may benefit from tuned buffers

3. **Progress Output**
   - Progress tracking is internal (no live CLI UI)

4. **Error Strategy**
   - Errors are collected, not retried

---

## 5. Future Improvements

- Add buffered channels with adaptive sizing
- Support context-based early termination on error threshold
- Live progress bar (TUI)
- Pluggable processors (strategy pattern)
- Benchmark suite (`go test -bench`)
- Optional JSON / Parquet output

---

## 6. Test Coverage Report

### Test Scope

| Package              | Coverage Focus |
|---------------------|---------------|
| `reader`            | CSV parsing, row count, channel behavior |
| `worker`            | Worker lifecycle, WaitGroup correctness |
| `progress`          | Increment / Done synchronization |
| `result`            | Aggregation correctness |

All tests include:
- Channel close validation
- Goroutine completion
- Deadlock prevention

### Run Coverage

```bash
go test ./... -cover
```

Example output:
```
        csvproc/cmd/csvproc             coverage: 0.0% of statements
ok      csvproc/internal/progress       1.006s  coverage: 95.0% of statements
        csvproc/internal/reader         coverage: 0.0% of statements
ok      csvproc/internal/result 0.003s  coverage: 100.0% of statements
        csvproc/internal/worker         coverage: 0.0% of statements
```

---