package linux_generic_systemd

import (
	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
)

func Start(config *sshclient.Config, serviceName string, sudo bool) error {
	command := "systemctl start " + serviceName
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return err
	}
	return nil
}

func Stop(config *sshclient.Config, serviceName string, sudo bool) error {
	command := "systemctl stop " + serviceName
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return err
	}
	return nil
}

func Restart(config *sshclient.Config, serviceName string, sudo bool) error {
	command := "systemctl restart " + serviceName
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return err
	}
	return nil
}

func Reload(config *sshclient.Config, serviceName string, sudo bool) error {
	command := "systemctl reload " + serviceName
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return err
	}
	return nil
}

func Enable(config *sshclient.Config, serviceName string, sudo bool) error {
	command := "systemctl enable " + serviceName
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return err
	}
	return nil
}

func Disable(config *sshclient.Config, serviceName string, sudo bool) error {
	command := "systemctl disable " + serviceName
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return err
	}
	return nil
}

func IsEnabled(config *sshclient.Config, serviceName string, sudo bool) (bool, error) {
	command := "systemctl is-enabled " + serviceName
	_, stdout, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return false, err
	}
	if stdout == "enabled" {
		return true, nil
	}
	return false, nil
}

func IsRunning(config *sshclient.Config, serviceName string, sudo bool) (bool, error) {
	command := "systemctl is-active " + serviceName
	_, stdout, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return false, err
	}
	if stdout == "active" {
		return true, nil
	}
	return false, nil
}

func Status(config *sshclient.Config, serviceName string, sudo bool) (string, error) {
	command := "systemctl status " + serviceName
	_, stdout, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return "", err
	}
	return stdout, nil
}

func DaemonReload(config *sshclient.Config, sudo bool) error {
	command := "systemctl daemon-reload"
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return err
	}
	return nil
}


