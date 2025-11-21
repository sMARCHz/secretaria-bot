package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	httpapi "github.com/sMARCHz/go-secretaria-bot/internal/adapters/inbound/http"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

func StartHTTPServer() {
	// Start server
	cfg := config.Get()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.App.Port),
		Handler: httpapi.NewRouter(),
	}
	go func() {
		logger.Infof("Listening and serving HTTP on :%v", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("cannot start server: ", err)
		}
	}()

	// Shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-ctx.Done()
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown: ", err)
	}
	logger.Info("Gracefully shutting down...")
}
