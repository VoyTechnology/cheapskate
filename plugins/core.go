package plugins

import (
	"context"
	"net/http"

	"github.com/golang/glog"
)

type corePlugin struct{}

func (corePlugin) Type() Type {
	return IntegrationType
}

func (corePlugin) Name() string {
	return "core"
}

func (corePlugin) Authors() []string {
	return []string{
		"Wojtek Bednarzak <wojtek.bednarzak.com>",
	}
}

func (corePlugin) Register() string {
	return "/core"
}

func (corePlugin) Do(ctx context.Context, a *Action) error {
	return nil
}

func (p corePlugin) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	data := []byte(req.URL.Query().Get("msg"))

	returnChan := make(chan []byte)

	a := &Action{
		Origin:   p,
		Data:     data,
		Response: returnChan,
	}

	AddAction(a)

	response := <-returnChan
	glog.Info("response: ", response)

	res.Write(response)
}
