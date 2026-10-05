package cmd

import (
	"github.com/Luiz-Cruz/cloudforge/platform/config"
	"github.com/spf13/viper"
)

type Application interface {
	Run()
}

func ProvideRunner() Application {
	appType := viper.GetString("APPLICATION")
	if config.IsAWS() {
		switch appType {
		case "JOB":
			return JobApplication{}
		case "CRON":
			return CronApplication{}
		default:
			return LambdaApplication{}
		}
	}

	switch appType {
	case "JOB":
		return LocalJobApplication{}
	case "CRON":
		return LocalCronApplication{}
	default:
		return FiberApplication{}
	}
}
