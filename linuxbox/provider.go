package linuxbox

import (
	"log"
	"os"

	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
	"github.com/erwin-willems/terraform-provider-linuxbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ssh_session_limit": {
				Type:     schema.TypeInt,
				Default:  5,
				Optional: true,
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_USER", ""),
				Description: "The username to ssh with",
			},
			"privilege_escalation": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_PRIV_ESC_METHOD", ""),
							Description: "Which method to use for Privileges Escalation? Valid values are sudo or \"\" to disallow Privilege Escalation",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_PRIV_ESC_PASSWORD", ""),
							Description: "The password to use for Privilege Escalation",
						},
					},
				},
			},
			"use_ssh_agent": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_USE_SSH_AGENT", true),
				Description: "Use the SSH Agent?",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_HOST", ""),
				Description: "The host to ssh into",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_PORT", 22),
				Description: "The ssh port",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_PASSWORD", ""),
				Description: "The password, if used for authentication",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_PRIVATE_KEY", "$HOME/.ssh/id_rsa"),
				Description: "The location of the private key, if used for authentication",
			},
			"private_key_passphrase": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_PRIVATE_KEY_PASSPHRASE", ""),
				Description: "The passphrase for the private key, if used for authentication",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"linuxbox_file":      datasourceLinuxboxFile(),
			"linuxbox_directory": datasourceLinuxboxDirectory(),
		},

		ResourcesMap: map[string]*schema.Resource{
			// "linuxbox_run_setup":          runsetup.Resource(),
			// "linuxbox_ssh_authorized_key": authorizedkey.Resource(),
			// "linuxbox_swap":               swap.Resource(),
			"linuxbox_file":      resourceLinuxboxFile(),
			"linuxbox_directory": resourceLinuxboxDirectory(),
			"linuxbox_user":      resourceLinuxboxUser(),
			"linuxbox_package":   resourceLinuxboxPackage(),
			"linuxbox_service":      resourceLinuxboxService(),
		},

		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			// Todo: Add support for ssh session limit
			// sshsession.SessionLimit = d.Get("ssh_session_limit").(int)

			// user :=

			var useSSHAgentBool bool
			useSSHAgent, ok := d.GetOk("use_ssh_agent")
			if !ok {
				useSSHAgentBool = false
			} else {
				useSSHAgentBool = useSSHAgent.(bool)
			}

			config := sshclient.Config{
				Host:             d.Get("host").(string),
				Port:             d.Get("port").(int),
				User:             d.Get("user").(string),
				Password:         d.Get("password").(string),
				PrivateKey:       os.ExpandEnv(d.Get("private_key").(string)),
				PrivateKeyPhrase: d.Get("private_key_passphrase").(string),
				PrivEscPassword:  d.Get("privilege_escalation.password").(string),
				PrivEscMethod:    d.Get("privilege_escalation.method").(string),
				UseSSHAgent:      useSSHAgentBool,
			}

			// workarounds for values that were not filled in correctly
			if config.PrivEscMethod == "" {
				config.PrivEscMethod = utils.EnvOrDefault("TF_LINUX_PRIV_ESC_METHOD", "")
			}
			if config.PrivEscPassword == "" {
				config.PrivEscPassword = utils.EnvOrDefault("TF_LINUX_PRIV_ESC_PASSWORD", "")
			}

			log.Println("Initializing SSH client")
			log.Printf(config.Host)
			return &config, nil
		},
	}
}
