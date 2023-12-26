package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var requiredConfig = []string{
	"app.port",
	"line.user_id",
	"line.channel_secret",
	"line.channel_token",
	"finance_url",
}

func checkMissingConfig() error {
	for _, v := range requiredConfig {
		if cfg := viper.GetString(v); cfg == "" {
			return fmt.Errorf("%s is missing in the config", v)
		}
	}
	return nil
}
