package plugins

import (
	"net/http"
)

// testerIntegrationPlugin is designed to be used to manually trigger certain actions.
// it should be disabled when reploying to production.
type testerIntegrationPlugin struct {
	msg chan []byte
}

func (*testerIntegrationPlugin) Type() Type {
	return IntegrationType
}

func (*testerIntegrationPlugin) Name() string {
	return "tester_integration"
}

func (*testerIntegrationPlugin) Authors() []string {
	return []string{
		"Wojtek Bednarzak <wojtek.bednarzak.com>",
	}
}

func (*testerIntegrationPlugin) Register() string {
	return "/tester"
}

// Do implements the function. The accepted commands are as follows:
//	response - if the integration is expecting back a response
func (p *testerIntegrationPlugin) Do(a *Action) error {
	p.msg <- a.Data

	return nil
}

func (p *testerIntegrationPlugin) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	p.msg = make(chan []byte)
	data := []byte(req.URL.Query().Get("msg"))
	if len(data) == 0 {
		res.Write([]byte("empty message"))
		return
	}

	a := &Action{
		Origin: p,
		Data:   data,
	}

	go AddAction(a)

	// We only respond once we get a response back from the system
	res.Write(<-p.msg)
}

func (p *testerIntegrationPlugin) NoAction() {
	p.msg <- []byte("")
}

func newTesterCommandPlugin() *PluginInfo {
	return &PluginInfo{
		PluginName: "tester_command",
		PluginType: CommandType,
		PluginAuthors: []string{
			"Wojtek Bednarzak <wojtek.bednarzak@gmail.com",
		},
		RegisterKeyword: "/tester",
		Action: func(a *Action) error {
			a.Data = []byte("got message with data " + TrimPrefix(
				string(a.Data), "/tester"))
			AddAction(a)
			return nil
		},
	}
}
