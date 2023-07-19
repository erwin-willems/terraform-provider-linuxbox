// Manage yum packages
package linux_redhat_yum

import (
	"fmt"
	"strings"
	"regexp"
	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
	"github.com/pkg/errors"
)

func InstallList(config *sshclient.Config, packageNames []string, sudo bool) error {
	command := fmt.Sprintf("yum install --quiet --assumeyes %s", strings.Join(packageNames, " "))
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, "Command failed")
	}
	return nil
}

func Install(config *sshclient.Config, packageName string, version string, sudo bool) error {
	if version == "" {
		command := fmt.Sprintf("yum install --quiet --assumeyes %s", packageName)
		_, _, err := sshclient.Run(config, sudo, command, "")
		if err != nil {
			return errors.Wrap(err, "Command failed")
		}
	} else {
		command := fmt.Sprintf("yum install --quiet --assumeyes %s-%s", packageName, version)
		_, _, err := sshclient.Run(config, sudo, command, "")
		if err != nil {
			return errors.Wrap(err, "Command failed")
		}
	}
	return nil
}

func RemoveList(config *sshclient.Config, packageNames []string, sudo bool) error {
	command := fmt.Sprintf("yum remove --quiet --assumeyes %s", strings.Join(packageNames, " "))
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, "Command failed")
	}
	return nil
}

func Remove(config *sshclient.Config, packageName string, sudo bool) error {
	command := fmt.Sprintf("yum remove -y %s", packageName)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, "Command failed")
	}
	return nil
}

func Update(config *sshclient.Config, packageName string, version string, sudo bool) error {
	if version == "" {
		command := fmt.Sprintf("yum update --quiet --assumeyes %s", packageName)
		_, _, err := sshclient.Run(config, sudo, command, "")
		if err != nil {
			return errors.Wrap(err, "Command failed")
		}
	} else {
		command := fmt.Sprintf("yum update --quiet --assumeyes %s-%s", packageName, version)
		_, _, err := sshclient.Run(config, sudo, command, "")
		if err != nil {
			return errors.Wrap(err, "Command failed")
		}
	}
	return nil
}

func Check(config *sshclient.Config, packageName string, sudo bool) (bool, error) {
	command := fmt.Sprintf("yum list installed %s", packageName)
	stdout, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return false, errors.Wrap(err, "Command failed")
	}
	if stdout == "" {
		return false, nil
	}
	return true, nil
}

func Version(config *sshclient.Config, packageName string, sudo bool) (string, error) {
	command := fmt.Sprintf("yum info installed %s", packageName)
	stdout, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return "", errors.Wrap(err, "Command failed")
	}
	// extract version from output. Example output:
	// Installed Packages
	// Name         : bind-utils
	// Epoch        : 32
	// Version      : 9.16.23
	// Release      : 11.el9
	// Architecture : x86_64
	// Size         : 642 k
	// Source       : bind-9.16.23-11.el9.src.rpm
	// Repository   : @System
	// From repo    : appstream
	// Summary      : Utilities for querying DNS name servers
	// URL          : https://www.isc.org/downloads/bind/
	// License      : MPLv2.0
	// Description  : Bind-utils contains a collection of utilities for querying DNS (Domain
	// 			 : Name System) name servers to find out information about Internet
	// 			 : hosts. These tools will provide you with the IP addresses for given
	// 			 : host names, as well as other information about registered domains and
	// 			 : network addresses.
	// 			 : 
	// 			 : You should install bind-utils if you need to get information from DNS name
	// 			 : servers.
	// find the first line that starts with "Version"
	// then extract the version from that line
	// if the package is not installed, the output will be empty
	// so we can split on spaces and take the third element
	// if the package is installed, the output will be something like:
	// Error: No matching Packages to list
	// so we can split on spaces and take the fourth element
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Version") {
			words := regexp.MustCompile(`\s+`).Split(line, -1)
			if len(words) > 2 {
				return words[len(words)-1], nil
			}
		}
	}

	return "", nil
}
