package ad

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

//test function:
func TestAccAdDataSourceUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAdDataSourceUserConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ad_user.test", "logon_name", "test"),
					resource.TestCheckResourceAttrSet("data.ad_user.test", "dn"),
				),
			},
		},
	})
}

func testAccAdDataSourceUserConfig() string {
	return fmt.Sprintf(`
provider "ad" {
  domain   = "%s"
  ip       = "%s"
  url      = "%s"
  user     = "%s"
  password = "%s"
}
data "ad_user" "test" {
  domain = "%s"
  logon_name = "test"
}`,
		os.Getenv("AD_DOMAIN"),
		os.Getenv("AD_IP"),
		os.Getenv("AD_URL"),
		os.Getenv("AD_USER"),
		os.Getenv("AD_PASSWORD"),
		os.Getenv("AD_USER_DOMAIN"))
}
