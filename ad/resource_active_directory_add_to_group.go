package ad

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	ldap "gopkg.in/ldap.v3"
)

func resourceAddToGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAddToGroupCreate,
		Read:   resourceAddToGroupRead,
		Delete: resourceAddToGroupDelete,
		Schema: map[string]*schema.Schema{
			"dns_to_add": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "A list of distinguished names to add to target_group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Required:    true,
				ForceNew:    true,
			},
			"target_group": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The distinguished name of the target group you're adding members to.",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceAddToGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ldap.Conn)

	var resourceId string

	targetGroup := d.Get("target_group").(string)
	// get each DN out of the Set and interate to do the work
	for _, distinguishedName := range d.Get("dns_to_add").(*schema.Set).List() {
		distinguishedNameToAdd := fmt.Sprintf("%s", distinguishedName)
		log.Printf("[DEBUG] Adding %s to %s", distinguishedNameToAdd, targetGroup)
		// call the helper to do the work
		err := addToGroup(distinguishedNameToAdd, targetGroup, client)
		if err != nil {
			log.Printf("[ERROR] Error while adding %s to %s : %s", distinguishedNameToAdd, targetGroup, err)
			return fmt.Errorf("Error while adding %s to %s : %s", distinguishedNameToAdd, targetGroup, err)
		}
		resourceId = resourceId + "|" + distinguishedNameToAdd
		log.Printf("[DEBUG] Successfully added %s to %s.", distinguishedNameToAdd, targetGroup)
	}
	resourceId = strings.TrimPrefix(resourceId, "|")
	d.SetId(resourceId)

	return nil
}

func resourceAddToGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ldap.Conn)
	targetGroup := strings.ToLower(d.Get("target_group").(string))
	log.Printf("[DEBUG] Searching for members of %s", targetGroup)
	splitTargetGroup := strings.SplitN(targetGroup, "dc", 2) // split target group
	searchRequest := ldap.NewSearchRequest(
		"dc"+splitTargetGroup[1], // Make BaseDN from dc elements of split target group
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(|(objectCategory=user)(objectCategory=group))(memberOf="+targetGroup+"))", // Find users and groups that are members of targetGroup
		[]string{"dn"}, // A list attributes to retrieve
		nil,
	)

	searchResult, err := client.Search(searchRequest)
	if err != nil {
		log.Printf("[ERROR] Error while searching group : %s ", err)
		return fmt.Errorf("Error while searching group : %s", err)
	}
	if len(searchResult.Entries) == 0 {
		log.Println("[ERROR] Target Group was not found, or is empty.")
		d.SetId("")
		return nil
	}

	var readDns []string
	for _, distinguishedName := range d.Get("dns_to_add").(*schema.Set).List() {
		distinguishedNameToRead := fmt.Sprintf("%s", distinguishedName)
		log.Printf("[DEBUG] Checking if %s is a member.", distinguishedNameToRead)
		for _, entry := range searchResult.Entries {
			if strings.EqualFold(entry.DN, distinguishedNameToRead) {
				log.Printf("[DEBUG] found %s", distinguishedNameToRead)
				readDns = append(readDns, distinguishedNameToRead)
			}
		}
	}

	d.Set("dns_to_add", readDns)
	return nil
}

func resourceAddToGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ldap.Conn)
	targetGroup := strings.ToLower(d.Get("target_group").(string))
	log.Printf("[DEBUG] Searching for members of %s", targetGroup)
	splitTargetGroup := strings.SplitN(targetGroup, "dc", 2) // split target group
	searchRequest := ldap.NewSearchRequest(
		"dc"+splitTargetGroup[1], // Make BaseDN from dc elements of split target group
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(|(objectCategory=user)(objectCategory=group))(memberOf="+targetGroup+"))", // Find users and groups that are members of targetGroup
		[]string{"dn"}, // A list attributes to retrieve
		nil,
	)

	searchResult, err := client.Search(searchRequest)
	if err != nil {
		log.Printf("[ERROR] Error while searching group : %s ", err)
		return fmt.Errorf("Error while searching group : %s", err)
	}
	if len(searchResult.Entries) == 0 {
		log.Println("[ERROR] Target Group was not found, or is empty.")
		d.SetId("")
		return nil
	}

	dnsToSplit := d.Id()
	dnsToRemove := strings.Split(dnsToSplit, "|")

	for _, distinguishedNameToRemove := range dnsToRemove {
		log.Printf("Checking for %s in target group.", distinguishedNameToRemove)
		for _, entry := range searchResult.Entries {
			if strings.EqualFold(entry.DN, distinguishedNameToRemove) {
				log.Printf("[DEBUG] Removing %s from %s", distinguishedNameToRemove, targetGroup)
				// call the helper to do the work
				err := removeFromGroup(distinguishedNameToRemove, targetGroup, client)
				if err != nil {
					log.Printf("[ERROR] Error while removing %s from %s : %s", distinguishedNameToRemove, targetGroup, err)
					return fmt.Errorf("Error while removing %s from %s : %s", distinguishedNameToRemove, targetGroup, err)
				}
			}
		}
	}

	return nil
}
