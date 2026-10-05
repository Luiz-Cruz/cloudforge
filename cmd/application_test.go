package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestProvideRunner(t *testing.T) {
	t.Run("Local Fiber Application", func(t *testing.T) {
		viper.Set("SERVER", "LOCAL")
		viper.Set("APPLICATION", "API")

		runner := ProvideRunner()
		_, ok := runner.(FiberApplication)
		assert.True(t, ok)
	})

	t.Run("Local Job Application", func(t *testing.T) {
		viper.Set("SERVER", "LOCAL")
		viper.Set("APPLICATION", "JOB")

		runner := ProvideRunner()
		_, ok := runner.(LocalJobApplication)
		assert.True(t, ok)
	})

	t.Run("Local Cron Application", func(t *testing.T) {
		viper.Set("SERVER", "LOCAL")
		viper.Set("APPLICATION", "CRON")

		runner := ProvideRunner()
		_, ok := runner.(LocalCronApplication)
		assert.True(t, ok)
	})

	t.Run("AWS Lambda API Application", func(t *testing.T) {
		viper.Set("SERVER", "AWS")
		viper.Set("APPLICATION", "API")

		runner := ProvideRunner()
		_, ok := runner.(LambdaApplication)
		assert.True(t, ok)
	})

	t.Run("AWS Job Application", func(t *testing.T) {
		viper.Set("SERVER", "AWS")
		viper.Set("APPLICATION", "JOB")

		runner := ProvideRunner()
		_, ok := runner.(JobApplication)
		assert.True(t, ok)
	})

	t.Run("AWS Cron Application", func(t *testing.T) {
		viper.Set("SERVER", "AWS")
		viper.Set("APPLICATION", "CRON")

		runner := ProvideRunner()
		_, ok := runner.(CronApplication)
		assert.True(t, ok)
	})

	// Reset
	viper.Set("SERVER", "LOCAL")
	viper.Set("APPLICATION", "API")
}
