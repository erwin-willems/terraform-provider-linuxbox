package linuxbox

import (
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testAccPreCheck(t *testing.T) {
	//testAccProvider := provider.Provider()
	// err := testAccProvider.Configure(terraform.NewResourceConfigRaw(nil))
	// if err != nil {
	// 	t.Fatal(err)
	// }
}

func TestAccFileCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: fileCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", ""),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "644"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccFileWithContentCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: fileWithContentCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", "abc"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "644"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccFileWithOwnerCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: fileWithOwnerCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", ""),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "644"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccFileWithGroupCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: fileWithGroupCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", ""),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "644"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccFileWithPermissionsCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: fileWithPermissionsCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", ""),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "640"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
		},
	})
}

func TestAccFileWithAllAttrsCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: fileWithAllAttrsCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "640"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
		},
	})
}


func TestAccFileUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: fileWithAllAttrsCreationConfig,
			},
			{
				Config: fileWithOwnerUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "640"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
			{
				Config: fileWithOwnerPermissionsUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "760"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
			{
				Config: fileWithOwnerPermissionsContentUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", "this was mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "760"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
			{
				Config: fileWithOwnerPermissionsContentPathUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile2"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", "this was mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "760"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
			{
				Config: fileWithOwnerPermissionsContentPathUpdatedAllConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "path", "/opt/testfile3"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "content", "this is still mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "owner", "mail"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "group", "root"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "mode", "755"),
					resource.TestCheckResourceAttr("linuxbox_file.testfile", "sudo", "false"),
				),
			},
		},
	})
}



const fileCreationConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile"
}
`
const fileWithContentCreationConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile"
  content = "abc"
}
`
const fileWithOwnerCreationConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile"
  owner = "mail"
}
`
const fileWithGroupCreationConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile"
  group = "mail"
}
`
const fileWithPermissionsCreationConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile"
  mode = 640
}
`
const fileWithAllAttrsCreationConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile"
  owner = "mail"
  group = "mail"
  mode = 640
  content = "mail"
}
`
const fileWithOwnerUpdatedConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile"
  owner = "root"
  group = "mail"
  mode = 640
  content = "mail"
}
`
const fileWithOwnerPermissionsUpdatedConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile"
  owner = "root"
  group = "mail"
  mode = 760
  content = "mail"
}
`
const fileWithOwnerPermissionsContentUpdatedConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile"
  owner = "root"
  group = "mail"
  mode = 760
  content = "this was mail"
}
`
const fileWithOwnerPermissionsContentPathUpdatedConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile2"
  owner = "root"
  group = "mail"
  mode = 760
  content = "this was mail"
}
`
const fileWithOwnerPermissionsContentPathUpdatedAllConfig = `
resource "linuxbox_file" "testfile" {
  path = "/opt/testfile3"
  owner = "mail"
  group = "root"
  mode = 755
  content = "this is still mail"
}
`

