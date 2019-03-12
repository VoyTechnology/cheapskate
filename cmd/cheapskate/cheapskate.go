package main // import "github.com/cpssd-students/cheapskate"

import (
	"os"
	"time"

	"github.com/cpssd-students/cheapskate/pkg/config"
	"github.com/cpssd-students/cheapskate/pkg/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{
		TimeFormat: time.RFC3339,
		Out:        os.Stdout,
	}).With().Timestamp().Logger()

	if err := run(log); err != nil {
		log.Error().Msgf("%v", err)
		os.Exit(1)
	}
}

func run(log zerolog.Logger) error {
	cfg, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "unable to load config")
	}

	s, err := server.New(
		server.WithLogger(log.With().Str("component", "server").Logger()),
		server.WithConfig(cfg.Server),
	)
	if err != nil {
		return errors.Wrap(err, "unable to create server")
	}

	if err := s.Run(); err != nil {
		return errors.Wrap(err, "server failed to run")
	}
	return nil
}
