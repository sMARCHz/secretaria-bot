package main

import (
	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/driving/rest"
	"github.com/sMARCHz/go-secretaria-bot/internal/config"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

func main() {
	sync := logger.InitProductionLogger()
	defer sync()

	config := config.LoadConfig()

	rest.Start(config)
}
