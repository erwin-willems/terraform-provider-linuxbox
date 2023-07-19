package linuxbox

import (
	"fmt"
)

func validateServiceState(vi interface{}, k string) (ws []string, errors []error) {

	// Validate that the state is one of: started, stopped, restarted, reloaded:
	v, err := vi.(string)
	if !err {
		errors = append(errors, fmt.Errorf("state should be a string"))
	}
	if v != "started" && v != "stopped" && v != "restarted" && v != "reloaded" {
		errors = append(errors, fmt.Errorf("state should be one of: started, stopped, restarted, reloaded"))
	}
	return
}
