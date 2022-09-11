package config

import (
	"strings"

	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
	"github.com/spf13/viper"
)

type Configuration struct {
	App               AppConfiguration
	Line              LineMessageConfiguration
	FinanceServiceURL string `mapstructure:"finance-url"`
}

type AppConfiguration struct {
	Port string
}

type LineMessageConfiguration struct {
	ChannelSecret string
	ChannelToken  string
}

func LoadConfig(path string, logger logger.Logger) Configuration {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.BindEnv("line.channelsecret", "CHANNEL_SECRET"); err != nil {
		logger.Fatal("failed to bind CHANNEL_SECRET env: ", err)
	}
	if err := viper.BindEnv("line.channeltoken", "CHANNEL_TOKEN"); err != nil {
		logger.Fatal("failed to bind CHANNEL_TOKEN env: ", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("failed to load configuration: ", err)
	}

	var configuration Configuration
	if err := viper.Unmarshal(&configuration); err != nil {
		logger.Fatal("failed to unmarshal configuration: ", err)
	}
	return configuration
}
