package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfiguration() {
	viper.AutomaticEnv()

	viper.SetDefault("SERVER", "AWS")
	viper.SetDefault("APPLICATION", "API")
	viper.SetDefault("CLOUD", "AWS")

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Configuration Initialized")
}

func IsAWS() bool {
	return viper.GetString("SERVER") == "AWS"
}

func IsLocalStack() bool {
	return viper.GetString("CLOUD") == "LOCAL"
}
