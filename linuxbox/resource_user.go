package linuxbox

import (
	"fmt"
	"strings"

	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
	linux_generic_user "github.com/erwin-willems/terraform-provider-linuxbox/linux_generic/user"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

func resourceLinuxboxUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			// todo:
			// - add password support

			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "The name of the user. This is the name that will be used to login with. It is also the name of the home directory of the user.",
			},
			"uid": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: "The user id of the user. If not specified, the user will be created with a user id equal to the group id.",
			},
			"gid": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "The group id of the user. If not specified, the user will be created with a group id equal to the user id.",
			},
			//TODO: Check what system does..
			"system": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "If true, the user will be created as a system user. System users are typically used for running services. They have no shell and are not allowed to login.",
			},
			"shell": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "The shell of the user. Defaults to /bin/bash. Use /bin/false to disable shell access.",
			},
			"home": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "The home directory of the user. Add $username to use the username in the path.",
			},
			"remove_home": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "If true, the home directory of the user will be removed when the user is deleted. If false, the home directory will be left untouched.",
			},
			"sudo": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "If true, use sudo to create the user. Must be true if the user specified in the provider is not root.",
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*sshclient.Config)
	sudo := d.Get("sudo").(bool)
	user := linux_generic_user.User{
		Name:   d.Get("name").(string),
		Uid:    d.Get("uid").(int),
		Gid:    d.Get("gid").(int),
		System: d.Get("system").(bool),
		Home:   d.Get("home").(string),
		Shell:  d.Get("shell").(string),
	}

	// replace $username in home path with username
	user.Home = strings.Replace(user.Home, "$username", user.Name, -1)

	// set default shell
	if user.Shell == "" {
		user.Shell = "/bin/bash"
	}

	err := linux_generic_user.CreateUser(config, &user, sudo)
	if err != nil {
		return errors.Wrap(err, "Couldn't create user")
	}

	userCreated, err := linux_generic_user.ReadUser(config, user.Name)
	if err != nil {
		return errors.Wrap(err, "Couldn't read user")
	}
	d.Set("uid", userCreated.Uid)
	d.Set("gid", userCreated.Gid)
	d.Set("system", userCreated.System)
	d.Set("home", userCreated.Home)
	d.Set("shell", userCreated.Shell)

	d.SetId(fmt.Sprintf("%v", userCreated.Uid))
	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*sshclient.Config)
	
	user, err := linux_generic_user.ReadUser(config, d.Id())
	if err != nil {
		return errors.Wrap(err, "Couldn't read user")
	}
	d.Set("name", user.Name)
	d.Set("uid", user.Uid)
	d.Set("gid", user.Gid)
	d.Set("system", user.System)
	d.Set("home", user.Home)
	d.Set("shell", user.Shell)
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*sshclient.Config)

	user := linux_generic_user.User{
		Name:   d.Get("name").(string),
		Uid:    d.Get("uid").(int),
		Gid:    d.Get("gid").(int),
		System: d.Get("system").(bool),
		Home:   d.Get("home").(string),
		Shell:  d.Get("shell").(string),
	}
	sudo := d.Get("sudo").(bool)

	// replace $username in home path with username
	user.Home = strings.Replace(user.Home, "$username", user.Name, -1)

	oldUser, err := linux_generic_user.ReadUser(config, d.Id())
	if err != nil {
		return errors.Wrap(err, "Couldn't read user")
	}

	if oldUser.Name != user.Name {
		if err := linux_generic_user.RenameUser(config, oldUser.Name, user.Name, sudo); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to change the user name for user %s", oldUser.Name))
		}
	}

	if oldUser.Gid != user.Gid {
		// Changing the gid requires changing the group id of the user's group
		oldGroup, err := linux_generic_user.ReadGroup(config, fmt.Sprintf("%d", oldUser.Gid))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read group %d", oldUser.Gid))
		}
		if err := linux_generic_user.ChangeGroupId(config, oldGroup.Name, user.Gid, sudo); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to change the group id for user %s", user.Name))
		}
	}

	if oldUser.Home != user.Home {
		if err := linux_generic_user.ChangeHome(config, user.Name, user.Home, sudo); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to change the home directory for user %s", user.Name))
		}
	}

	if oldUser.Shell != user.Shell {
		if err := linux_generic_user.ChangeShell(config, user.Name, user.Shell, sudo); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to change the shell for user %s", user.Name))
		}
	}
	

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*sshclient.Config)
	remove := d.Get("remove_home").(bool)
	sudo := d.Get("sudo").(bool)
	user, err := linux_generic_user.ReadUser(config, d.Id())
	if err != nil {
		return errors.Wrap(err, "Failed to get user name")
	}

	if err := linux_generic_user.DeleteUser(config, user.Name, remove, sudo); err != nil {
		return errors.Wrap(err, "Couldn't delete user")
	}

	return nil
}
