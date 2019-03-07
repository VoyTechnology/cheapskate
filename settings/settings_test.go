package settings

import (
	"fmt"
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	// TODO: THIS IS WRONG. This is the default and that is the only reason
	//       why it works. ENV VARS do not work.
	expected := 8080
	os.Setenv("CS_WEBHOOK_PORT", fmt.Sprint(expected))

	s, err := New()
	if err != nil {
		t.Fatalf("Init().Error = %v, want %v", err, nil)
	}

	if port := s.Webhook.Port; port != expected {
		t.Errorf("Settings.Webhook.Port = %d, want %d", port, expected)
	}
}
