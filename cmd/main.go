package main

import (
	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/driving/rest"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

func main() {
	logger := logger.NewProductionLogger()
	config := config.LoadConfig(".", logger)

	rest.Start(config, logger)
}
