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
			logger.Fatal("Cannot start a server: ", err)
		}
	}()

	// Shutdown: listen for interrupt/terminate signals (SIGKILL cannot be caught)
	sigCtx, sigCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-sigCtx.Done()
	sigCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Forcefully shutting down: ", err)
	}
	logger.Info("Gracefully shutting down...")
}
