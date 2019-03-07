// Package hook creates a new webserver ready to listen on
package hook

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cpssd-students/cheapskate/settings"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Errors
var (
	ErrServerAlreadyStarted = errors.New("server already started, can't register plugin")
)

var started = false

// Run the webhook
func Run(log zerolog.Logger, cfg settings.Webhook) error {
	log.Info().Int("port", cfg.Port).Msg("starting web server")
	started = true
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}

// RegisterHandleFunc allows to register the integration
func RegisterHandleFunc(path string, f http.Handler) error {
	if started {
		return ErrServerAlreadyStarted
	}

	log.Info().Str("endpoint", path).Msg("registered")

	http.Handle(path, f)
	return nil
}
