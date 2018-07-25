// Adapted from steely/plugins/clap.py by Wojtek Bednarzak

package plugins

import (
	"strings"
)

type clapPlugin struct{}

func (clapPlugin) Name() string {
	return "clap"
}

func (clapPlugin) Type() Type {
	return CommandType
}

func (clapPlugin) Authors() []string {
	return []string{
		"Senan Kelly <senan@senan.xyz>",
		"Aaron Delaney <devoxel@gmail.com>",
		"Wojtek Bednarzak <wojtek.bednarzak@gmail.com>",
	}
}

func (clapPlugin) Register() string {
	return "/clap"
}

func (p *clapPlugin) Do(a *Action) error {
	split := strings.Split(TrimPrefix(string(a.Data), "/clap"), " ")

	a.Data = []byte(strings.Join(split, " 👏 "))
	AddAction(a)
	return nil
}

func (clapPlugin) NoAction() {}
