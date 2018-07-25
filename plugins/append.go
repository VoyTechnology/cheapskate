// Adapted from steely/plugins/append.py by Wojtek Bednarzak

package plugins

// Append some text
type appendPlugin struct{}

func (appendPlugin) Name() string {
	return "append"
}

func (appendPlugin) Type() Type {
	return CommandType
}

func (appendPlugin) Authors() []string {
	return []string{
		"Aaron Delaney <devoxel@gmail.com>",
	}
}

func (appendPlugin) Register() string {
	return "/append"
}

func (appendPlugin) Do(a *Action) error {
	// TODO: Get last message somehow..

	// TODO: Implement this.. Honestly I have no idea what this does.
	return nil
}

func (appendPlugin) NoAction() {}
