package linuxbox

import (
	"strings"

	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
	linux_generic_file "github.com/erwin-willems/terraform-provider-linuxbox/linux_generic/file"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

func datasourceLinuxboxFile() *schema.Resource {
	return &schema.Resource{
		Read: datasourceFileRead(),

		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateFilePath,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sudo": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func datasourceFileRead() func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		config := m.(*sshclient.Config)
		path := d.Get("path").(string)
		sudo := d.Get("sudo").(bool)

		owner, group, mode, err := linux_generic_file.GetDetails(config, path, sudo)
		if err != nil {
			if strings.Contains(err.Error(), "File not found with path") {
				d.SetId("")
				return nil
			}
			return errors.Wrap(err, "Unable to ls the file")
		}

		content, err := linux_generic_file.ReadFile(config, path, sudo)
		if err != nil {
			return errors.Wrap(err, "Unable to read the file")
		}
		d.SetId(path)
		d.Set("content", content)

		d.Set("owner", owner)
		d.Set("group", group)
		d.Set("mode", mode)
		return nil
	}
}