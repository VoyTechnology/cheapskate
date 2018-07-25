package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cpssd-students/cheapskate/hook"
	"github.com/cpssd-students/cheapskate/plugins"

	_ "github.com/cpssd-students/cheapskate/settings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func run() error {
	if err := plugins.Load(); err != nil {
		return err
	}

	if err := hook.Run(); err != nil {
		return err
	}

	return nil
}

func handleFunc(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("OK"))
}
