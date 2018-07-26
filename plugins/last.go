package plugins

func newLastPlugin() *PluginInfo {
	return &PluginInfo{
		PluginName: "last",
		PluginType: CommandType,
		PluginAuthors: []string{
			"Wojtek Bednarzak <wojtek.bednarzak@gmail.com>",
		},
		RegisterKeyword: "/last",
		Action: func(a *Action) error {
			d := TrimPrefix(string(a.Data), "/last")

			if lastAction == nil {
				a.Data = []byte(d)
				a.Done = true
				AddAction(a)
				return nil
			}

			if len(d) == 0 {
				a.Data = lastAction.Data
				AddAction(a)
				return nil
			}

			a.Data = []byte(d + " " + string(lastAction.Data))
			AddAction(a)
			return nil
		},
	}
}
