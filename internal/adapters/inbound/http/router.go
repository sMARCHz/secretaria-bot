package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/client/finance"
	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/inbound/http/line"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/services"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	service := services.NewBotService(finance.NewFinanceServiceClient())
	lineHandler := line.NewLineHandler(service)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	router.POST("/line", func(ctx *gin.Context) {
		lineHandler.HandleLineMessage(ctx)
	})

	return router
}
