package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfiguration() {
	viper.AutomaticEnv()

	server := viper.GetString("SERVER")
	cloud := viper.GetString("CLOUD")

	if server == "" && cloud == "LOCAL" {
		viper.Set("SERVER", "LOCAL")
	} else if server == "" {
		viper.SetDefault("SERVER", "AWS")
	}

	if cloud == "" && viper.GetString("SERVER") == "LOCAL" {
		viper.Set("CLOUD", "LOCAL")
	} else if cloud == "" {
		viper.SetDefault("CLOUD", "AWS")
	}

	viper.SetDefault("APPLICATION", "API")

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Configuration Initialized")
}

func IsAWS() bool {
	return viper.GetString("SERVER") == "AWS"
}

func IsLocalStack() bool {
	return viper.GetString("CLOUD") == "LOCAL"
}
