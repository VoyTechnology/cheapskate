package config

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDefaults(t *testing.T) {
	c, err := Get()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	want := 8080
	if got := c.Server.Port; got != want {
		t.Fatalf("c.Server.Port = %v, want %v", got, want)
	}

}

func TestEnvVar(t *testing.T) {
	want := 1111
	os.Setenv(strings.ToUpper(prefix)+"_SERVER_PORT", fmt.Sprint(want))

	c, err := Get()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if got := c.Server.Port; got != want {
		t.Fatalf("c.Server.Port = %v, want %v", got, want)
	}
}
