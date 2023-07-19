package linuxbox

import (
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)


func TestAccDirectoryCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: directoryCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "path", "/opt/testdir"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "group", "root"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "mode", "755"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccDirectoryWithOwnerCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: directoryWithOwnerCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "path", "/opt/testdir"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "owner", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "group", "root"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "mode", "755"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccDirectoryWithGroupCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: directoryWithGroupCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "path", "/opt/testdir"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "mode", "755"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccDirectoryWithModeCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: directoryWithModeCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "path", "/opt/testdir"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "group", "root"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "mode", "700"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccDirectoryWithAllAttrsCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: directoryWithAllAttrsCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "path", "/opt/testdir"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "owner", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "mode", "700"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccDirectoryWithOwnerUpdated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: directoryCreationConfig,
			},
			{
				Config: directoryWithOwnerUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "path", "/opt/testdir"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "owner", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "group", "root"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "mode", "755"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "sudo", "false"),
				),
			},
			{
				Config: directoryWithOwnerGroupUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "path", "/opt/testdir"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "owner", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "mode", "755"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "sudo", "false"),
				),
			},
			{
				Config: directoryWithOwnerGroupModeUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "path", "/opt/testdir"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "owner", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "mode", "700"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "sudo", "false"),
				),
			},
			{
				Config: directoryWithOwnerGroupModePathUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "path", "/opt/testdir2"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "owner", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "mode", "700"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "sudo", "false"),
				),
			},
			{
				Config: directoryWithAllUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "path", "/opt/testdir3"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "group", "root"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "mode", "750"),
					resource.TestCheckResourceAttr("linuxbox_directory.testdir", "sudo", "false"),
				),
			},
			
		},
	})
}

const directoryCreationConfig = `
resource "linuxbox_directory" "testdir" {
  path = "/opt/testdir"
}
`
const directoryWithOwnerCreationConfig = `
resource "linuxbox_directory" "testdir" {
  path = "/opt/testdir"
  owner = "mail"
}
`
const directoryWithGroupCreationConfig = `
resource "linuxbox_directory" "testdir" {
  path = "/opt/testdir"
  group = "mail"
}
`
const directoryWithModeCreationConfig = `
resource "linuxbox_directory" "testdir" {
  path = "/opt/testdir"
  mode = "700"
}
`
// Test with all parameters
const directoryWithAllAttrsCreationConfig = `
resource "linuxbox_directory" "testdir" {
  path = "/opt/testdir"
  owner = "mail"
  group = "mail"
  mode = "700"
}
`



// Tests for update
const directoryWithOwnerUpdatedConfig = `
resource "linuxbox_directory" "testdir" {
  path = "/opt/testdir"
  owner = "mail"
}
`
const directoryWithOwnerGroupUpdatedConfig = `
resource "linuxbox_directory" "testdir" {
  path = "/opt/testdir"
  owner = "mail"
  group = "mail"
}
`
const directoryWithOwnerGroupModeUpdatedConfig = `
resource "linuxbox_directory" "testdir" {
  path = "/opt/testdir"
  owner = "mail"
  group = "mail"
  mode = "700"
}
`
const directoryWithOwnerGroupModePathUpdatedConfig = `
resource "linuxbox_directory" "testdir" {
  path = "/opt/testdir2"
  owner = "mail"
  group = "mail"
  mode = "700"
}
`
const directoryWithAllUpdatedConfig = `
resource "linuxbox_directory" "testdir" {
  path = "/opt/testdir3"
  owner = "root"
  group = "root"
  mode = "750"
}
`
