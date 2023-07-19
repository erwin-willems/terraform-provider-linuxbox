package linuxbox

import (
	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	linux_generic_command "github.com/erwin-willems/terraform-provider-linuxbox/linux_generic/command"
	linux_generic_systemd "github.com/erwin-willems/terraform-provider-linuxbox/linux_generic/systemd"
	"github.com/pkg/errors"

)



func resourceLinuxboxService() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceCreate,
		Read:   resourceServiceRead,
		Update: resourceServiceUpdate,
		Delete: resourceServiceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the package",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:    true,
				Description: "Enable the service",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "started",
				ValidateFunc: validateServiceState,
				Description: "The state of the package. Can be started, stopped, restarted or reloaded.",
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

func resourceServiceCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*sshclient.Config)
	name := d.Get("name").(string)
	enabled := d.Get("enabled").(bool)
	state := d.Get("state").(string)
	sudo := d.Get("sudo").(bool)
	
	service, err := detectServiceDeamon(config)
	if err != nil {
		return err
	}

	if enabled {
		err := enableService(config, service, name, sudo)
		if err != nil {
			return err
		}
	} else {
		err := disableService(config, service, name, sudo)
		if err != nil {
			return err
		}
	}

	switch state {
	case "started":
		err = startService(config, service, name, sudo)
	case "stopped":
		err = stopService(config, service, name, sudo)
	case "restarted":
		err = restartService(config, service, name, sudo)
	case "reloaded":
		err = reloadService(config, service, name, sudo)
	default: 
		return errors.Errorf("Service state %s not supported", state)
	}
	if err != nil {
		return err
	}

	return nil
}

func resourceServiceRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}



func detectServiceDeamon(config *sshclient.Config) (string, error) {
	// Check if systemd is available
	systemdExists, err := linux_generic_command.CommandExists(config, "systemctl", false)
	if err != nil {
		return "", err
	}
	if systemdExists {
		return "systemd", nil
	}

	return "", errors.New("No service deamon found")
}


func enableService(config *sshclient.Config, serviceDeamon string, name string, sudo bool) error {
	switch serviceDeamon {
	case "systemd":
		return linux_generic_systemd.Enable(config, name, sudo)
	default:
		return errors.New("Service deamon not supported")
	}
}

func disableService(config *sshclient.Config, serviceDeamon string, name string, sudo bool) error {
	switch serviceDeamon {
	case "systemd":
		return linux_generic_systemd.Disable(config, name, sudo)
	default:
		return errors.New("Service deamon not supported")
	}
}

func startService(config *sshclient.Config, serviceDeamon string, name string, sudo bool) error {
	switch serviceDeamon {
	case "systemd":
		return linux_generic_systemd.Start(config, name, sudo)
	default:
		return errors.New("Service deamon not supported")
	}
}

func stopService(config *sshclient.Config, serviceDeamon string, name string, sudo bool) error {
	switch serviceDeamon {
	case "systemd":
		return linux_generic_systemd.Stop(config, name, sudo)
	default:
		return errors.New("Service deamon not supported")
	}
}

func restartService(config *sshclient.Config, serviceDeamon string, name string, sudo bool) error {
	switch serviceDeamon {
	case "systemd":
		return linux_generic_systemd.Restart(config, name, sudo)
	default:
		return errors.New("Service deamon not supported")
	}
}

func reloadService(config *sshclient.Config, serviceDeamon string, name string, sudo bool) error {
	switch serviceDeamon {
	case "systemd":
		return linux_generic_systemd.Reload(config, name, sudo)
	default:
		return errors.New("Service deamon not supported")
	}
}
