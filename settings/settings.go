// Package settings allows the system to read settings for the running
// application
package settings

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// New parses the settings and returns them
func New() (Settings, error) {
	s := viper.New()
	setDefaults(s)

	s.AutomaticEnv()
	s.SetEnvPrefix("CHEAPSKATE")

	var c Settings
	if err := s.Unmarshal(&c); err != nil {
		return Settings{}, errors.Wrap(err, "unable to parse settings")
	}

	return c, nil
}

func setDefaults(s *viper.Viper) {
	s.SetDefault("webhook.port", 8080)
	s.SetDefault("plugin.enable", true)
	s.SetDefault("plugin.disabled",
		"tester_integration,tester_command")
}
