# Analysis: Context-First Core REST Architecture

## 1. Project Goal
Create a reference implementation of a production-grade REST API in Go, focusing on **Context-First Architecture** and **Graceful Shutdown**. This project serves as a template for high-scale systems by implementing patterns described in the `go-context_first-rest-server-guide.md` and `go-echo-rest-dev-guide-with-context-mgmt.md`.

**Target Location**: `A01.rest_api/04_core_arch/10_context_pattern/`

## 2. Core Constraints & Stack
- **Language**: Go 1.20+
- **Web Framework**: Labstack Echo v4
- **CLI**: Cobra
- **Configuration**: Viper
- **Logging**: Zap (Structured, JSON)
- **Persistence**: **None** (Dummy business logic as requested, removing MyBatis/SQL dependencies).
- **Key Pattern**: Context propagation from OS signals down to individual workers.

## 3. Architecture Overview

### 3.1. Context Propagation Flow
The defining feature of this architecture is the single lineage of `context.Context`:

1.  **Root Context**: Created in `main.go` / `server.go` using `signal.NotifyContext`. Captures `SIGINT` and `SIGTERM`.
2.  **Server Injection**: This root context is passed into the `Server` struct.
3.  **Request Context (Middleware)**: A custom middleware wraps every HTTP request. It creates a child context that cancels if *either* the client disconnects *or* the Server's root context is canceled.
4.  **Service/Worker Layer**: All business logic and background workers accept `ctx` as the first argument. They immediately abort if `ctx.Done()` is closed.

### 3.2. Graceful Shutdown Sequence
Shutdown must happen in a specific order to prevent data loss or dropped requests:
1.  **Stop HTTP Server**: `echo.Shutdown(ctx)` stops accepting new connections.
2.  **Drain Worker Pools**: Signal background workers to finish current jobs and exit.
3.  **Close Resources**: (Simulated) Close database connections or file handles.

## 4. Component Design

### 4.1. Configuration (`pkg/config`)
- **Singleton Pattern**: `SharedAppConf()` to ensure config is loaded once.
- **Viper Integration**: Support for `config.yaml` and Environment variables (`APP_...`).
- **Structure**:
    - `Server`: Port, CORS settings.
    - `Log`: Level, Format, Rotation settings.
    - `Worker`: Pool size, Queue size.

### 4.2. Logging (`pkg/util/log`)
- **Zap Implementation**:
    - **Production**: JSON encoder, ISO8601 time, removal of sensitive fields.
    - **Development**: Console encoder, colorized.
- **Context Awareness**: Helper functions to extract/inject Trace IDs (`X-Request-ID`) from/to `context.Context`.

### 4.3. Middleware Stack (`pkg/middleware`)
Ordered execution is critical:
1.  **Recover**: Safety net.
2.  **Context Propagator**: Merges Root Context + Request Context.
3.  **Request ID**: Generates/Reads `X-Request-ID`.
4.  **Zap Logger**: Logs request ingress/egress with Trace ID.
5.  **CORS**: Security headers.

### 4.4. Domain Logic (`pkg/domain/example`)
- **Handler**: Binds HTTP inputs, calls Service.
- **Service**: Contains dummy business logic (e.g., `time.Sleep` to simulate work). Checks `ctx.Done()` to respect cancellation.
- **Model**: DTOs for Request/Response.

### 4.5. Worker Pool (`pkg/util/worker`)
- Generic implementation using `sync.WaitGroup` and buffered channels.
- **Key Feature**: Workers listen to the passed `ctx` for shutdown signals, ensuring they don't get "stuck" processing during a deployment.

## 5. Implementation Plan (Files)

| File Path | Responsibility |
|-----------|----------------|
| `cmd/main.go` | Entry point. Executes Root Command. |
| `pkg/cmd/root.go` | Initializes Config, Logger, and starts Server. |
| `pkg/server/server.go` | Orchestrates the `signal.NotifyContext`, Echo setup, and Shutdown sequence. |
| `pkg/config/app_conf.go` | Viper configuration definitions. |
| `pkg/middleware/context.go` | **CRITICAL**: Context merging logic. |
| `pkg/middleware/logger.go` | Zap middleware for Echo. |
| `pkg/util/log/logger.go` | Zap factory and configuration. |
| `pkg/util/worker/pool.go` | Reusable worker pool pattern. |
| `pkg/domain/example/` | Domain logic (Handler/Service) demonstrating context usage. |

## 6. Next Steps
1.  Scaffold the directory structure (Done).
2.  Implement `pkg/config` and `pkg/util/log` to establish the foundation.
3.  Implement `pkg/server` with the signal handling logic.
4.  Implement the Middleware chain.
5.  Create the Dummy Domain logic to verify the architecture works without a DB.
