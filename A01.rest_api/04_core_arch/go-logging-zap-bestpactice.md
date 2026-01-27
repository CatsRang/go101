# Go Zap Logging Best Practices – Research Summary

The current “best case” for using Zap in Go services is a small set of stable patterns that show up repeatedly in mature guides and in real-world implementations like etcd.

---

## 1. Logger Construction and Configuration

Use a **single construction path** (usually `zap.Config`) and drive it via environment or flags.

Key points:

- **Build once** in `main`, not per request.
- **Production defaults**:
  - Encoding: `"json"` (machine-friendly, log-aggregator-friendly).
  - Level: `info` (raise to `debug` only when necessary).
  - EncoderConfig:
    - `EncodeTime: zapcore.ISO8601TimeEncoder`
    - `EncodeLevel: zapcore.LowercaseLevelEncoder`
    - `EncodeCaller: zapcore.ShortCallerEncoder`
- **Development config**:
  - Use `zap.NewDevelopment()` or `Encoding: "console"` for readable local logs.
- Always `defer logger.Sync()` in `main` to flush buffered logs.

Zap’s production config enables **sampling** by default (`100:100`) to avoid flooding.

**References**
- BetterStack: https://betterstack.com/community/guides/logging/go/zap/
- SigNoz: https://signoz.io/guides/zap-logger/
- Dash0: https://www.dash0.com/guides/logging-in-go-with-zap
- Official docs: https://pkg.go.dev/go.uber.org/zap

---

## 2. Structured Fields and Naming

The consensus is **“JSON with consistent fields”**:

- Always log with **structured fields** (`zap.String`, `zap.Int`, `zap.Duration`, `zap.Error`, etc.).
- Avoid `fmt.Sprintf` in messages; put data in fields.
- Pick and enforce a **field naming convention**:
  - etcd uses kebab-case: `local-member-id`, `cluster-id`, `snapshot-index`.
- Use `zap.Namespace` to group related fields (e.g., `request`, `user`).
- For errors, prefer `zap.Error(err)` or `zap.NamedError("apply-error", err)`.

This makes logs easy to filter and correlate in ELK, Loki, Splunk, etc.

**References**
- etcd logging doc: https://github.com/etcd-io/etcd (server logging docs / `config_logging.go`)
- Stackademic: https://blog.stackademic.com/building-observability-in-go-structured-logging-with-zap-513e1b32ec1f
- BetterStack: https://betterstack.com/community/guides/logging/go/zap/

---

## 3. Log Level Semantics

Common level usage guidelines are very consistent across sources and etcd:

- **Debug**: internal flow, large or noisy details; usually disabled in production.
- **Info**: lifecycle events, configuration, state transitions (“started server”, “completed job”).
- **Warn**: slow operations, transient failures with retries, degraded but working.
- **Error**: failed operations that affect behavior; include context and stack trace.
- **Fatal/Panic**: invariants broken; log once and terminate.

Best practices:

- Don’t spam debug logs in production; gate with level checks or environment.
- Avoid logging the same error many times up the stack; log where it’s actionable.

**References**
- etcd logging doc: https://github.com/etcd-io/etcd
- BetterStack: https://betterstack.com/community/guides/logging/go/zap/
- Last9: https://last9.io/blog/zap-logger/
- Leapcell (Go logging best practices): https://leapcell.io/blog/robust-go-best-practices-for-error-handling

---

## 4. Logger Propagation and Context

The modern “best case” is **no globals-only**, and no custom wrappers around Zap.

Recommended patterns:

- Build a root logger in `main()` and:
  - Call `zap.ReplaceGlobals(logger)` for convenience.
  - Also pass `*zap.Logger` into subsystems (`.Named("subsystem")`).
- Use **context-based propagation** where appropriate:
  - `NewContext(ctx, *zap.Logger)` / `FromContext(ctx) *zap.Logger` helpers.
  - HTTP middleware that adds request-scoped fields (request ID, method, path, remote IP) via `With()` before passing the context downstream.
- For observability stacks, attach OpenTelemetry **trace/span IDs** to logs from `context.Context`.

Anti-patterns called out explicitly:

- Don’t wrap `*zap.Logger` in your own type when not required.
- Don’t re-export `Info/Error` to your own package-level helpers just for style.

**References**
- etcd logging doc (propagation diagram, gRPC integration): https://github.com/etcd-io/etcd
- OpenTelemetry & context propagation patterns: https://nox.im/posts/2022/0327/observability-opentelemetry-context-propagation-in-go/
- Zap FAQ: https://github.com/uber-go/zap/blob/master/FAQ.md

---

## 5. Performance and Volume Control

Zap is already optimized (zero-allocation JSON encoder, sync.Pool, etc.), but best practice is to **avoid work around logging** too.

Patterns:

- Use `*zap.Logger` API (not `SugaredLogger`) in **hot paths** for fewer allocations.
- Use **level checks** for expensive payloads:

  ```go
  if ce := logger.Check(zap.DebugLevel, "user data"); ce != nil {
      ce.Write(zap.Any("data", expensiveUserData()))
  }
  ```

- Avoid logging inside tight loops; log periodically or sample.
- Let Zap’s default sampling (100:100) run, or tune it explicitly for very noisy messages.
- Use buffered / async writers (`BufferedWriteSyncer`) when I/O latency dominates.

**References**
- Zap design & performance: https://www.sobyte.net/post/2022-01/go-log-zap/
- Reddit discussions on logger performance: https://www.reddit.com/r/golang/comments/1181tfb/how_to_improve_logger_performances/
- Official zap docs & benchmarks: https://pkg.go.dev/go.uber.org/zap

---

## 6. Outputs, Rotation, and Integrations

Best practices for production deployments focus on **centralization and safety**:

- Prefer logging to **stdout/stderr** in containers; let the platform handle collection.
- For file-based logging:
  - Use `lumberjack` for rotation (size, age, backups, compression).
  - Or rely on external `logrotate` when integrating with system tooling.
- On systemd-based Linux:
  - Consider output to `systemd/journal` for central log management, rotation, and querying via `journalctl`.
  - Map Zap levels to journal priorities (debug/info/warn/error/crit).
- Use `zapcore.NewTee` to fan-out to multiple sinks (console + file + remote).

**References**
- SigNoz (lumberjack example): https://signoz.io/guides/zap-logger/
- StackOverflow – rotating files with Zap: https://stackoverflow.com/questions/45440491/how-to-configure-uber-go-zap-logger-for-rolling-filesystem-log
- etcd journal integration & rotation: https://github.com/etcd-io/etcd

---

## 7. Error Logging and Sensitivity

Error handling and security themes are very consistent:

- Log errors **once per failure path** with enough context:
  - Always include operation, key parameters, and correlation IDs.
- Avoid logging **sensitive data**:
  - Never log passwords, tokens, secrets, or full PII.
  - Use custom encoders or `String()` implementations to redact fields if needed.
- Use `Fatal` and `Panic` sparingly—only on unrecoverable conditions.

For infra events (like TLS handshake failures), down-level very noisy errors (e.g, EOF) to debug to avoid log noise.

**References**
- etcd TLS handshake logging and sensitivity: https://github.com/etcd-io/etcd
- Dash0 best practices: https://www.dash0.com/guides/logging-in-go-with-zap
- Leapcell error handling: https://leapcell.io/blog/robust-go-best-practices-for-error-handling

---

## 8. Sampling Strategies

Zap’s built-in **time-window sampling** is widely used and recommended:

- Default: first 100 events per second per (level, message), then 1 of 100.
- Custom:

  ```go
  samplingCore := zapcore.NewSamplerWithOptions(
      core,
      time.Second, // window
      5,           // first 5
      100,         // then 1/100
  )
  ```

Use sampling for:

- Repeated identical warnings (e.g., transient network failures).
- High frequency background jobs or health checks.
- Chatty debug logs in hot code paths.

**References**
- Zap sampling config (docs): https://pkg.go.dev/go.uber.org/zap
- Zap FAQ (sampling rationale): https://github.com/uber-go/zap/blob/master/FAQ.md
- BetterStack – log sampling & cost: https://betterstack.com/community/guides/logging/log-sampling/

---

## 9. Testing and Tooling

Recommended patterns for verifying logging behavior:

- Use `zaptest/observer` to assert on **messages, levels, and fields**:
  - Build a logger with an observed core and inspect `recorded.All()`.
- Use `zaptest.NewLogger(t)` for integration-style tests that should log to `testing.TB`.
- Use `zapcore.NewTee` to send logs to both your “real” core and an observer when you want to test a custom core’s behavior.

This allows you to treat logging as part of observable behavior without relying on brittle string parsing.

**References**
- zaptest docs: https://pkg.go.dev/go.uber.org/zap/zaptest
- Observer usage examples: https://stackoverflow.com/questions/70400426/how-to-properly-capture-zap-logger-output-in-unit-tests

---

## 10. High-Level “Best Case” Checklist

Synthesizing guides and large production systems like etcd into a short checklist:

- JSON logs in production; console in development.
- Single logger construction in `main` using `zap.Config`.
- Level, format, and outputs configurable via env/flags.
- `zap.ReplaceGlobals()` plus explicit `*zap.Logger` injection into subsystems.
- Request/context propagation using `With()` and `context.Context`.
- Consistent field naming (e.g., kebab-case) and type-safe fields.
- No sensitive data, no entire payload dumps in logs.
- Sampling enabled and tuned for noisy paths.
- Rotation configured (lumberjack, logrotate, or journald).
- Log at meaningful boundaries (startup, shutdown, RPC/HTTP edges, slow ops).
- Test logging behavior with `zaptest`/`observer`.
- Integrate with tracing (trace/span IDs in fields) for correlation.

**References**
- Zap GitHub: https://github.com/uber-go/zap
- BetterStack guide: https://betterstack.com/community/guides/logging/go/zap/
- Dash0 guide: https://www.dash0.com/guides/logging-in-go-with-zap
- SigNoz guide: https://signoz.io/guides/zap-logger/
- etcd logging implementation: https://github.com/etcd-io/etcd
