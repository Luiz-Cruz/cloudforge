package cmd

import (
	"github.com/Luiz-Cruz/cloudforge/platform/config"
	"github.com/spf13/viper"
)

type Application interface {
	Run()
}

func ProvideRunner() Application {
	if config.IsAWS() {
		if viper.GetString("APPLICATION") == "JOB" {
			return JobApplication{}
		}
		return LambdaApplication{}
	}
	return FiberApplication{}
}
