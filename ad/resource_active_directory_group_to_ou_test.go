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
					resource.TestCheckResourceAttr(
						"ad_group_to_ou.test", "gid_number", "9001"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccResourceAdGroupToOUPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdGroupToOUDestroy("ad_group_to_ou.test"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceAdGroupToOUConfig_with_auto_gid(),
				Check: resource.ComposeTestCheckFunc(
					// make sure we get gidNumber 9001 when 9000 is used already.
					resource.TestCheckResourceAttr(
						"ad_group_to_ou.test", "auto_gid_number", "9001"),
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
			[]string{"dn"}, // A list attributes to retrieve
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
			// filter for the test group by cn and gidNumber
			"(&(objectClass=Group)(cn="+rs.Primary.Attributes["group_name"]+")(gidNumber="+rs.Primary.Attributes["gid_number"]+"))", // The filter to apply
			[]string{"dn"}, // A list attributes to retrieve
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
	url      = "%s"
  user     = "%s"
  password = "%s"  
}

resource "ad_group_to_ou" "test" {
  ou_distinguished_name = "%s"
  group_name = "terraform"
  description = "terraform test"
  gid_number = "9001"
}`,
		os.Getenv("AD_DOMAIN"),
		os.Getenv("AD_IP"),
		os.Getenv("AD_URL"),
		os.Getenv("AD_USER"),
		os.Getenv("AD_PASSWORD"),
		os.Getenv("AD_GROUP_OU_DISTINGUISHED_NAME"))
}

func testAccResourceAdGroupToOUConfig_with_auto_gid() string {
	return fmt.Sprintf(`
provider "ad" {
  domain   = "%s"
	ip       = "%s"
	url      = "%s"
  user     = "%s"
  password = "%s"  
}

resource "ad_group_to_ou" "static_gid" {
  ou_distinguished_name = "%s"
  group_name = "terraform9000"
  description = "terraform test"
  gid_number = "9000"
}

resource "ad_group_to_ou" "test" {
  ou_distinguished_name = "%[5]s"
  group_name = "terraform9001"
  description = "terraform test"
  auto_gid = true
  auto_gid_min = 9000
  auto_gid_max = 9001
}`,
		os.Getenv("AD_DOMAIN"),
		os.Getenv("AD_IP"),
		os.Getenv("AD_URL"),
		os.Getenv("AD_USER"),
		os.Getenv("AD_PASSWORD"),
		os.Getenv("AD_GROUP_OU_DISTINGUISHED_NAME"))
}
