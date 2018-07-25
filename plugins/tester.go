package plugins

import (
	"fmt"
	"net/http"
	"strings"
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
	if a.Command == "response" {
		p.msg <- a.Data
	}

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

// Registeres a tester
type testerCommandPlugin struct{}

func (*testerCommandPlugin) Type() Type {
	return CommandType
}

func (*testerCommandPlugin) Name() string {
	return "tester_command"
}

func (*testerCommandPlugin) Authors() []string {
	return []string{
		"Wojtek Bednarzak <wojtek.bednarzak@gmail.com",
	}
}

func (*testerCommandPlugin) Register() string {
	return "/tester"
}

func (p *testerCommandPlugin) Do(a *Action) error {
	go AddAction(&Action{
		Origin:  p,
		Target:  a.Origin,
		Command: "response",
		Data: []byte(fmt.Sprintf(
			"got message with data %s",
			strings.TrimPrefix(string(a.Data), "/tester "))),
	})
	return nil
}

func (testerCommandPlugin) NoAction() {}
