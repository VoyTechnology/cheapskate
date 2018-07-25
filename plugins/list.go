// This file contains the list of plugins and their use

package plugins

var plugins = []Plugin{
	new(testerIntegrationPlugin),
	new(testerCommandPlugin),
	new(angryPlugin),
	new(appendPlugin),
}
