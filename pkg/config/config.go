// Package config allows to load the configuration for cheapskate
package config

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	prefix = "CHEAPSKATE"
)

// Config defines the configuration options available for the application
type Config struct {
	Server Server `mapstructure:"server"`
}

// Server configuration options
type Server struct {
	Port int `mapstructure:"port"`
}

// Plugin configiguration options
type Plugin struct {
	Options map[string]map[string]interface{} `mapstructure:"options"`
}

// Get gets the configuration from enviromental variables and exposes them as
// internal configuration complete with default values
func Get() (Config, error) {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix(prefix)
	v.AutomaticEnv()

	v = setDefaults(v)

	var c Config
	if err := v.UnmarshalExact(&c); err != nil {
		return Config{}, errors.Wrap(err, "unable to unmarshal configuration")
	}

	return c, nil
}

func setDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault("server.port", 8080)

	return v
}
