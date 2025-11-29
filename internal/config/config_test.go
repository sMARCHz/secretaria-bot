package config

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	originalDir, _ := os.Getwd()
	yaml := `
app:
  port: "8080"
line:
  channel_secret: "secret"
  channel_token: "token"
`
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/config.yaml"
	if err := os.WriteFile(tmpFile, []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}
	os.Chdir(tmpDir)
	t.Setenv("LINE_CHANNEL_SECRET", "secret")
	t.Setenv("LINE_CHANNEL_TOKEN", "token")

	config := Get()

	assert.Equal(t, "8080", config.App.Port)
	assert.Equal(t, "secret", config.Line.ChannelSecret)
	assert.Equal(t, "token", config.Line.ChannelToken)
	os.Chdir(originalDir)
}

func TestReset(t *testing.T) {
	loadOnce.Do(func() {
		data = Configuration{
			App: AppConfiguration{
				Port: "80",
			},
			FinanceServiceURL: "127.0.0.1:8080",
		}
	})

	Reset()

	assert.Equal(t, Configuration{}, data)
	called := false
	loadOnce.Do(func() { called = true })
	assert.True(t, called, "loadOnce was actually reset")
	loadOnce = sync.Once{}
}
