package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/cpssd-students/cheapskate/hook"
	"github.com/cpssd-students/cheapskate/plugins"
	"github.com/cpssd-students/cheapskate/settings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}
}

func run() error {

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})

	log.Info().Msg("starting cheapskate")

	cfg, err := settings.New()
	if err != nil {
		return errors.Wrap(err, "unable to get settings")
	}

	plog := log.With().Str("component", "plugins").Logger()
	if err := plugins.Load(plog, cfg.Plugin); err != nil {
		return errors.Wrap(err, "unable to load plugins")
	}

	hlog := log.With().Str("component", "hook").Logger()
	if err := hook.Run(hlog, cfg.Webhook); err != nil {
		return errors.Wrap(err, "unable to start webhooks")
	}

	return nil
}

func handleFunc(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("OK"))
}
