package linuxbox

import (
	"testing"
	"regexp"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var version = regexp.MustCompile(`\d\.\d\.\d`)

func TestAccPackageCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: packageSingleCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "name", "bzip2"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "state", "present"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "version", ""),
					resource.TestMatchResourceAttr("linuxbox_package.testpackage", "installed_version", version),
				),
			},
		},
	})
}

func TestAccPackageListCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: packageListCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "name", "list_76cd2aa57a276975451fbfc7813f3857"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "names.0", "bzip2"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "names.1", "bc"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "state", "present"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "version", ""),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "installed_version", ""),
				),
			},
		},
	})
}


func TestAccPackageWithStateCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: packageWithStateCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "name", "bzip2"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "state", "absent"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "version", ""),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "installed_version", ""),
				),
			},
		},
	})
}

func TestAccPackageUpdated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: packageSingleCreationConfig,
			},
			{
				Config: packageWithNameUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "name", "bc"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "state", "present"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "version", ""),
					resource.TestMatchResourceAttr("linuxbox_package.testpackage", "installed_version", version),
				),
			},
		},
	})
}

func TestAccPackageUpdated2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{ "linuxbox": Provider() },
		Steps: []resource.TestStep{
			{
				Config: packageListCreationConfig,
			},
			{
				Config: packageWithNamesUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "name", "bc"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "state", "present"),
					resource.TestCheckResourceAttr("linuxbox_package.testpackage", "version", ""),
					resource.TestMatchResourceAttr("linuxbox_package.testpackage", "installed_version", version),
				),
			},
		},
	})
}


// Note:
// The scenario for parameter "version" is hard to test because it can change at any time.
const packageSingleCreationConfig = `
resource "linuxbox_package" "testpackage" {
	name = "bzip2"
}
`
const packageListCreationConfig = `
resource "linuxbox_package" "testpackage" {
	names = ["bzip2", "bc"]
}
`
const packageWithStateCreationConfig = `
resource "linuxbox_package" "testpackage" {
	name = "bzip2"
	state = "absent"
}
`
const packageWithNameUpdatedConfig = `
resource "linuxbox_package" "testpackage" {
	name = "bc"
}
`
const packageWithNamesUpdatedConfig = `
resource "linuxbox_package" "testpackage" {
	names = ["geolite2-city", "geolite2-country"]
}
`
