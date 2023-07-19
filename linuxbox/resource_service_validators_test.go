package linuxbox

import (
	"testing"
)

func TestInvalidServiceState(t *testing.T) {
	_, err := validateServiceState("notvalid", "")
	if err == nil {
		t.Errorf("The state notvalid should be invalid: %v", err)
	}
}

func TestValidServiceState(t *testing.T) {
	for _, state := range []string{"started", "stopped", "reloaded", "restarted"} {
		_, err := validateServiceState(state, "")
		if err != nil {
			t.Errorf("The state %s should valid: %v", state, err)
		}
	}
}
