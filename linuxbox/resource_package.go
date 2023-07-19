package linuxbox

import (
	"strings"
	"fmt"
	"github.com/erwin-willems/terraform-provider-linuxbox/internal/utils"
	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
	linux_redhat_yum "github.com/erwin-willems/terraform-provider-linuxbox/linux_redhat/yum"
	linux_debian_apt "github.com/erwin-willems/terraform-provider-linuxbox/linux_debian/apt"
	linux_generic_command "github.com/erwin-willems/terraform-provider-linuxbox/linux_generic/command"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

// todo:
// - Add support for dnf
// - Add support for zypper
// - Add validator for state, which can only be present or absent

// Multiple packages:
// locals {
// 	packages = [
// 	  "bzip2",
// 	  "bind-utils",
// 	  "git",
// 	]
// }
  
// resource "linuxbox_package" "packages" {
// 	for_each = toset(local.packages)
// 	name = each.key
// 	sudo = true
// }
  


func resourceLinuxboxPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourcePackageCreate,
		Read:   resourcePackageRead,
		Update: resourcePackageUpdate,
		Delete: resourcePackageDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the package",
			},
			"names": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The names of the packages",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The version of the package",
			},
			"installed_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The installed version of the package",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "present",
				ValidateFunc: validateServiceState,
				Description: "The state of the package. Can be present or absent.",
			},
			"sudo": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the command will be run with sudo.",
			},
		},
	}
}

func convertNames(names []interface{}) []string {
	var namesString []string
	for _, v := range names {
		namesString = append(namesString, v.(string))
	}
	return namesString
}

func resourcePackageCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*sshclient.Config)
	name := d.Get("name").(string)
	names := convertNames(d.Get("names").([]interface{}))
	version := d.Get("version").(string)
	state := d.Get("state").(string)
	sudo := d.Get("sudo").(bool)

	d.Set("version", version)

	switch state {
	case "present":
		if len(names) > 0 {
			// extract names from interface{}

			err := installPackages(config, names, sudo)
			if err != nil {
				return err
			}
			name = fmt.Sprintf("list_%s", utils.SetToMd5(names))
			d.Set("name", name)
			d.SetId(name)
			d.Set("installed_version", "")
		} else {
			version, err := installPackage(config, name, version, sudo)
			if err != nil {
				return err
			}
			d.SetId(fmt.Sprintf("%s:%s", name, version))
			d.Set("installed_version", version)
		}	
		return nil
	case "absent":
		err := removePackage(config, name, sudo)
		if err != nil {
			return err
		}
		d.Set("installed_version", "")
		d.SetId(name)
		return nil
	default:
		return errors.Errorf("Invalid state %s", state)
	}
}

func resourcePackageRead(d *schema.ResourceData, m interface{}) error {
	// config := m.(*sshclient.Config)
	state := d.Get("state").(string)
	id := d.Id()
	var name, version string
	// if d.Id() starts with 'list_' then it is a list of packages
	if strings.HasPrefix(id, "list_") {
		name = id
		version = ""
	} else {
		if state == "present" {
			name, version = parseId(d.Id())
		} else {
			name = d.Id()
		}
	}
	d.Set("name", name)
	d.Set("installed_version", version)
	return nil
}

func resourcePackageDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*sshclient.Config)
	state := d.Get("state").(string)
	sudo := d.Get("sudo").(bool)
	var name string
	var err error
	id := d.Id()
	if state == "present" {
		if strings.HasPrefix(id, "list_") {
			names := convertNames(d.Get("names").([]interface{}))

			err = removePackages(config, names, sudo)
		} else {
			name, _ = parseId(id)
			err = removePackage(config, name, sudo)
		}
	} else {
		name = id
		err = removePackage(config, name, sudo)
	}	
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourcePackageUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*sshclient.Config)
	name := d.Get("name").(string)
	version := d.Get("version").(string)
	state := d.Get("state").(string)
	sudo := d.Get("sudo").(bool)
	var oldName, oldVersion string

	if state == "present" {
		oldName, oldVersion = parseId(d.Id())
	} else {
		oldName = d.Id()
	}
	fmt.Printf("Version: %s\n", version)
	if name != oldName {
		err := removePackage(config, oldName, sudo)
		if err != nil {
			return err
		}
		newVersion, err := installPackage(config, name, version, sudo)		
		if err != nil {
			return err
		}
		d.Set("version", version)
		d.SetId(fmt.Sprintf("%s:%s", name, newVersion))
	}
	if version != oldVersion {
		newVersion, err := updatePackage(config, name, version, sudo)
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("%s:%s", name, newVersion))
	}
	return nil
}


// Extract name and version number from the id. Example: "bind-utils:9.16.23-11.el9" -> "bind-utils", "9.16.23-11.el9"
func parseId(id string) (string, string) {
	parts := strings.Split(id, ":")
	return parts[0], parts[1]
}

func removePackages(config *sshclient.Config, names []string, sudo bool) error {
	packageManager, err := detectPackageManager(config)
	if err != nil {
		return err
	}
	switch packageManager {
	case "yum":
		err := linux_redhat_yum.RemoveList(config, names, sudo)
		if err != nil {
			return err
		}
		return nil
	case "apt":
		err := linux_debian_apt.RemoveList(config, names, sudo)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("No package manager found")
	}
}

func removePackage(config *sshclient.Config, name string, sudo bool) error {
	packageManager, err := detectPackageManager(config)
	if err != nil {
		return err
	}
	switch packageManager {
	case "yum":
		err := linux_redhat_yum.Remove(config, name, sudo)
		if err != nil {
			return err
		}
		return nil
	case "apt":
		err := linux_debian_apt.Remove(config, name, sudo)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("No package manager found")
	}
}


// Install package and return the version that was installed or the error if any
func installPackage(config *sshclient.Config, name string, version string, sudo bool) (string, error) {
	packageManager, err := detectPackageManager(config)
	if err != nil {
		return "", err
	}
	switch packageManager {
	case "yum":
		err := linux_redhat_yum.Install(config, name, version, sudo)
		if err != nil {
			return "", err
		}
		version, err := linux_redhat_yum.Version(config, name, sudo)
		if err != nil {
			return "", err
		}
		return version, nil
	case "apt":
		err := linux_debian_apt.Install(config, name, version, sudo)
		if err != nil {
			return "", err
		}
		version, err := linux_debian_apt.Version(config, name, sudo)
		if err != nil {
			return "", err
		}
		return version, nil
	default:
		return "", errors.New("No package manager found")
	}
}

func installPackages(config *sshclient.Config, names []string, sudo bool) error {
	packageManager, err := detectPackageManager(config)
	if err != nil {
		return err
	}
	switch packageManager {
	case "yum":
		err := linux_redhat_yum.InstallList(config, names, sudo)
		if err != nil {
			return err
		}
		return nil
	case "apt":
		err := linux_debian_apt.InstallList(config, names, sudo)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("No package manager found")
	}
}


func updatePackage(config *sshclient.Config, name string, version string, sudo bool) (string, error) {
	packageManager, err := detectPackageManager(config)
	if err != nil {
		return "", err
	}
	switch packageManager {
	case "yum":
		err := linux_redhat_yum.Update(config, name, version, sudo)
		if err != nil {
			return "", err
		}
		version, err := linux_redhat_yum.Version(config, name, sudo)
		if err != nil {
			return "", err
		}
		return version, nil
	case "apt":
		err := linux_debian_apt.Update(config, name, version, sudo)
		if err != nil {
			return "", err
		}
		version, err := linux_debian_apt.Version(config, name, sudo)
		if err != nil {
			return "", err
		}
		return version, nil
	default:
		return "", errors.New("No package manager found")
	}
}


// Detect which package manager is available on the system
func detectPackageManager(config *sshclient.Config) (string, error) {
	// Check if yum is available
	yumExists, err := linux_generic_command.CommandExists(config, "yum", false)
	if err != nil {
		return "", err
	}
	if yumExists {
		return "yum", nil
	}

	// Check if apt-get is available
	aptGetExists, err := linux_generic_command.CommandExists(config, "apt", false)
	if err != nil {
		return "", err
	}
	if aptGetExists {
		return "apt", nil
	}

	return "", errors.New("No package manager found")
}

	
