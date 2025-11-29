package config

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCheckMissingConfig(t *testing.T) {
	testcases := []struct {
		it       string
		setup    func()
		expected error
	}{
		{
			it: "returns nil if every required config are loaded",
			setup: func() {
				viper.Set("app.port", "80")
				viper.Set("line.user_id", "line_uid")
				viper.Set("line.channel_secret", "line_secret")
				viper.Set("line.channel_token", "line_token")
				viper.Set("finance_url", "127.0.0.1:8080")
			},
			expected: nil,
		},
		{
			it: "returns error if any required config isn't loaded",
			setup: func() {
				viper.Set("app.port", "80")
			},
			expected: errors.New("line.user_id is missing in the config"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.it, func(t *testing.T) {
			tc.setup()
			defer viper.Reset()

			err := checkMissingConfig()
			if tc.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expected.Error())
			}
		})
	}
}
