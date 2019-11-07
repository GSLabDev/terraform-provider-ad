package ad

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/ldap.v3"
)

func TestAccAdComputer_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccResourceAdComputerPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdComputerDestroy("ad_computer.test"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceAdComputerConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdComputerExists("ad_computer.test"),
					resource.TestCheckResourceAttr(
						"ad_computer.test", "computer_name", "terraform"),
				),
			},
		},
	})
}

func testAccResourceAdComputerPreCheck(t *testing.T) {
	if v := os.Getenv("AD_COMPUTER_DOMAIN"); v == "" {
		t.Fatal("AD_COMPUTER_DOMAIN must be set for acceptance tests")
	}
}

func testAccCheckAdComputerDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AD Computer ID is set")
		}
		client := testAccProvider.Meta().(*ldap.Conn)
		domain := rs.Primary.Attributes["domain"]
		var dnOfComputer string
		domainArr := strings.Split(domain, ".")
		dnOfComputer = "dc=" + domainArr[0]
		for index, item := range domainArr {
			if index == 0 {
				continue
			}
			dnOfComputer += ",dc=" + item
		}
		searchRequest := ldap.NewSearchRequest(
			dnOfComputer, //"cn=code1,cn=Computers,dc=terraform,dc=local", // The base dn to search
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(objectClass=Computer)(cn="+rs.Primary.Attributes["computer_name"]+"))", // The filter to apply
			[]string{"dn", "cn"}, // A list attributes to retrieve
			nil,
		)
		sr, err := client.Search(searchRequest)
		if err != nil {
			return err
		}
		if len(sr.Entries) == 0 {
			return nil
		}

		return fmt.Errorf("Computer AD still exists")
	}

}

func testAccCheckAdComputerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AD Computer ID is set")
		}
		client := testAccProvider.Meta().(*ldap.Conn)
		domain := rs.Primary.Attributes["domain"]
		var dnOfComputer string
		domainArr := strings.Split(domain, ".")
		dnOfComputer = "dc=" + domainArr[0]
		for index, item := range domainArr {
			if index == 0 {
				continue
			}
			dnOfComputer += ",dc=" + item
		}
		searchRequest := ldap.NewSearchRequest(
			dnOfComputer, //"cn=code1,cn=Computers,dc=terraform,dc=local", // The base dn to search
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(objectClass=Computer)(cn="+rs.Primary.Attributes["computer_name"]+"))", // The filter to apply
			[]string{"dn", "cn"}, // A list attributes to retrieve
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

func testAccResourceAdComputerConfig() string {
	return fmt.Sprintf(`
provider "ad" {
  domain   = "%s"
  ip       = "%s"
  user     = "%s"
  password = "%s"
}
resource "ad_computer" "test" {
  domain = "%s"
  computer_name = "terraform"
  description = "terraform test"
}`,
		os.Getenv("AD_DOMAIN"),
		os.Getenv("AD_IP"),
		os.Getenv("AD_USER"),
		os.Getenv("AD_PASSWORD"),
		os.Getenv("AD_COMPUTER_DOMAIN"))
}
