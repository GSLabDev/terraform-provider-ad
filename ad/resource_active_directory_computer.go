package ad

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	ldap "gopkg.in/ldap.v3"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceComputer() *schema.Resource {
	return &schema.Resource{
		Create: resourceADComputerCreate,
		Read:   resourceADComputerRead,
		Delete: resourceADComputerDelete,
		Schema: map[string]*schema.Schema{
			"computer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
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

func resourceADComputerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ldap.Conn)

	computerName := d.Get("computer_name").(string)
	domain := d.Get("domain").(string)
	description := d.Get("description").(string)
	var dnOfComputer string
	dnOfComputer += "cn=" + computerName + ",cn=Computers"
	domainArr := strings.Split(domain, ".")
	for _, item := range domainArr {
		dnOfComputer += ",dc=" + item
	}

	log.Printf("[DEBUG] Name of the DN is : %s ", dnOfComputer)
	log.Printf("[DEBUG] Adding the Computer to the AD : %s ", computerName)

	err := addComputerToAD(computerName, dnOfComputer, client, description)
	if err != nil {
		log.Printf("[ERROR] Error while adding a Computer to the AD : %s ", err)
		return fmt.Errorf("Error while adding a Computer to the AD %s", err)
	}
	log.Printf("[DEBUG] Computer Added to AD successfully: %s", computerName)
	d.SetId(domain + "/" + computerName)
	return nil
}

func resourceADComputerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ldap.Conn)

	computerName := d.Get("computer_name").(string)
	domain := d.Get("domain").(string)
	var dnOfComputer string
	domainArr := strings.Split(domain, ".")
	dnOfComputer = "dc=" + domainArr[0]
	for index, item := range domainArr {
		if index == 0 {
			continue
		}
		dnOfComputer += ",dc=" + item
	}
	log.Printf("[DEBUG] Name of the DN is : %s ", dnOfComputer)
	log.Printf("[DEBUG] Deleting the Computer from the AD : %s ", computerName)

	searchRequest := ldap.NewSearchRequest(
		dnOfComputer, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=Computer)(cn="+computerName+"))", // The filter to apply
		[]string{"dn", "cn"}, // A list attributes to retrieve
		nil,
	)

	sr, err := client.Search(searchRequest)
	if err != nil {
		log.Printf("[ERROR] Error while searching a Computer : %s ", err)
		return fmt.Errorf("Error while searching a Computer : %s", err)
	}
	fmt.Println("[ERROR] Found " + strconv.Itoa(len(sr.Entries)) + " Entries")
	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}
	if len(sr.Entries) == 0 {
		log.Println("[ERROR] Computer was not found")
		d.SetId("")
	}
	return nil
}

func resourceADComputerDelete(d *schema.ResourceData, meta interface{}) error {
	resourceADComputerRead(d, meta)
	if d.Id() == "" {
		log.Println("[ERROR] Cannot find Computer in the specified AD")
		return fmt.Errorf("[ERROR] Cannot find Computer in the specified AD")
	}
	client := meta.(*ldap.Conn)

	computerName := d.Get("computer_name").(string)
	domain := d.Get("domain").(string)
	var dnOfComputer string
	dnOfComputer += "cn=" + computerName + ",cn=Computers"
	domainArr := strings.Split(domain, ".")
	for _, item := range domainArr {
		dnOfComputer += ",dc=" + item
	}

	log.Printf("[DEBUG] Name of the DN is : %s ", dnOfComputer)
	log.Printf("[DEBUG] Deleting the Computer from the AD : %s ", computerName)

	err := deleteComputerFromAD(dnOfComputer, client)
	if err != nil {
		log.Printf("[ERROR] Error while Deleting a Computer from AD : %s ", err)
		return fmt.Errorf("Error while Deleting a Computer from AD %s", err)
	}
	log.Printf("[DEBUG] Computer deleted from AD successfully: %s", computerName)
	return nil
}
