package ad

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	ldap "gopkg.in/ldap.v2"
)

func dataActiveDirectoryUsers() *schema.Resource {
	return &schema.Resource{
		Read: adUserReadBySearch,

		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"base_search_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"username_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "*",
			},
			"attributes": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func adUserReadBySearch(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ldap.Conn)

	baseSearchDn := d.Get("base_search_dn").(string)
	usernameFilter := d.Get("username_filter").(string)

	var attributes []string
	if v, ok := d.GetOk("attributes"); ok {
		attributes = expandStringSlice(v.([]interface{}))
	}

	searchRequest := ldap.NewSearchRequest(
		baseSearchDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=User)(cn="+usernameFilter+"))", // The filter to apply
		attributes, // A list attributes to retrieve
		nil,
	)

	sr, err := client.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	userProperties := map[string]string{}

	entry := sr.Entries[0]

	for _, properties := range entry.Attributes {
		userProperties[properties.Name] = properties.Values[0]
	}

	d.SetId(baseSearchDn)
	d.Set("user", userProperties)

	return nil
}
