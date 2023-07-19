package linuxbox

import (
	"testing"
)


func TestInvalidDirectoryPath(t *testing.T) {
	_, err := validateDirectoryPath("some/path", "")
	if err == nil {
		t.Errorf("Non-absolute path should be invalid: %v", err)
	}
}

func TestValidDirectoryPath(t *testing.T) {
	_, err := validateDirectoryPath("/some/path", "")
	if err != nil {
		t.Errorf("Absolute path should valid: %v", err)
	}
}
