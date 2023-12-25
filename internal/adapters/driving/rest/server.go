package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/driven/financeservice"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/dto"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/services"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

func Start(config config.Configuration) {
	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.App.Port),
		Handler: buildHandler(config),
	}
	go func() {
		logger.Infof("Listening and serving HTTP on :%v", config.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Cannot start server: ", err)
		}
	}()

	// Shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-ctx.Done()
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}
	logger.Info("Gracefully shutting down...")
}

func buildHandler(config config.Configuration) *gin.Engine {
	router := gin.Default()
	service := services.NewBotService(
		financeservice.NewFinanceServiceClient(config.FinanceServiceURL),
		config,
	)
	lineHandler := Handler{
		service: service,
		config:  config,
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	router.POST("/line", func(ctx *gin.Context) {
		lineHandler.handleLineMessage(ctx)
	})

	router.POST("/__test", func(ctx *gin.Context) {
		username, password, auth := ctx.Request.BasicAuth()
		if !auth || username != "sMARCHz's test" {
			logger.Warnf("someone tried to breach (username: %s, password: %s)", username, password)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var msg dto.TextMessageRequest
		if err := ctx.BindJSON(&msg); err != nil {
			logger.Error("cannot bind json: ", err)
		}
		res, err := service.HandleTextMessage(msg.Message)
		if err != nil {
			ctx.AbortWithError(err.StatusCode, errors.New(err.Message))
		} else {
			ctx.JSON(http.StatusOK, res)
		}
	})

	return router
}
