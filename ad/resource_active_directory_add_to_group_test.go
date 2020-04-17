package ad

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ldap "gopkg.in/ldap.v3"
)

func TestAccAddToGroup_Basic(t *testing.T) {
	var groupDN string = os.Getenv("AD_GROUP_OU_DISTINGUISHED_NAME")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccResourceAddToGroupPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddToGroupDestroy("ad_add_to_group.test"),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceAddToGroupConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAddToGroupExists("ad_add_to_group.test"),
					resource.TestCheckResourceAttr(
						"ad_add_to_group.test", "id", "CN=terraform3,"+groupDN+"|CN=terraform2,"+groupDN),
				),
			},
		},
	})
}

func testAccResourceAddToGroupPreCheck(t *testing.T) {
	if v := os.Getenv("AD_GROUP_OU_DISTINGUISHED_NAME"); v == "" {
		t.Fatal("AD_GROUP_OU_DISTINGUISHED_NAME must be set for acceptance tests")
	}
}

func testAccCheckAddToGroupDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Add To Group ID is set")
		}

		client := testAccProvider.Meta().(*ldap.Conn)
		targetGroup := rs.Primary.Attributes["target_group"]
		splitTargetGroup := strings.Split(targetGroup, ",") // split target group by commas
		searchRequest := ldap.NewSearchRequest(
			splitTargetGroup[len(splitTargetGroup)-2]+","+splitTargetGroup[len(splitTargetGroup)-1], // Make BaseDN from last two elements of split target group
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(|(objectCategory=user)(objectCategory=group))(memberOf="+targetGroup+"))", // Find users and groups that are members of targetGroup
			[]string{"dn"}, // A list attributes to retrieve
			nil,
		)

		searchResult, err := client.Search(searchRequest)
		if err != nil {
			return err
		}

		if len(searchResult.Entries) == 0 {
			return nil
		}

		return fmt.Errorf("test group still has members")
	}
}

func testAccCheckAddToGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Add To Group ID is set")
		}

		client := testAccProvider.Meta().(*ldap.Conn)
		targetGroup := rs.Primary.Attributes["target_group"]
		splitTargetGroup := strings.Split(targetGroup, ",") // split target group by commas
		searchRequest := ldap.NewSearchRequest(
			splitTargetGroup[len(splitTargetGroup)-2]+","+splitTargetGroup[len(splitTargetGroup)-1], // Make BaseDN from last two elements of split target group
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			"(&(|(objectCategory=user)(objectCategory=group))(memberOf="+targetGroup+"))", // Find users and groups that are members of targetGroup
			[]string{"dn"}, // A list attributes to retrieve
			nil,
		)

		searchResult, err := client.Search(searchRequest)
		if err != nil {
			return err
		}

		if len(searchResult.Entries) == 2 {
			return nil
		}

		return fmt.Errorf("expecting 2 members of test group")
	}
}

func testAccResourceAddToGroupConfig() string {
	return fmt.Sprintf(`
provider "ad" {
  domain   = "%s"
  ip       = "%s"
  url      = "%s"
  user     = "%s"
  password = "%s"  
}

resource "ad_group_to_ou" "test" {
  count = 3
  ou_distinguished_name = "%s"
  group_name = "terraform${count.index + 1}"
  description = "terraform test"
}

resource "ad_add_to_group" "test" {
  target_group = "CN=${ad_group_to_ou.test[0].group_name},${ad_group_to_ou.test[0].ou_distinguished_name}"
  dns_to_add = [
	"CN=${ad_group_to_ou.test[1].group_name},${ad_group_to_ou.test[1].ou_distinguished_name}",
	"CN=${ad_group_to_ou.test[2].group_name},${ad_group_to_ou.test[2].ou_distinguished_name}",
  ]
}`,
		os.Getenv("AD_DOMAIN"),
		os.Getenv("AD_IP"),
		os.Getenv("AD_URL"),
		os.Getenv("AD_USER"),
		os.Getenv("AD_PASSWORD"),
		os.Getenv("AD_GROUP_OU_DISTINGUISHED_NAME"))
}
