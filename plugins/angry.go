// Adapted from steely/plugins/angry.py by Wojtek Bednarzak

package plugins

func newAngryPlugin() *PluginInfo {
	return &PluginInfo{
		PluginName: "angry",
		PluginType: RegexType,
		PluginAuthors: []string{
			"Senan Kelly <senan@senan.xyz>",
		},
		RegisterKeyword: "tayne|reed|bot|steely",
		Action: func(a *Action) error {
			// TODO: Update when Messenger plugin is available.
			// TODO: Add randomness.
			go AddAction(&Action{
				Origin:  a.Origin,
				Command: "response",
				Meta: map[string]string{
					"response_to": a.Meta["respond_to"],
					"emoji_size":  "large",
				},
				Data: []byte("😠"),
			})
			return nil
		},
	}
}
