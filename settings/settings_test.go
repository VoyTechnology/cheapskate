package settings

import (
	"fmt"
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	expected := 1000
	os.Setenv("CS_WEBHOOK_PORT", fmt.Sprint(expected))

	s, err := New()
	if err != nil {
		t.Fatalf("Init().Error = %v, want %v", err, nil)
	}

	if port := s.Webhook.Port; port != expected {
		t.Errorf("Settings.Webhook.Port = %d, want %d", port, expected)
	}
}
