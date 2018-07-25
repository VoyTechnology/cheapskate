// Package settings allows the system to read settings for the running
// application
package settings

import (
	"flag"
	"fmt"

	"github.com/spf13/pflag"

	"github.com/spf13/viper"

	// Make sure glog registers its flags
	_ "github.com/golang/glog"
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
		fmt.Println("--- CONFIG FILE NOT FOUND - USING DEFAULTS ---")
	}
}

func setDefaults() {
	settings.SetDefault("webhook.port", 8080)
	settings.SetDefault("plugin.enable", true)
	settings.SetDefault("plugin.disabled", "")
}

func setOptions() {
	flag.Bool("plugin.enable", true, "enable plugins")
	flag.String("plugin.disabled", "", "list of plugins to disable")
	flag.Int("webhook.port", 0, "webhook port")

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
