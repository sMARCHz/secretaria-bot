package sandbox

import (
	"testing"

	"github.com/sMARCHz/secretaria-bot/internal/config"
)

func Run(t *testing.T) {
	setTestEnv(t)
	cleanup(t)
}

func setTestEnv(t *testing.T) {
	t.Setenv("LINE_CHANNEL_SECRET", "secret")
	t.Setenv("LINE_CHANNEL_TOKEN", "token")
	config.Get()
}

func cleanup(t *testing.T) {
	t.Cleanup(config.Reset)
}
