package rest

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

func Start(config config.Configuration, logger logger.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.App.Port),
		Handler: buildHandler(config, logger),
	}

	// Start server
	go func() {
		logger.Infof("Listening and serving HTTP on :%v", config.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Cannot start server: ", err)
		}
	}()

	// Shutdown server
	<-ctx.Done()
	stop()

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}
	logger.Info("Gracefully shutting down...")
}

func buildHandler(config config.Configuration, logger logger.Logger) *gin.Engine {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	router.GET("/line", func(ctx *gin.Context) {
		handleLineMessage(ctx, config.Line, logger)
	})

	return router
}
