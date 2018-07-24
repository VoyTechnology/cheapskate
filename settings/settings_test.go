package settings

import "testing"

func TestSet(t *testing.T) {
	key := "test"
	value := "expected"

	Set(key, value)
	if got := settings.Get(key).(string); got != value {
		t.Errorf("expected %s, got %s", value, got)
	}
}
