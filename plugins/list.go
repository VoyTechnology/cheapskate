// This file contains the list of plugins and their use

package plugins

var plugins []Plugin

func init() {
	RegisterPlugins([]Plugin{
		newTesterCommandPlugin(),
		new(testerIntegrationPlugin),
		newAngryPlugin(),
		newAppendPlugin(),
		newClapPlugin(),
		newUppercasePlugin(),
	})
}
