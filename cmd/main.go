package main

import (
	"github.com/sMARCHz/go-secretaria-bot/internal/adapters/inbound/http"
	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
)

func main() {
	sync := logger.InitLogger()
	defer sync()

	http.Start()
}
