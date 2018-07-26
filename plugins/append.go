// Adapted from steely/plugins/append.py by Wojtek Bednarzak

package plugins

func newAppendPlugin() *PluginInfo {
	return &PluginInfo{
		PluginName: "append",
		PluginType: CommandType,
		PluginAuthors: []string{
			"Aaron Delaney <devoxel@gmail.com>",
		},
		RegisterKeyword: "/append",
		Action: func(a *Action) error {
			// TODO: Implement this.
			return nil
		},
	}
}
