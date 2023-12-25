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
	"github.com/line/line-bot-sdk-go/v7/linebot"
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
	financeClient := financeservice.NewFinanceServiceClient(config.FinanceServiceURL)
	lbot, err := linebot.New(config.Line.ChannelSecret, config.Line.ChannelToken)
	if err != nil {
		logger.Error("Cannot create new linebot: ", err)

	}
	botHandler := BotHandler{
		service: services.NewBotService(financeClient, config),
		config:  config,
		linebot: lbot,
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	router.POST("/line", func(ctx *gin.Context) {
		botHandler.handleLineMessage(ctx)
	})

	router.POST("/test/line", func(ctx *gin.Context) {
		botService := services.NewBotService(financeClient, config)
		var msg dto.TextMessageRequest
		if err := ctx.BindJSON(&msg); err != nil {
			logger.Error("Cannot bind json: ", err)
		}
		res, err := botService.HandleTextMessage(msg.Message)
		if err != nil {
			ctx.AbortWithError(err.StatusCode, errors.New(err.Message))
		} else {
			ctx.JSON(http.StatusOK, res)
		}
	})

	return router
}
