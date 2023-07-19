// Run a single command
package linux_generic_command

import (
	"fmt"
	"github.com/erwin-willems/terraform-provider-linuxbox/internal/sshclient"
)

// Check if a command exists and return true or false
func CommandExists(config *sshclient.Config, command string, sudo bool) (bool, error) {
	command = fmt.Sprintf("command -v %s", command)
	stdout, _, _ := sshclient.Run(config, sudo, command, "")
	return stdout != "", nil
}
