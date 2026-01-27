package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"example.com/rest-core-arch/pkg/config"
	"example.com/rest-core-arch/pkg/domain/example"
	middlewareZap "example.com/rest-core-arch/pkg/middleware"
	"example.com/rest-core-arch/pkg/util/log"
	"example.com/rest-core-arch/pkg/util/worker"
)

type Server struct {
	echo       *echo.Echo
	config     *config.AppConf
	ctx        context.Context // Root context
	workerPool *worker.Pool
}

func New(conf *config.AppConf) *Server {
	return &Server{
		echo:   echo.New(),
		config: conf,
	}
}

// Setup initializes the server components with the provided context
func (s *Server) Setup(ctx context.Context) {
	s.ctx = ctx
	s.setupMiddleware()
	s.setupRoutes()
}

func (s *Server) setupMiddleware() {
	// 1. Recovery - MUST be first to catch panics
	s.echo.Use(middleware.Recover())

	// 2. Context Propagation - Inject root context
	s.echo.Use(middlewareZap.ContextPropagator(s.ctx))

	// 3. Request ID & Logging
	// ZapLogger middleware handles Request ID generation and logging
	s.echo.Use(middlewareZap.ZapLogger())

	// 4. CORS
	s.applyCORSMiddleware()
}

func (s *Server) applyCORSMiddleware() {
	cfg := s.config.Server.CORS
	if !cfg.Enabled {
		return
	}

	corsConfig := middleware.CORSConfig{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		ExposeHeaders:    cfg.ExposeHeaders,
		MaxAge:           cfg.MaxAge,
	}

	// Safe defaults
	if len(corsConfig.AllowOrigins) == 0 {
		corsConfig.AllowOrigins = []string{"*"}
	}
	if len(corsConfig.AllowMethods) == 0 {
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}
	if len(corsConfig.AllowHeaders) == 0 {
		corsConfig.AllowHeaders = []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderXRequestID,
		}
	}

	s.echo.Use(middleware.CORSWithConfig(corsConfig))
}

func (s *Server) setupRoutes() {
	s.echo.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Initialize Worker Pool
	s.workerPool = worker.NewPool(s.ctx, s.config.Worker.Count, s.config.Worker.QueueSize)
	s.workerPool.Start()

	// Initialize Domain
	svc := example.NewService(s.workerPool)
	h := example.NewHandler(svc)

	// Register Routes
	v1 := s.echo.Group("/api/v1")
	h.RegisterRoutes(v1)
}

// Run starts the server and waits for a shutdown signal
func (s *Server) Run() {
	// 1. Create root context with signal handling
	// This context will be cancelled when SIGINT or SIGTERM is received.
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// 2. Setup with context propagation
	s.Setup(ctx)

	log.L().Info("Server starting", zap.Int("port", s.config.Server.Port))

	// 3. Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%d", s.config.Server.Port)
		if err := s.echo.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// 4. Wait for shutdown signal OR server error
	select {
	case <-ctx.Done():
		log.L().Info("Shutdown signal received", zap.Error(ctx.Err()))
	case err := <-serverErr:
		log.L().Error("Server failed to start", zap.Error(err))
		return
	}

	// 5. Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.L().Error("Shutdown error", zap.Error(err))
	} else {
		log.L().Info("Server shutdown completed successfully")
	}
}

// Shutdown stops the HTTP server and other components
func (s *Server) Shutdown(ctx context.Context) error {
	var errs []error

	// 1. Stop accepting NEW HTTP requests first
	log.L().Info("shutdown: stopping HTTP server")
	if err := s.echo.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("HTTP server shutdown: %w", err))
	}

	// 2. Stop worker pools
	log.L().Info("shutdown: stopping worker pools")
	if s.workerPool != nil {
		if err := s.workerPool.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("worker pool shutdown: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	return nil
}
