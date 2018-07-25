package plugins

import "strings"

// uppercasePlugin brings all data to uppercase state
type uppercasePlugin struct{}

func (uppercasePlugin) Name() string {
	return "uppercase"
}

func (uppercasePlugin) Type() Type {
	return CommandType
}

func (uppercasePlugin) Register() string {
	return "/uppercase"
}

func (uppercasePlugin) Authors() []string {
	return []string{
		"Wojtek Bednarzak <wojtek.bednarzak@gmail.com>",
	}
}

func (uppercasePlugin) Do(a *Action) error {
	a.Data = []byte(
		strings.ToUpper(
			strings.TrimPrefix(string(a.Data), "/uppercase ")))
	AddAction(a)
	return nil
}

func (uppercasePlugin) NoAction() {}
