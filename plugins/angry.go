// Adapted from steely/plugins/angry.py by Wojtek Bednarzak

package plugins

type angryPlugin struct{}

func (angryPlugin) Type() Type {
	return RegexType
}

func (angryPlugin) Name() string {
	return "angry"
}

func (angryPlugin) Authors() []string {
	return []string{
		"Senan Kelly <senan@senan.xyz>",
	}
}

func (angryPlugin) Register() string {
	return "tayne|reed|bot|steely"
}

func (p *angryPlugin) Do(a *Action) error {

	// TODO: Update when Messenger plugin is available.
	// TODO: Add randomness.
	go AddAction(&Action{
		Origin:  p,
		Target:  a.Origin,
		Command: "response",
		Meta: map[string]string{
			"response_to": a.Meta["respond_to"],
			"emoji_size":  "large",
		},
		Data: []byte("😠"),
	})

	return nil
}

func (angryPlugin) NoAction() {}
