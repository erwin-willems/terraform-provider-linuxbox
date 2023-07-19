package linuxbox

import (
	"fmt"
	"strings"
)

func validateDirectoryPath(vi interface{}, k string) (ws []string, errors []error) {
	v, err := vi.(string)
	if !err {
		errors = append(errors, fmt.Errorf("path should be a string"))
	}
	if !strings.HasPrefix(v, "/") {
		errors = append(errors, fmt.Errorf("path should be an absolute path"))
	}
	return
}

func validateDirectoryOwner(vi interface{}, k string) (ws []string, errors []error) {
	_, err := vi.(string)
	if !err {
		errors = append(errors, fmt.Errorf("owner should be a string"))
	}
	return
}

func validateDirectoryGroup(vi interface{}, k string) (ws []string, errors []error) {
	_, err := vi.(string)
	if !err {
		errors = append(errors, fmt.Errorf("group should be a string"))
	}
	return
}
