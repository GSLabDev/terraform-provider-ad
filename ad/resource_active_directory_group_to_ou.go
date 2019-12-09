package ad

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

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
			"gid_number": {
				Type:        schema.TypeString,
				Description: "Statically sets the 'gidNumber' attribute on the resultant group.",
				Optional:    true,
				Default:     nil,
				ForceNew:    true,
			},
			"auto_gid": {
				Type:        schema.TypeBool,
				Description: "Automatically set the 'gidNumber' attribute on the resultant group, and ensure that it's unique.",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"auto_gid_min": {
				Type:        schema.TypeInt,
				Description: "The lower bounds of automatically assignable gid numbers.",
				Optional:    true,
				Default:     nil,
				ForceNew:    true,
			},
			"auto_gid_max": {
				Type:        schema.TypeInt,
				Description: "The upper bounds of automatically assignable gid numbers.",
				Optional:    true,
				Default:     nil,
				ForceNew:    true,
			},
			"auto_gid_number": {
				Type:        schema.TypeInt,
				Description: "The resultant gid number that was automatically set.",
				Computed:    true,
			},
		},
	}
}

func resourceADGroupToOUCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ldap.Conn)

	groupName := d.Get("group_name").(string)
	OUDistinguishedName := d.Get("ou_distinguished_name").(string)
	description := d.Get("description").(string)
	gidNumber := d.Get("gid_number").(string)
	auto_gid := d.Get("auto_gid").(bool)
	auto_gid_min := d.Get("auto_gid_min").(int)
	auto_gid_max := d.Get("auto_gid_max").(int)

	var dnOfGroup string
	dnOfGroup += "cn=" + groupName + "," + OUDistinguishedName

	log.Printf("[DEBUG] Name of the DN is : %s ", dnOfGroup)
	log.Printf("[DEBUG] Adding the Group to the AD : %s ", groupName)

	err := addGroupToAD(groupName, dnOfGroup, client, description, gidNumber)
	if err != nil {
		log.Printf("[ERROR] Error while adding a Group to the AD : %s ", err)
		return fmt.Errorf("Error while adding a Group to the AD %s", err)
	}

	// if gid_max is enabled and no gidNumber has been dictated
	if auto_gid == true && gidNumber == "" {
		// some sane defaults if no gid_min or gid_max have been chosen
		if auto_gid_min <= 0 {
			// probably don't want this either, but it's better than 0 which is reserved for root
			auto_gid_min = 1
		}
		if auto_gid_max <= 0 {
			// a "safe" gid max per resolution note 3:
			// https://access.redhat.com/solutions/25404
			auto_gid_max = 2097151
		}
		// check for smaller max than min
		if auto_gid_max < auto_gid_min {
			return fmt.Errorf("[ERROR] auto_gid_max must be greater than or equal to auto_gid_min.")
		}
		// random wait to help with gidNumber race
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(30) // maybe this wait should be variablized and pulled from an attrbiute
		time.Sleep(time.Duration(n) * time.Second)

		var duplicate_check bool = true
		for duplicate_check {
			// find next available gidNumber
			err, next_available_gid := find_next_gidNumber(dnOfGroup, client, auto_gid_min, auto_gid_max)
			if err != nil {
				log.Fatal(err)
				return fmt.Errorf("[ERROR] Error while searching for next available gidNumber. %s", err)
			}
			log.Printf("[DEBUG] Received %d as next available.", next_available_gid)
			// try updating the group with it
			err = update_gidNumber(dnOfGroup, client, next_available_gid)
			if err != nil {
				log.Fatal(err)
				return fmt.Errorf("[ERROR] Error while updating gidNumber of group. %s", err)
			}
			d.Set("auto_gid_number", next_available_gid)
			// wait a moment for things to stick
			time.Sleep(2 * time.Second) // maybe this wait should be variablized and pulled from an attrbiute
			// check for duplicates will return false and break the loop if no dups found
			err, duplicate_check = find_duplicate_gidNumber(dnOfGroup, client, next_available_gid, auto_gid_min, auto_gid_max)
			if err != nil {
				log.Fatal(err)
				return fmt.Errorf("[ERROR] Error while checking for duplicate gidNumbers %s", err)
			}
			// if we got a duplicate, wait, and to the top try again
			if duplicate_check {
				rand.Seed(time.Now().UnixNano())
				n := rand.Intn(10) // maybe this wait should be variablized and pulled from an attrbiute
				log.Printf("[DEBUG] Found duplicate gidNumber %d. Trying again in %d seconds.", next_available_gid, n)
				time.Sleep(time.Duration(n) * time.Second)
			}
			// maybe need some max retries code here
		}
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
