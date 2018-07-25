// This file contains the list of plugins and their use

package plugins

var plugins = []Plugin{
	new(corePlugin),
}

// registedPlugins is generated when the `Load` function is called, and should
// not be called directly. Instead, use the Get function.
var registeredPlugins = make(map[string]Plugin)
