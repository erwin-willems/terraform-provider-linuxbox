package linuxbox

import (
	"strings"
	"time"

	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
	linux_generic_file "github.com/erwin-willems/terraform-provider-linuxbox/linux_generic/file"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

func resourceLinuxboxFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceFileCreate(),
		Read:   resourceFileRead(),
		Update: resourceFileUpdate(),
		Delete: resourceFileDelete,

		// Todo:
		// - Add host, to allow other hosts to be used than specified in the provider config
		// - Add content_base64, to allow content to be specified in base64
		// - Add validation for content_base64, to check if it is valid base64
		// - Add validation for content_base64 and content should be mutually exclusive
		// - Add backup, to allow backup of the file before it is changed
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateFilePath,
			},
			"owner": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateFileOwner,
			},
			"group": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateFileGroup,
			},
			"mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"sudo": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Timeouts: &schema.ResourceTimeout{ // TODO: Set this to reasonable values
			Create: schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
	}
}

func rollbackFile(config *sshclient.Config, err error, errMsg string, path string, sudo bool) error {
	err2 := errors.Wrap(err, errMsg)
	if err3 := linux_generic_file.Remove(config, path, sudo); err3 != nil {
		err3 = errors.Wrap(err2, err3.Error())
		return errors.Wrap(err3, "Couldn't delete file.")
	}
	return err2
}

func resourceFileCreate() func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		config := m.(*sshclient.Config)
		path := d.Get("path").(string)
		owner := d.Get("owner").(string)
		group := d.Get("group").(string)
		mode := d.Get("mode").(int)
		sudo := d.Get("sudo").(bool)

		if err := linux_generic_file.CreateFile(config, path, sudo); err != nil {
			return errors.Wrap(err, "Couldn't create file")
		}

		if owner != "" || group != "" {
			if err := linux_generic_file.ChangeOwner(config, path, owner, group, sudo); err != nil {
				return rollbackFile(config, err, "Couldn't apply owner, rolling back file creation", path, sudo)
			}
		}

		if mode != 0 {
			if err := linux_generic_file.ChangeMode(config, path, mode, sudo); err != nil {
				return rollbackFile(config, err, "Couldn't apply permissions, rolling back file creation", path, sudo)
			}
		}

		content := d.Get("content").(string)
		if content != "" {
			if err := linux_generic_file.WriteContent(config, path, content, sudo); err != nil {
				return rollbackFile(config, err, "Couldn't write content, rolling back file creation", path, sudo)
			}
		}

		d.SetId(path)
		return resourceFileRead()(d, m)
	}
}

func resourceFileRead() func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		config := m.(*sshclient.Config)
		id := d.Id()
		sudo := d.Get("sudo").(bool)

		owner, group, mode, err := linux_generic_file.GetDetails(config, id, sudo)
		if err != nil {
			if strings.Contains(err.Error(), "File not found with path") {
				d.SetId("")
				return nil
			}
			return errors.Wrap(err, "Unable to ls the file")
		}

		content, err := linux_generic_file.ReadFile(config, id, sudo)
		if err != nil {
			return errors.Wrap(err, "Unable to read the file")
		}
		d.Set("content", content)

		d.Set("owner", owner)
		d.Set("group", group)
		d.Set("mode", mode)
		return nil
	}
}

func resourceFileUpdate() func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		config := m.(*sshclient.Config)

		path := d.Get("path").(string)
		owner := d.Get("owner").(string)
		group := d.Get("group").(string)
		mode := d.Get("mode").(int)
		sudo := d.Get("sudo").(bool)

		oldPath := d.Id()
		oldOwner, oldGroup, oldMode, err := linux_generic_file.GetDetails(config, oldPath, sudo)
		if err != nil {
			return errors.Wrap(err, "Unable to ls the file")
		}

		content := d.Get("content").(string)
		oldContent, err := linux_generic_file.ReadFile(config, oldPath, sudo)
		if err != nil {
			return errors.Wrap(err, "Unable to read the file")
		}
		if oldContent != content {
			if err := linux_generic_file.WriteContent(config, oldPath, content, sudo); err != nil {
				return errors.Wrap(err, "Couldn't rewrite content")
			}
		}

		if oldPath != path {
			if err := linux_generic_file.Move(config, oldPath, path, sudo); err != nil {
				return errors.Wrap(err, "Couldn't mv file")
			}
			d.SetId(path)
		}

		if oldOwner != owner || oldGroup != group {
			if err := linux_generic_file.ChangeOwner(config, path, owner, group, sudo); err != nil {
				return errors.Wrap(err, "Couldn't apply owner or group")
			}
		}

		if oldMode != mode {
			if err := linux_generic_file.ChangeMode(config, path, mode, sudo); err != nil {
				return errors.Wrap(err, "Couldn't apply permissions")
			}
		}

		return resourceFileRead()(d, m)
	}
}

func resourceFileDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*sshclient.Config)
	id := d.Id()
	sudo := d.Get("sudo").(bool)
	return linux_generic_file.Remove(config, id, sudo)
}
