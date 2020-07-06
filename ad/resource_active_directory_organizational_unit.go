package ad

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	ldap "gopkg.in/ldap.v3"
)

func resourceOU() *schema.Resource {
	return &schema.Resource{
		Create: ressourceADOUCreate,
		Read:   resourceADOURead,
		Delete: resourceADOUDelete,

		Schema: map[string]*schema.Schema{
			"ou_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ou_distinguished_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				ForceNew: true,
			},
		},
	}
}

func ressourceADOUCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ldap.Conn) // m is our client to talk to server
	ouName := d.Get("ou_name").(string)
	OUDistinguishedName := d.Get("ou_distinguished_name").(string)
	description := d.Get("description").(string)
	var dnOfOU string
	dnOfOU += "OU=" + ouName + "," + OUDistinguishedName //object's entire path to the root
	log.Printf("[DEBUG] dnOfOU: %s ", dnOfOU)
	log.Printf("[DEBUG] Adding OU : %s ", ouName)
	err := addOU(ouName, dnOfOU, client, description)
	if err != nil {
		log.Printf("[ERROR] Error while adding OU: %s ", err)
		return fmt.Errorf("Error while adding OU %s", err)
	}
	log.Printf("[DEBUG] OU Added successfully: %s", ouName)
	d.SetId(OUDistinguishedName + "/" + ouName)
	return nil

}

func resourceADOURead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ldap.Conn)

	ouName := d.Get("ou_name").(string)
	OUDistinguishedName := d.Get("ou_distinguished_name").(string)

	var dnOfOU string
	dnOfOU += OUDistinguishedName

	log.Printf("[DEBUG] Searching OU with domain: %s ", dnOfOU)

	NewReq := ldap.NewSearchRequest( //represents the search request send to the server
		dnOfOU, // base dnOfOU.
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=OrganizationalUnit)(ou="+ouName+"))", //applied filter
		[]string{"ou", "dn"},
		nil,
	)

	sr, err := client.Search(NewReq)
	if err != nil {
		log.Printf("[ERROR] while seaching OU : %s", err)
		return fmt.Errorf("Error while searching  OU : %s", err)
	}

	log.Println("[DEBUG] Found " + strconv.Itoa(len(sr.Entries)) + " Entries")
	for _, entry := range sr.Entries {
		log.Printf("[DEBUG] %s: %v\n", entry.DN, entry.GetAttributeValue("ou"))

	}

	if len(sr.Entries) == 0 {
		log.Println("[DEBUG] OU not found")
		d.SetId("")
	}
	return nil
}

func resourceADOUDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[ERROR] Finding OU")
	resourceADOURead(d, m)
	if d.Id() == "" {
		log.Println("[ERROR] Cannot find OU in the specified AD")
		return fmt.Errorf("[ERROR] Cannot find OU in the specified AD")
	}
	client := m.(*ldap.Conn)

	ouName := d.Get("ou_name").(string)
	OUDistinguishedName := d.Get("ou_distinguished_name").(string)

	var dnOfOU string
	dnOfOU += "OU=" + ouName + "," + OUDistinguishedName
	log.Printf("[DEBUG] Name of the DN is : %s ", dnOfOU)
	log.Printf("[DEBUG] Deleting the OU from the AD : %s ", ouName)

	err := deleteOU(dnOfOU, client)
	if err != nil {
		log.Printf("[ERROR] Error while Deleting OU from AD : %s ", err)
		return fmt.Errorf("Error while Deleting OU from AD %s", err)
	}
	log.Printf("[DEBUG] OU deleted from AD successfully: %s", ouName)
	return nil

}
