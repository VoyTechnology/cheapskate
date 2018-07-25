// Adapted from steely/plugins/clap.py by Wojtek Bednarzak

package plugins

import "strings"

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
	split := strings.Split(strings.TrimPrefix(string(a.Data), "/clap "), " ")

	withClap := strings.Join(split, " 👏 ")

	go AddAction(&Action{
		Origin:  p,
		Target:  a.Origin,
		Command: "response",
		Data:    []byte(withClap),
	})
	return nil
}

func (clapPlugin) NoAction() {}
