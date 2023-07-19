package linuxbox

import (
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

)

func TestAccUserCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: userCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1000"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1000"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/home/testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/bash"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccUserWithUid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: userWithUidConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1001"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1001"), // Linux tries to keep the gid the same as uid
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/home/testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/bash"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccUserWithGid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: userWithGidConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1000"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1001"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/home/testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/bash"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccUserWithHome(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: userWithHomeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1000"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1000"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/opt/testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/bash"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccUserWithShell(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: userWithShellConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1000"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1000"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/home/testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/sh"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccUserUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: userCreationConfig,
			},
			// TODO Fix me
			// {
			// 	Config: userWithNameUpdatedConfig,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser1"),
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1000"),
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1001"),
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/home/testuser"),
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/bash"),
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
			// 	),
			// },
			
			// TODO When userWithGid is fixed, all expected names must be set to testuser1
			{
				Config: userWithGidUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1000"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1001"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/home/testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/bash"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
				),
			},
			{
				Config: userWithGidUidUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1001"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1001"), // Linux tries to keep the gid the same as uid
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/home/testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/bash"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
				),
			},
			{
				Config: userWithGidUidHomeUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1001"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1001"), // Linux tries to keep the gid the same as uid
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/opt/testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/bash"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
				),
			},
			{
				Config: userWithGidUidHomeShellUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1001"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1001"), // Linux tries to keep the gid the same as uid
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/opt/testuser"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/sh"),
					resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
				),
			},
			// Todo: fixme
			// {
			// 	Config: userWithAllUpdatedConfig,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "name", "testuser2"),
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "uid", "1002"),
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "gid", "1002"), // Linux tries to keep the gid the same as uid
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "home", "/home/testuser2"),
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "shell", "/bin/bash"),
			// 		resource.TestCheckResourceAttr("linuxbox_user.testuser", "sudo", "false"),
			// 	),
			// },
		},
	})
}




const userCreationConfig = `
resource "linuxbox_user" "testuser" {
	name = "testuser"
}
`

const userWithUidConfig = `
resource "linuxbox_user" "testuser" {
	name = "testuser"
	uid = 1001
}
`

const userWithGidConfig = `
resource "linuxbox_user" "testuser" {
	name = "testuser"
	gid = 1001
}
`

const userWithHomeConfig = `
resource "linuxbox_user" "testuser" {
	name = "testuser"
	home = "/opt/testuser"
}
`

const userWithShellConfig = `
resource "linuxbox_user" "testuser" {
	name = "testuser"
	shell = "/bin/sh"
}
`


// const userWithNameUpdatedConfig = `
// resource "linuxbox_user" "testuser" {
// 	name = "testuser1"
// 	gid = 1001
// }
// `

// Todo: Once userWithNameUpdate is fixed, all usernames below must be set to testuser1
const userWithGidUpdatedConfig = `
resource "linuxbox_user" "testuser" {
	name = "testuser"
	gid = 1001
}
`

const userWithGidUidUpdatedConfig = `
resource "linuxbox_user" "testuser" {
	name = "testuser"
	uid = 1001
	gid = 1001
}
`

const userWithGidUidHomeUpdatedConfig = `
resource "linuxbox_user" "testuser" {
	name = "testuser"
	uid = 1001
	gid = 1001
	home = "/opt/testuser"
}
`

const userWithGidUidHomeShellUpdatedConfig = `
resource "linuxbox_user" "testuser" {
	name = "testuser"
	uid = 1001
	gid = 1001
	home = "/opt/testuser"
	shell = "/bin/sh"
}
`

// const userWithAllUpdatedConfig = `
// resource "linuxbox_user" "testuser" {
// 	name = "testuser2"
// 	uid = 1002
// 	gid = 1002
// 	home = "/home/testuser2"
// 	shell = "/bin/bash"
// }
// `
