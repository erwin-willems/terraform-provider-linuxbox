package linuxbox

import (
	"testing"
)

func TestInvalidPackageState(t *testing.T) {
	_, err := validatePackageState("notvalid", "")
	if err == nil {
		t.Errorf("The state notvalid should be invalid: %v", err)
	}
}

func TestValidPackageState(t *testing.T) {
	for _, state := range []string{"present", "absent"} {
		_, err := validatePackageState(state, "")
		if err != nil {
			t.Errorf("The state %s should valid: %v", state, err)
		}
	}
}
