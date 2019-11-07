package ad

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/ldap.v3"
)

func TestAccAdGroupToOU_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccResourceAdGroupToOUPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdGroupToOUDestroy("ad_group_to_ou.test"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceAdGroupToOUConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdGroupToOUExists("ad_group_to_ou.test"),
					resource.TestCheckResourceAttr(
						"ad_group_to_ou.test", "group_name", "terraform"),
				),
			},
		},
	})
}

func testAccResourceAdGroupToOUPreCheck(t *testing.T) {
	if v := os.Getenv("AD_GROUP_OU_DISTINGUISHED_NAME"); v == "" {
		t.Fatal("AD_GROUP_OU_DISTINGUISHED_NAME must be set for acceptance tests")
	}
}

func testAccCheckAdGroupToOUDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AD Group ID is set")
		}
		client := testAccProvider.Meta().(*ldap.Conn)
		ouDistinguishedName := rs.Primary.Attributes["ou_distinguished_name"]
		var dnOfGroup string
		dnOfGroup = ouDistinguishedName
		searchRequest := ldap.NewSearchRequest(
			dnOfGroup, //"cn=code1,ou=DevGroups,dc=terraform,dc=local", // The base dn to search
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(objectClass=Group)(cn="+rs.Primary.Attributes["group_name"]+"))", // The filter to apply
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

		return fmt.Errorf("Group AD still exists")
	}

}

func testAccCheckAdGroupToOUExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AD Group ID is set")
		}
		client := testAccProvider.Meta().(*ldap.Conn)
		ouDistinguishedName := rs.Primary.Attributes["ou_distinguished_name"]
		var dnOfGroup string
		dnOfGroup = ouDistinguishedName
		searchRequest := ldap.NewSearchRequest(
			dnOfGroup, //"cn=code1,ou=DevGroups,dc=terraform,dc=local", // The base dn to search
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(objectClass=Group)(cn="+rs.Primary.Attributes["group_name"]+"))", // The filter to apply
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

func testAccResourceAdGroupToOUConfig() string {
	return fmt.Sprintf(`
provider "ad" {
  domain   = "%s"
  ip       = "%s"
  user     = "%s"
  password = "%s"  
}

resource "ad_group_to_ou" "test" {
  ou_distinguished_name = "%s"
  group_name = "terraform"
  description = "terraform test"
}`,
		os.Getenv("AD_DOMAIN"),
		os.Getenv("AD_IP"),
		os.Getenv("AD_USER"),
		os.Getenv("AD_PASSWORD"),
		os.Getenv("AD_GROUP_OU_DISTINGUISHED_NAME"))
}
