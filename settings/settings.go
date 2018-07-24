// Package settings allows the system to read settings for the running
// application
package settings

import (
	"flag"

	"github.com/spf13/pflag"

	"github.com/spf13/viper"
)

var settings *viper.Viper

func init() {
	settings = viper.New()

	setDefaults()

	settings.SetConfigName("cheapskate")
	settings.AddConfigPath("/etc/cpssd")
	settings.AddConfigPath(".")
	setOptions()

	if err := settings.ReadInConfig(); err != nil {
		panic(err)
	}
}

func setDefaults() {
	settings.SetDefault("webhook.port", 80)
}

func setOptions() {
	flag.Bool("plugin.disabled", false, "disable plugins")

	flag.Parse()

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	settings.BindPFlags(pflag.CommandLine)
}

// Get a setting value
func Get(key string) interface{} {
	return settings.Get(key)
}

// Set a setting value
func Set(key string, value interface{}) {
	settings.Set(key, value)
}
