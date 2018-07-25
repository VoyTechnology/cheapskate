// Package plugins defines the plugin API
package plugins

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/cpssd-students/cheapskate/hook"
	"github.com/cpssd-students/cheapskate/settings"
	"github.com/golang/glog"
)

// Errors
var (
	ErrPluginRegistered = errors.New("plugin already registered")
	ErrPluginNoExist    = errors.New("plugin does not exist")
)

// Type defines the type of the plugin. Since everything in the system is
// handled with them, this the main entrypoint
type Type string

const (
	// IntegrationType handles all external integrations, such as Slack or
	// Messenger
	IntegrationType Type = "integration"

	// CommandType respond to specific commands, such as /ping and they
	// return a response
	CommandType = "command"

	// RegexType handles substrings in the message sent
	RegexType = "regex"
)

func (t Type) String() string {
	return string(t)
}

// Get the plugin based on its name.
func Get(name string) (Plugin, error) {
	p, exists := registeredPlugins[name]
	if exists {
		return p, nil
	}

	return nil, ErrPluginNoExist
}

// Plugin defines the general interface the plugins must use
type Plugin interface {
	// Type of the plugin
	Type() Type

	// name of the plugin
	Name() string

	// Authors name, specified as `Name <Email>`.
	Authors() []string

	// Register registers specific functions based on the `Type` of the plugin.
	//	CommandType - /command, where only command has to be specified, eg. ping
	//	RegexType - `([A-Z])\w+`, it must be valid or error would be returned
	//	IntegrationType - registers as subpath of the webhook request.
	Register() string

	// Do is send to actually do the request.
	Do(context.Context, *Action) error

	// Used for integration
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// Action specifies the action to be taken
type Action struct {
	// Origin specifies where the action originates from
	Origin Plugin

	// Target is specified when the action has a specific target and should not
	// go anywhere else.
	Target Plugin

	// Internal command that the plugin has to do
	Command string

	// Data of the action
	Data []byte

	// Response to the action. Given as a channel when something is done
	Response chan []byte
}

// Load the plugins
func Load() error {
	if !settings.Get("plugin.enable").(bool) {
		glog.Info("Not loading any plugins")
		return nil
	}

	disabled := make(map[string]struct{})
	for _, d := range strings.Split(settings.Get("plugin.disabled").(string), ",") {
		disabled[d] = struct{}{}
	}

	for _, p := range plugins {
		if _, exists := registeredPlugins[p.Name()]; exists {
			return ErrPluginRegistered
		}

		if _, isDisabled := disabled[p.Name()]; isDisabled {
			glog.Infof("Not loading plugin %s - disabled", p.Name())
		}

		switch p.Type() {
		default:
			glog.Info("registering integration")
			if err := hook.RegisterHandleFunc(p.Register(), p); err != nil {
				return err
			}
		}

		glog.Infof("Registering plugin %s of type %s", p.Name(), p.Type())
		registeredPlugins[p.Name()] = p
	}

	return nil
}

// AddAction to the action stream
func AddAction(a *Action) {
	glog.Infof("received action from %s", a.Origin.Name())
	a.Response <- []byte("done")
}
