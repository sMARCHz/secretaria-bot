package line

import (
	"fmt"
	"testing"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sMARCHz/secretaria-bot/internal/logger"
	"github.com/sMARCHz/secretaria-bot/test/mocks"
	"github.com/sMARCHz/secretaria-bot/test/sandbox"
	"github.com/stretchr/testify/assert"
)

func TestNewLineHandler(t *testing.T) {
	sandbox.Run(t)
	bot := mocks.NewMockBotService(t)

	handler := NewLineHandler(bot)

	assert.IsType(t, &LineHandler{}, handler)
	assert.Equal(t, bot, handler.service)
	assert.IsType(t, &linebot.Client{}, handler.client)
	assert.NotEmpty(t, handler.client)
}

func TestNewLineHandler_Error(t *testing.T) {
	testLogger := &testLogger{}
	logger.SetLogger(testLogger)
	bot := mocks.NewMockBotService(t)

	NewLineHandler(bot)

	assert.True(t, testLogger.called)
	assert.Equal(t, "cannot create linebot client: missing channel secret", testLogger.msg)
}

type testLogger struct {
	called bool
	msg    string
}

func (t *testLogger) Debug(args ...interface{}) {}
func (t *testLogger) Info(args ...interface{})  {}
func (t *testLogger) Warn(args ...interface{})  {}
func (t *testLogger) Error(args ...interface{}) {}
func (t *testLogger) Fatal(args ...interface{}) {
	msg := ""
	for _, arg := range args {
		msg += fmt.Sprintf("%v", arg)
	}
	t.msg = msg
	t.called = true
}
func (t *testLogger) Debugf(format string, args ...interface{}) {}
func (t *testLogger) Infof(format string, args ...interface{})  {}
func (t *testLogger) Warnf(format string, args ...interface{})  {}
func (t *testLogger) Errorf(format string, args ...interface{}) {}
func (t *testLogger) Fatalf(format string, args ...interface{}) {}
