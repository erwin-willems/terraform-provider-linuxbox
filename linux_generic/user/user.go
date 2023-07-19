// Manage users	and groups
package linux_generic_user

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
	"github.com/pkg/errors"
)

type User struct {
	Name   string
	Uid    int
	Gid    int
	System bool
	Shell  string
	Home   string
}

type Group struct {
	Name   string
	Gid    int
	System bool
}


func CreateGroup(config *sshclient.Config, group *Group, sudo bool) error {
	command := "/usr/sbin/groupadd"

	if group.Gid > 0 {
		command = fmt.Sprintf("%s --gid %d", command, group.Gid)
	}
	if group.System {
		command = fmt.Sprintf("%s --system", command)
	}
	command = fmt.Sprintf("%s %s", command, group.Name)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}



func CreateUser(config *sshclient.Config, user *User, sudo bool) error {

	command := "/usr/sbin/useradd"

	if user.Uid > 0 {
		command = fmt.Sprintf("%s --uid %d", command, user.Uid)
	}
	if user.Gid > 0 {
		// If the user gid is specified, we first need to create the group otherwise useradd will fail
		group := Group{
			Name:   user.Name,
			Gid:    user.Gid,
			System: user.System,
		}

		err := CreateGroup(config, &group, sudo)
		if err != nil {
			return err
		}

		command = fmt.Sprintf("%s --gid %d", command, user.Gid)
	}
	if user.System {
		command = fmt.Sprintf("%s --system", command)
	}
	if user.Home != "" {
		command = fmt.Sprintf("%s --home %s", command, user.Home)
	}
	if user.Shell != "" {
		command = fmt.Sprintf("%s --shell %s", command, user.Shell)
	}

	command = fmt.Sprintf("%s %s", command, user.Name)
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func DeleteUser(config *sshclient.Config, name string, remove bool, sudo bool) error {
	command := fmt.Sprintf("userdel %s", name)

	if remove {
		command = fmt.Sprintf("%s --remove", command)
	}
	_, _, err := sshclient.Run(config, sudo, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func RenameUser(config *sshclient.Config, oldName string, newName string, sudo bool) error {
	command := fmt.Sprintf("usermod --login %s %s", newName, oldName)
	if _, _, err := sshclient.Run(config, sudo, command, ""); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func ChangeGroupId(config *sshclient.Config, groupname string, gid int, sudo bool) error {
	command := fmt.Sprintf("groupmod --gid %d %s", gid, groupname)
	if _, _, err := sshclient.Run(config, sudo, command, ""); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func ReadUser(config *sshclient.Config, user string) (*User, error) {

	command := fmt.Sprintf("getent passwd %s", user)
	stdout, _, err := sshclient.Run(config, false, command, "")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	if stdout == "" {
		return nil, fmt.Errorf("user not found with id %v", user)
	}
	stdout = strings.TrimSuffix(stdout, "\n")
	userParts := strings.Split(stdout, ":")
	name := userParts[0]
	uid, err := strconv.Atoi(userParts[2])
	if err != nil {
		return nil, err
	}
	gid, err := strconv.Atoi(userParts[3])
	if err != nil {
		return nil, err
	}
	home := userParts[5]
	shell := userParts[6]

	return &User{
		Name:   name,
		Uid:    uid, 
		Gid:    gid,
		System: false, //TODO fixme
		Shell:  shell,
		Home:   home,
	}, nil
}

func ReadGroup(config *sshclient.Config, group string) (*Group, error) {
	
	command := fmt.Sprintf("getent group %s", group)
	stdout, _, err := sshclient.Run(config, false, command, "")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	if stdout == "" {
		return nil, fmt.Errorf("group not found with id %v", group)
	}
	stdout = strings.TrimSuffix(stdout, "\n")
	groupParts := strings.Split(stdout, ":")
	name := groupParts[0]
	gid, err := strconv.Atoi(groupParts[2])
	if err != nil {
		return nil, err
	}

	return &Group{
		Name:   name,
		Gid:    gid,
		System: false, //TODO fixme
	}, nil
}

func ChangeHome(config *sshclient.Config, user string, home string, sudo bool) error {
	command := fmt.Sprintf("usermod --home %s %s", home, user)
	if _, _, err := sshclient.Run(config, sudo, command, ""); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func ChangeShell(config *sshclient.Config, user string, shell string, sudo bool) error {
	command := fmt.Sprintf("usermod --shell %s %s", shell, user)
	if _, _, err := sshclient.Run(config, sudo, command, ""); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}
