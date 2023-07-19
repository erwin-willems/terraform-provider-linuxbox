package linuxbox

import (
	"fmt"
)

func validatePackageState(vi interface{}, k string) (ws []string, errors []error) {

	// Validate that the state is one of: present, absent:
	v, err := vi.(string)
	if !err {
		errors = append(errors, fmt.Errorf("state should be a string"))
	}
	if v != "present" && v != "absent" {
		errors = append(errors, fmt.Errorf("state should be one of: present, absent"))
	}
	return
}