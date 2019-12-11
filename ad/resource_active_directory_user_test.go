package ad

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/ldap.v2"
)
//test function:
func TestAccAdUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccResourceAdUserPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdUserDestroy("ad_user.test"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceAdUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdUserExists("ad_user.test"),
					resource.TestCheckResourceAttr(
						"ad_user.test", "logon_name", "terraform"),
				),
			},
		},
	})
}

func testAccResourceAdUserPreCheck(t *testing.T) {
	if v := os.Getenv("AD_USER_DOMAIN"); v == "" {
		t.Fatal("AD_USER_DOMAIN must be set for acceptance tests")
	}
}

func testAccCheckAdUserDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AD User ID is set")
		}
		client := testAccProvider.Meta().(*ldap.Conn)
		domain := rs.Primary.Attributes["domain"]
		var dnOfUser string
		domainArr := strings.Split(domain, ".")
		dnOfUser = "dc=" + domainArr[0]
		for index, item := range domainArr {
			if index == 0 {
				continue
			}
			dnOfUser += ",dc=" + item
		}
		searchRequest := ldap.NewSearchRequest(
			dnOfUser,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(objectClass=User)(cn="+rs.Primary.Attributes["name"]+"))", // The filter to apply
			[]string{"dn", "cn"},                                          // A list attributes to retrieve
			nil,
		)
		sr, err := client.Search(searchRequest)
		if err != nil {
			return err
		}
		if len(sr.Entries) == 0 {
			return nil
		}

		return fmt.Errorf("User AD still exists")
	}

}

func testAccCheckAdUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AD User ID is set")
		}
		client := testAccProvider.Meta().(*ldap.Conn)
		domain := rs.Primary.Attributes["domain"]
		var dnOfUser string
		domainArr := strings.Split(domain, ".")
		dnOfUser = "dc=" + domainArr[0]
		for index, item := range domainArr {
			if index == 0 {
				continue
			}
			dnOfUser += ",dc=" + item
		}
		searchRequest := ldap.NewSearchRequest(
			dnOfUser,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(objectClass=User)(cn="+rs.Primary.Attributes["name"]+"))", // The filter to apply
			[]string{"dn", "cn"},                                          // A list attributes to retrieve
			nil,
		)
		sr, err := client.Search(searchRequest)
		if err != nil {
			return err
		}
		if len(sr.Entries) > 0 {
			return nil
		}
		return nil
	}
}

func testAccResourceAdUserConfig() string {
	return fmt.Sprintf(`
provider "ad" {
  domain   = "%s"
  ip       = "%s"
  user     = "%s"
  password = "%s"
}
resource "ad_user" "test" {
  domain = "%s"
  logon_name = "terraform"
}`,
		os.Getenv("AD_DOMAIN"),
		os.Getenv("AD_IP"),
		os.Getenv("AD_USER"),
		os.Getenv("AD_PASSWORD"),
		os.Getenv("AD_USER_DOMAIN"))
}
