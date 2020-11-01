package ad

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	ldap "gopkg.in/ldap.v3"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceADUserRead,

		Schema: map[string]*schema.Schema{
			"logon_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DN of user",
			},
		},
	}
}

func dataSourceADUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ldap.Conn)
	domain := d.Get("domain").(string)
	logonName := d.Get("logon_name").(string)

	var dnOfUser string // dnOfUser: distingished names uniquely identifies an entry to AD.
	domainArr := strings.Split(domain, ".")
	dnOfUser = "dc=" + domainArr[0]
	for index, i := range domainArr {
		if index == 0 {
			continue
		}
		dnOfUser += ",dc=" + i
	}

	NewReq := ldap.NewSearchRequest(
		dnOfUser, // base dnOfUser.
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0,
		false,
		"(&(objectClass=User)(sAMAccountName="+logonName+"))", //applied filter
		[]string{"dn"},
		nil,
	)

	sr, err := client.Search(NewReq)
	if err != nil {
		log.Printf("[ERROR] while seaching user : %s", err)
		return fmt.Errorf("Error while searching  user : %s", err)
	}

	fmt.Println("[ERROR] Found " + strconv.Itoa(len(sr.Entries)) + " Entries")
	for _, entry := range sr.Entries {
		d.SetId(logonName)
		d.Set("dn", entry.DN)
		log.Printf("[DEBUG] ### %s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}

	if len(sr.Entries) == 0 {
		log.Println("[ERROR] user not found")
		d.SetId("")
	}
	return nil
}
