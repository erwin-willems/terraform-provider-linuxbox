package linuxbox

import (
	"testing"
)


func TestInvalidFilePath(t *testing.T) {
	_, err := validateFilePath("some/path", "")
	if err == nil {
		t.Errorf("Non-absolute path should be invalid: %v", err)
	}
}

func TestValidFilePath(t *testing.T) {
	_, err := validateFilePath("/some/path", "")
	if err != nil {
		t.Errorf("Absolute path should valid: %v", err)
	}
}

func TestValidFileOwner(t *testing.T) {
	_, err := validateFileOwner("root", "")
	if err != nil {
		t.Errorf("Owner should valid: %v", err)
	}
}

func TestValidFileGroup(t *testing.T) {
	_, err := validateFileGroup("root", "")
	if err != nil {
		t.Errorf("Group should valid: %v", err)
	}
}

func TestInvalidFileOwner(t *testing.T) {
	_, err := validateFileOwner(123, "")
	if err == nil {
		t.Errorf("Owner should be a string: %v", err)
	}
}

func TestInvalidFileGroup(t *testing.T) {
	_, err := validateFileGroup(123, "")
	if err == nil {
		t.Errorf("Group should be a string: %v", err)
	}
}
