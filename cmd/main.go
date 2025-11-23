package main

import (
	"github.com/sMARCHz/secretaria-bot/internal/infrastructure"
	"github.com/sMARCHz/secretaria-bot/internal/logger"
)

func main() {
	sync := logger.InitLogger()
	defer sync()

	infrastructure.StartHTTPServer()
}
