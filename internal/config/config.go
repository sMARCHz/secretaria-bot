package config

import (
	"sync"

	"github.com/sMARCHz/go-secretaria-bot/internal/logger"
	"github.com/spf13/viper"
)

var (
	data     Configuration
	loadOnce sync.Once
)

type Configuration struct {
	App               AppConfiguration  `mapstructure:"app"`
	Line              LineConfiguration `mapstructure:"line"`
	FinanceServiceURL string            `mapstructure:"finance_url"`
}

type AppConfiguration struct {
	Port         string `mapstructure:"port"`
	TestUsername string `mapstructure:"test_username"`
}

type LineConfiguration struct {
	UserID        string `mapstructure:"user_id"`
	ChannelSecret string `mapstructure:"channel_secret"`
	ChannelToken  string `mapstructure:"channel_token"`
}

func Get() Configuration {
	loadOnce.Do(func() {
		data = loadConfig()
	})
	return data
}

func Reset() {
	loadOnce = sync.Once{}
}

func loadConfig() Configuration {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("failed to load configuration: ", err)
	}

	// ENV
	if err := viper.BindEnv("line.channel_secret", "LINE_CHANNEL_SECRET"); err != nil {
		logger.Fatal("failed to bind LINE_CHANNEL_SECRET env: ", err)
	}
	if err := viper.BindEnv("line.channel_token", "LINE_CHANNEL_TOKEN"); err != nil {
		logger.Fatal("failed to bind LINE_CHANNEL_TOKEN env: ", err)
	}
	if err := viper.BindEnv("app.test_username", "APP_TEST_USERNAME"); err != nil {
		logger.Fatal("failed to bind APP_TEST_USERNAME env: ", err)
	}

	if err := checkMissingConfig(); err != nil {
		logger.Fatal(err)
	}

	var configuration Configuration
	if err := viper.Unmarshal(&configuration); err != nil {
		logger.Fatal("failed to unmarshal configuration: ", err)
	}

	return configuration
}
