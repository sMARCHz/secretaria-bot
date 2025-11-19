package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/driven/financeservice"
	financetest "github.com/sMARCHz/go-secretaria-bot/internal/adapters/driving/rest/finance_test"
	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/driving/rest/line"
	"github.com/sMARCHz/go-secretaria-bot/internal/core/services"
)

func newRouter() *gin.Engine {
	router := gin.Default()
	service := services.NewBotService(financeservice.NewFinanceServiceClient())
	lineHandler := line.NewLineHandler(service)
	testHandler := financetest.NewTestHandler(service)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	router.POST("/line", func(ctx *gin.Context) {
		lineHandler.HandleLineMessage(ctx)
	})

	router.POST("/__test", func(ctx *gin.Context) {
		testHandler.Test(ctx)
	})

	return router
}
