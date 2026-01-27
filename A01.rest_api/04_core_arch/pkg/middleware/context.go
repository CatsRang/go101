package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
)

// ContextPropagator returns a middleware that merges the server's root context
// with the request context. This ensures that request handlers are notified
// when the server is shutting down (rootCtx cancelled) or when the client
// disconnects (request context cancelled).
func ContextPropagator(rootCtx context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if rootCtx == nil {
				return next(c)
			}

			req := c.Request()
			originalCtx := req.Context()

			// Create a child context from the root context (application lifecycle)
			// This automatically propagates rootCtx cancellation to mergedCtx.
			mergedCtx, cancel := context.WithCancel(rootCtx)
			defer cancel()

			// Monitor the original request context (client disconnect)
			// If the client disconnects, we cancel our merged context.
			// We use a goroutine because we need to wait on select, but we can't block
			// the main thread here before calling next().
			//
			// Note: This goroutine might leak if next(c) hangs forever and neither context cancels,
			// but in practice HTTP handlers have timeouts or return, triggering defer cancel(),
			// which closes mergedCtx.Done(), allowing this goroutine to exit (via the second case
			// if we added it, or just because mergedCtx is canceled).
			//
			// Wait... if mergedCtx is cancelled by defer cancel(), we should ensure the goroutine exits.
			// The current select only listens to originalCtx and mergedCtx.
			// If defer cancel() runs, mergedCtx.Done() closes.
			go func() {
				select {
				case <-originalCtx.Done():
					// Client disconnected
					cancel()
				case <-mergedCtx.Done():
					// Request finished (defer cancel ran) or Server shutdown
				}
			}()

			// Replace the request context with our merged context
			c.SetRequest(req.WithContext(mergedCtx))

			return next(c)
		}
	}
}
