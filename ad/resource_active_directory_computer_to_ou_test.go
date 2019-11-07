package ad

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/ldap.v3"
)

func TestAccAdComputerToOU_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccResourceAdComputerToOUPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdComputerToOUDestroy("ad_computer_to_ou.test"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceAdComputerToOUConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdComputerToOUExists("ad_computer_to_ou.test"),
					resource.TestCheckResourceAttr(
						"ad_computer_to_ou.test", "computer_name", "terraform"),
				),
			},
		},
	})
}

func testAccResourceAdComputerToOUPreCheck(t *testing.T) {
	if v := os.Getenv("AD_COMPUTER_OU_DISTINGUISHED_NAME"); v == "" {
		t.Fatal("AD_COMPUTER_OU_DISTINGUISHED_NAME must be set for acceptance tests")
	}
}

func testAccCheckAdComputerToOUDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AD Computer ID is set")
		}
		client := testAccProvider.Meta().(*ldap.Conn)
		ouDistinguishedName := rs.Primary.Attributes["ou_distinguished_name"]
		var dnOfComputer string
		dnOfComputer = ouDistinguishedName
		searchRequest := ldap.NewSearchRequest(
			dnOfComputer, //"cn=code1,ou=DevComputers,dc=terraform,dc=local", // The base dn to search
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

func testAccCheckAdComputerToOUExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AD Computer ID is set")
		}
		client := testAccProvider.Meta().(*ldap.Conn)
		ouDistinguishedName := rs.Primary.Attributes["ou_distinguished_name"]
		var dnOfComputer string
		dnOfComputer = ouDistinguishedName
		searchRequest := ldap.NewSearchRequest(
			dnOfComputer, //"cn=code1,ou=DevComputers,dc=terraform,dc=local", // The base dn to search
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

func testAccResourceAdComputerToOUConfig() string {
	return fmt.Sprintf(`
provider "ad" {
  domain   = "%s"
  ip       = "%s"
  user     = "%s"
  password = "%s"  
}

resource "ad_computer_to_ou" "test" {
  ou_distinguished_name = "%s"
  computer_name = "terraform"
  description = "terraform test"
}`,
		os.Getenv("AD_DOMAIN"),
		os.Getenv("AD_IP"),
		os.Getenv("AD_USER"),
		os.Getenv("AD_PASSWORD"),
		os.Getenv("AD_COMPUTER_OU_DISTINGUISHED_NAME"))
}
