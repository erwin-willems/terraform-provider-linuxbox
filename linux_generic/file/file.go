// Package linux_generic_file provides a generic file and directory operations implementation for Linux.
package linux_generic_file

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
	"github.com/pkg/errors"
)

func CreateFile(config *sshclient.Config, path string, sudo bool) error {
	command := fmt.Sprintf("touch %s", path)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func CreateDirectory(config *sshclient.Config, path string, sudo bool) error {
	command := fmt.Sprintf("mkdir -p %s", path)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

// Implements the chown command
func ChangeOwner(config *sshclient.Config, path string, owner string, group string, sudo bool) error {
	if owner == "" && group == "" {
		return nil
	}
	var command string
	if owner != "" && group == "" {
		command = fmt.Sprintf("chown %s %s", owner, path)
	}
	if owner == "" && group != "" {
		command = fmt.Sprintf("chgrp %s %s", group, path)
	}
	if owner != "" && group != "" {
		command = fmt.Sprintf("chown %s:%s %s", owner, group, path)
	}

	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

// Implements the chgrp command
func ChangeGroup(config *sshclient.Config, path string, group string, sudo bool) error {
	if group == "" {
		return errors.New("Group not specified")
	}
	command := fmt.Sprintf("chgrp %s %s", group, path)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

// Implements the chmod command
func ChangeMode(config *sshclient.Config, path string, mode int, sudo bool) error {
	command := fmt.Sprintf("chmod %d %s", mode, path)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func WriteContent(config *sshclient.Config, path string, content string, sudo bool) error {
	b64Content := base64.StdEncoding.EncodeToString([]byte(content))
	command := fmt.Sprintf("echo '%s' | base64 -d > %s", b64Content, path)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func parsePermissionString(perms string) int {
	permMap := map[string]int{
		"---": 0,
		"--x": 1,
		"-w-": 2,
		"-wx": 3,
		"r--": 4,
		"r-x": 5,
		"rw-": 6,
		"rwx": 7,
	}
	return (permMap[perms[1:4]] * 100) +
		(permMap[perms[4:7]] * 10) +
		(permMap[perms[7:10]] * 1)
}

func GetDetails(config *sshclient.Config, path string, sudo bool) (string, string, int, error) {
	command := fmt.Sprintf("ls -ld %s", path)
	stdout, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return "", "", 0, errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	if stdout == "" {
		return "", "", 0, fmt.Errorf("File not found with path %v", path)
	}
	fields := strings.Fields(stdout)
	mode, owner, group := parsePermissionString(fields[0]), fields[2], fields[3]
	if err != nil {
		return "", "", 0, errors.Wrap(err, fmt.Sprintf("Unable to parse permission string from %s", command))
	}
	return owner, group, mode, nil
}

func ReadFile(config *sshclient.Config, path string, sudo bool) (string, error) {
	command := fmt.Sprintf("cat %s", path)
	stdout, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return stdout, nil
}

// Implements the mv command.
// Works for both files and directories
func Move(config *sshclient.Config, oldPath string, newPath string, sudo bool) error {
	command := fmt.Sprintf("mv %s %s", oldPath, newPath)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

// Implements the rm command.
// Works for both files and directories
func Remove(config *sshclient.Config, path string, sudo bool) error {
	command := fmt.Sprintf("rm -rf %s", path)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}
