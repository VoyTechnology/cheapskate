// Package hook creates a new webserver ready to listen on
package hook

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cpssd-students/cheapskate/settings"
	"github.com/golang/glog"
)

// Errors
var (
	ErrServerAlreadyStarted = errors.New("server already started, can't register plugin")
)

var started = false

// Run the webhook
func Run() error {
	glog.Info("Starting web server")
	started = true
	return http.ListenAndServe(fmt.Sprintf(":%d", settings.Get("webhook.port").(int)), nil)
}

// RegisterHandleFunc allows to register the integration
func RegisterHandleFunc(path string, f http.Handler) error {
	if started {
		return ErrServerAlreadyStarted
	}

	glog.Infof("Registering %s", path)

	http.Handle(path, f)
	return nil
}
