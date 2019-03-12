// Package plugin ties in all plugins.
package plugin

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/cpssd-students/cheapskate/pkg/config"
	"github.com/cpssd-students/cheapskate/pkg/server"
)

// Error from plugin package
type Error string

func (e Error) Error() string {
	return string(e)
}

// Exported errors
const (
	ErrPluginNameUsed    = Error("plugin name already registered")
	ErrBadImplementation = Error("bad plugin implementation")
)

func init() {
	ph = handler{
		log: zerolog.Nop(),

		names:       make(map[string]Plugin),
		integration: []Plugin{},
		action:      []Plugin{},
	}
}

type handler struct {
	log zerolog.Logger
	s   *server.Server
	cfg config.Plugin

	names       map[string]Plugin
	integration []Plugin
	action      []Plugin
}

// Init the plugin handler with a list of options
func Init(opts ...Option) error {
	for _, o := range opts {
		if err := o.apply(&ph); err != nil {
			return errors.Wrap(err, "unable to apply plugin handler options")
		}
	}
	return nil
}

var ph handler

// Option for the plugin handler
type Option interface {
	apply(*handler) error
}

type optionFunc func(*handler) error

func (of optionFunc) apply(h *handler) error {
	return of(h)
}

// Config is passed to each plugin when it is initialized. It is aquired from
// the env variables.
type Config map[string]interface{}

// Type of plugin
type Type string

// Types of plugins
const (
	Integration Type = "integration"
	Action           = "action"
)

// Plugin must implement these basic options in order to function propertly.
type Plugin interface {
	// Name of the plugin that is supposed to be registered. They have to be
	// unique.
	Name() string

	// Type of the plugin
	Type() Type

	// Command (or endpoint to register)
	Command() string

	// Init is called right when the plugin calls registation
	Init(zerolog.Logger, Config)

	Do(*Event) error
}

// Event that happened and has to be handled by the plugin
type Event struct {
}

// Register the plugin and allow to do some action
func Register(p Plugin) error {
	if err := ph.add(p); err != nil {
		return errors.Wrap(err, "unable to register a plugin")
	}
	return nil
}

func (h *handler) add(p Plugin) error {
	if _, exists := h.names[p.Name()]; exists {
		return ErrPluginNameUsed
	}
	h.names[p.Name()] = p

	if p.Type() == Action {
		l := h.log.
			With().
			Str("plugin", p.Name()).
			Str("type", string(p.Type())).
			Logger()
		h.action = append(h.action, p)
		p.Init(l, h.cfg.Options[p.Name()])
		return nil
	}

	// otherwise we can check the integration plugins do they implement
	// http.Handler

	hh, ok := p.(http.Handler)
	if !ok {
		return errors.Wrap(ErrBadImplementation,
			"plugin does not implement http.Handler")
	}
	if err := h.s.Register(p.Command(), hh); err != nil {
		return errors.Wrap(err, "unable to register plugin handler")
	}
	h.integration = append(h.integration, p)
	return nil
}
