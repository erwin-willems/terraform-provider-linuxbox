// Manage yum packages
package linux_debian_apt

import (
	"fmt"
	"strings"
	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
	"github.com/pkg/errors"
)

func aptUpdate(config *sshclient.Config, sudo bool) error {
	command := "apt update -y"
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, "Command failed")
	}
	return nil
}

func InstallList(config *sshclient.Config, packageNames []string, sudo bool) error {
	command := fmt.Sprintf("apt install -y %s", strings.Join(packageNames, " "))
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, "Command failed")
	}
	return nil
}

func Install(config *sshclient.Config, packageName string, version string, sudo bool) error {
	err := aptUpdate(config, sudo)
	if err != nil {
		return errors.Wrap(err, "Command failed")
	}

	if version == "" {
			
		command := fmt.Sprintf("apt install -y %s", packageName)
		_, _, err = sshclient.Run(config, sudo, command, "")
		if err != nil {
			return errors.Wrap(err, "Command failed")
		}
	} else {
		command := fmt.Sprintf("apt install -y %s=%s", packageName, version)
		_, _, err = sshclient.Run(config, sudo, command, "")
		if err != nil {
			return errors.Wrap(err, "Command failed")
		}
	}
	return nil
}

func RemoveList(config *sshclient.Config, packageNames []string, sudo bool) error {
	command := fmt.Sprintf("apt remove -y %s", strings.Join(packageNames, " "))
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, "Command failed")
	}
	return nil
}

func Remove(config *sshclient.Config, packageName string, sudo bool) error {
	command := fmt.Sprintf("apt remove -y %s", packageName)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, "Command failed")
	}
	return nil
}

func Update(config *sshclient.Config, packageName string, version string, sudo bool) error {
	err := aptUpdate(config, sudo)
	if err != nil {
		return errors.Wrap(err, "Command failed")
	}
	if version == "" {

		command := fmt.Sprintf("apt upgrade -y %s", packageName)
		_, _, err = sshclient.Run(config, sudo, command, "")
		if err != nil {
			return errors.Wrap(err, "Command failed")
		}
	} else {
		command := fmt.Sprintf("apt upgrade -y %s=%s", packageName, version)
		_, _, err = sshclient.Run(config, sudo, command, "")
		if err != nil {
			return errors.Wrap(err, "Command failed")
		}
	}
	return nil
}

func Check(config *sshclient.Config, packageName string, sudo bool) (bool, error) {
	command := fmt.Sprintf("apt list installed %s", packageName)
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
	command := fmt.Sprintf("apt list installed %s", packageName)
	stdout, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return "", errors.Wrap(err, "Command failed")
	}
	// extract version from output. Example output:
	// bind9/now 1:9.11.5.P4+dfsg-5.1ubuntu2.1 amd64 [installed,upgradable to: 1:9.11.5.P4+dfsg-5.1ubuntu2.2]
	return stdout, nil
}