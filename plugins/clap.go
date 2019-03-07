// Adapted from steely/plugins/clap.py by Wojtek Bednarzak

package plugins

import (
	"strings"
)

func newClapPlugin() *PluginInfo {
	return &PluginInfo{
		PluginName: "clap",
		PluginType: CommandType,
		PluginAuthors: []string{
			"Senan Kelly <senan@senan.xyz>",
			"Aaron Delaney <devoxel@gmail.com>",
			"Wojtek Bednarzak <wojtek.bednarzak@gmail.com>",
		},
		RegisterKeyword: "/clap",
		Action: func(a *Action) error {
			split := strings.Fields(TrimPrefix(string(a.Data), "/clap"))

			a.Data = []byte(strings.Join(split, " 👏 "))
			AddAction(a)
			return nil
		},
	}
}
