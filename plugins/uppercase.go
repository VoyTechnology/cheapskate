package plugins

import "strings"

func newUppercasePlugin() *PluginInfo {
	return &PluginInfo{
		PluginName: "uppercase",
		PluginType: CommandType,
		PluginAuthors: []string{
			"Wojtek Bednarzak <wojtek.bednarzak@gmail.com>",
		},
		RegisterKeyword: "/uppercase",
		Action: func(a *Action) error {
			a.Data = []byte(
				strings.ToUpper(
					strings.TrimPrefix(string(a.Data), "/uppercase ")))
			AddAction(a)
			return nil
		},
	}
}
