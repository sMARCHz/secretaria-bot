package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetLogger(t *testing.T) {
	testcases := []struct {
		it       string
		logger   Logger
		expected Logger
	}{
		{
			it:       "set logger as provided",
			logger:   zap.New(nil).Sugar(),
			expected: zap.New(nil).Sugar(),
		},
		{
			it:       "set logger as noop logger when logger is not provided",
			logger:   nil,
			expected: &noopLogger{},
		},
	}

	for _, tc := range testcases {
		originalLogger := logger
		t.Run(tc.it, func(t *testing.T) {
			assert.Equal(t, &noopLogger{}, logger)

			SetLogger(tc.logger)

			assert.Equal(t, tc.expected, logger)
		})
		logger = originalLogger
	}
}
