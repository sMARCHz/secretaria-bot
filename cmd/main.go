package main

import (
	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/driving/rest"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

func main() {
	sync := logger.InitLogger()
	defer sync()

	rest.Start()
}
