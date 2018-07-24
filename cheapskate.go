package main

import (
	"fmt"
	"net/http"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/cpssd-students/cheapskate/settings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func run() error {
	flag.Parse()

	http.HandleFunc("/", handleFunc)

	http.ListenAndServe(fmt.Sprintf(":%d", settings.Get("webhook.port").(int)), nil)

	return nil
}

func handleFunc(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("OK"))
}
