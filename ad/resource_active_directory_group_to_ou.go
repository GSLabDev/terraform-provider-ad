package ad

import (
	"fmt"
	"log"
	"strconv"

	ldap "gopkg.in/ldap.v3"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroupToOU() *schema.Resource {
	return &schema.Resource{
		Create: resourceADGroupToOUCreate,
		Read:   resourceADGroupToOURead,
		Delete: resourceADGroupToOUDelete,
		Schema: map[string]*schema.Schema{
			"group_name": {
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

func resourceADGroupToOUCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ldap.Conn)

	groupName := d.Get("group_name").(string)
	OUDistinguishedName := d.Get("ou_distinguished_name").(string)
	description := d.Get("description").(string)
	var dnOfGroup string
	dnOfGroup += "cn=" + groupName + "," + OUDistinguishedName

	log.Printf("[DEBUG] Name of the DN is : %s ", dnOfGroup)
	log.Printf("[DEBUG] Adding the Group to the AD : %s ", groupName)

	err := addGroupToAD(groupName, dnOfGroup, client, description)
	if err != nil {
		log.Printf("[ERROR] Error while adding a Group to the AD : %s ", err)
		return fmt.Errorf("Error while adding a Group to the AD %s", err)
	}
	log.Printf("[DEBUG] Group Added to AD successfully: %s", groupName)
	d.SetId(OUDistinguishedName + "/" + groupName)
	return nil
}

func resourceADGroupToOURead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[ERROR] In Read function")
	client := meta.(*ldap.Conn)

	groupName := d.Get("group_name").(string)
	OUDistinguishedName := d.Get("ou_distinguished_name").(string)
	var dnOfGroup string
	dnOfGroup += OUDistinguishedName

	log.Printf("[DEBUG] Name of the DN is : %s ", dnOfGroup)
	log.Printf("[DEBUG] Searching the Group from the AD : %s ", groupName)

	searchRequest := ldap.NewSearchRequest(
		dnOfGroup, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=Group)(cn="+groupName+"))", // The filter to apply
		[]string{"dn", "cn"},                       // A list attributes to retrieve
		nil,
	)

	sr, err := client.Search(searchRequest)
	if err != nil {
		log.Printf("[ERROR] Error while searching a Group : %s ", err)
		return fmt.Errorf("Error while searching a Group : %s", err)
	}
	fmt.Println("[ERROR] Found " + strconv.Itoa(len(sr.Entries)) + " Entries")
	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}
	if len(sr.Entries) == 0 {
		log.Println("[ERROR] Group was not found")
		d.SetId("")
	}
	return nil
}

func resourceADGroupToOUDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[ERROR] Finding group")
	resourceADGroupToOURead(d, meta)
	if d.Id() == "" {
		log.Println("[ERROR] Cannot find Group in the specified AD")
		return fmt.Errorf("[ERROR] Cannot find Group in the specified AD")
	}
	client := meta.(*ldap.Conn)

	groupName := d.Get("group_name").(string)
	OUDistinguishedName := d.Get("ou_distinguished_name").(string)
	var dnOfGroup string
	dnOfGroup += "cn=" + groupName + "," + OUDistinguishedName

	log.Printf("[DEBUG] Name of the DN is : %s ", dnOfGroup)
	log.Printf("[DEBUG] Deleting the Group from the AD : %s ", groupName)

	err := deleteGroupFromAD(dnOfGroup, client)
	if err != nil {
		log.Printf("[ERROR] Error while Deleting a Group from AD : %s ", err)
		return fmt.Errorf("Error while Deleting a Group from AD %s", err)
	}
	log.Printf("[DEBUG] Group deleted from AD successfully: %s", groupName)
	return nil
}
