// Package plugins defines the plugin API
package plugins

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/cpssd-students/cheapskate/hook"
	"github.com/cpssd-students/cheapskate/settings"
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

// registedPlugins is generated when the `Load` function is called, and should
// not be called directly. Instead, use the `Get` function.
var registeredPlugins = make(map[string]Plugin)

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
	// Initialize the pluggin, passing in all the necessary basic information
	// TODO: Enable on all plugins
	// Init(log zerolog.Logger)

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
	Do(*Action) error

	// This function is called when no command was triggered. Its up to the
	// plugin to decide what to do with it.
	NoAction()

	// ServeHTTP must be present for
	// ServeHTTP(http.ResponseWriter, *http.Request)
}

var compiledRegex = make(map[string]*regexp.Regexp)
var commandPlugins = make(map[string]Plugin)

// Load the plugins
func Load(log zerolog.Logger, cfg settings.Plugin) error {
	if !cfg.Enable {
		log.Info().Msg("Not loading any plugins")
		return nil
	}

	disabled := make(map[string]struct{})
	for _, d := range cfg.Disabled {
		disabled[d] = struct{}{}
	}

	for _, p := range plugins {
		if _, exists := registeredPlugins[p.Name()]; exists {
			return ErrPluginRegistered
		}

		if _, isDisabled := disabled[p.Name()]; isDisabled {
			log.Info().Str("plugin", p.Name()).Msg("ignoring")
			continue
		}

		// p.Init(log.With().Str("plugin", p.Name()).Logger())

		log.Info().Str("plugin", p.Name()).Str("type", p.Type().String()).Msg("registered")

		switch p.Type() {
		case IntegrationType:
			handler, ok := p.(http.Handler)
			if !ok {
				return errors.Errorf("plugin %s does not implement ServeHTTP",
					p.Name())
			}

			if err := hook.RegisterHandleFunc(p.Register(), handler); err != nil {
				return errors.Wrap(err, "registering webhook")
			}
		case CommandType:
			if binded, exists := commandPlugins[p.Register()]; exists {
				return errors.Errorf("plugin %s tried to register command %s which is already binded to plugin %s",
					p.Name(), p.Register(), binded.Name())
			}
			commandPlugins[p.Register()] = p
		case RegexType:
			re, err := regexp.Compile(p.Register())
			if err != nil {
				return errors.Errorf("plugin %s tried to register invalid regex %s: %v",
					p.Name(), p.Register(), err)
			}
			compiledRegex[p.Name()] = re
		}

		registeredPlugins[p.Name()] = p
	}

	return nil
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

	// Meta contains the metadata used to communicate using KVs between
	// different. Totally depends on the implementation.
	Meta map[string]string

	// Response to the action. Given as a channel when something is done
	Response chan []byte

	// Marked when there is nothing else to do with the command.
	Done bool
}

// AddAction to the action stream
func AddAction(a *Action) {
	glog.V(2).Infof("data: %s", a.Data)

	done := false

	// TODO: This should probably be a channel with caching
	if a.Target != nil {
		glog.V(1).Infof("Sending directly to %s from %s",
			a.Target.Name(), a.Origin.Name())

		go a.Target.Do(a)
		lastAction = a
		done = true
		return
	}

	var action func(*Action) error

	for _, p := range registeredPlugins {
		if a.Done {
			break
		}

		if p.Name() == a.Origin.Name() {
			continue
		}

		switch p.Type() {
		case CommandType:
			if strings.HasPrefix(strings.ToLower(string(a.Data)), p.Register()) {
				glog.V(2).Infof("Sending action to command %s", p.Name())
				done = true
				action = p.Do
				break
			}
		case RegexType:
			if re := compiledRegex[p.Name()]; re.Match(a.Data) {
				glog.V(2).Infof("Sending action to regex %s", p.Name())
				done = true
				action = p.Do
				break
			}
		// This should mostly be handled by the individual target, but we catch
		// it just in case
		case IntegrationType:
			glog.V(2).Infof("Sending action to integration %s", p.Name())

			done = true
			action = p.Do
			break
		}
	}

	if !done {
		glog.V(2).Infof("Nothing was done, sending action back to origin %s",
			a.Origin.Name())
		action = a.Origin.Do
	}

	go func() {
		action(a)

		// I am not sure do we want to actually keep it like that. Perhaps we should
		// add action store and last action message with IDs or something.
		lastAction = a
	}()
}

// This is pretty hacky and a bad solution. Ideally this should be somehow
// saved in a store to ensure that we preserve the last message. This is only
// required for the last message, and is pretty bad. Let's not do it ever again.
var lastAction *Action

// TrimPrefix is a custom helper function which also removes the words if they
// are uppercase
func TrimPrefix(s, prefix string) string {
	pre := s
	s = strings.TrimPrefix(s, prefix)
	if pre == s {
		s = strings.TrimPrefix(s, strings.ToUpper(prefix))
	}

	return strings.TrimSpace(s)
}

// RegisterPlugin registers a single plugin
func RegisterPlugin(p Plugin) {
	plugins = append(plugins, p)
}

// RegisterPlugins registers a whole list of plugins
func RegisterPlugins(p []Plugin) {
	plugins = append(plugins, p...)
}

// PluginInfo describes the functions uses in the plugins. It is used to create
// quick and easy plugins without needing all the control of the normal plugin
type PluginInfo struct {
	PluginName      string
	PluginType      Type
	PluginAuthors   []string
	RegisterKeyword string
	Action          func(*Action) error
}

// Name of the plugin
func (p PluginInfo) Name() string {
	return p.PluginName
}

// Authors of the plugin
func (p PluginInfo) Authors() []string {
	return p.PluginAuthors
}

// Type of the plugin
func (p PluginInfo) Type() Type {
	return p.PluginType
}

// Register the plugin
func (p PluginInfo) Register() string {
	return p.RegisterKeyword
}

// Do plugin action
func (p PluginInfo) Do(a *Action) error {
	return p.Action(a)
}

// NoAction of the plugin
func (p PluginInfo) NoAction() {}
