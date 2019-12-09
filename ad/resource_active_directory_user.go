package ad

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/ldap.v2"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceADUserCreate,
		Read:   resourceADUserRead,
		Delete: resourceADUserDelete,

		Schema: map[string]*schema.Schema{
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logon_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
// function to add a user in AD:

func resourceADUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ldap.Conn)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	pass := d.Get("password").(string)
	domain := d.Get("domain").(string)
	logonName := d.Get("logon_name").(string)
	upn := logonName + "@" + domain
	userName := firstName + " " + lastName
	var dnOfUser string // dnOfUser: distingished names uniquely identifies an entry to AD.
	dnOfUser += "CN=" + userName + ",CN=Users"
	domainArr := strings.Split(domain, ".") 
	for _, i := range domainArr {
		dnOfUser += ",DC=" + i
	}

	log.Printf("[DEBUG] dnOfUser: %s ", dnOfUser)     
	log.Printf("[DEBUG] Adding user : %s ", userName) 
	err := addUser(userName, dnOfUser, client, upn, lastName, pass)
	if err != nil {
		log.Printf("[ERROR] Error while adding user: %s ", err) 
		return fmt.Errorf("Error while adding user %s", err)
	}
	log.Printf("[DEBUG] User Added success: %s", userName) 
	d.SetId(domain + "/" + userName)
	return nil
}

// Function to read user in AD:

func resourceADUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ldap.Conn)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	domain := d.Get("domain").(string)
	userName := firstName + " " + lastName
	var dnOfUser string // dnOfUser: distingished names uniquely identifies an entry to AD.
	domainArr := strings.Split(domain, ".")
	dnOfUser = "dc=" + domainArr[0]
	for index, i := range domainArr {
		if index == 0 {
			continue
		}
		dnOfUser += ",dc=" + i
	}
	log.Printf("[DEBUG] dnOfUser: %s ", dnOfUser)
	log.Printf("[DEBUG] Deleting user : %s ", userName)

	NewReq := ldap.NewSearchRequest(
		dnOfUser, // base dnOfUser.
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0,
		false,
		"(&(objectClass=User)(cn="+userName+"))", //applied filter
		[]string{"dnOfUser", "cn"},               
		nil,
	)

	sr, err := client.Search(NewReq)
	if err != nil {
		log.Printf("[ERROR] while seaching user : %s", err)
		return fmt.Errorf("Error while searching  user : %s", err)
	}

	fmt.Println("[ERROR] Found " + strconv.Itoa(len(sr.Entries)) + " Entries") 
	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn")) 

	}

	if len(sr.Entries) == 0 { 
		log.Println("[ERROR] user not found")
		d.SetId("") 
	}
	return nil
}

//function to delete user from AD :

func resourceADUserDelete(d *schema.ResourceData, m interface{}) error {
	resourceADUserRead(d, m)
	if d.Id() == "" {
		log.Printf("[ERROR] user not found !!")
		return fmt.Errorf("[ERROR] Cannot find user")
	}
	client := m.(*ldap.Conn)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	domain := d.Get("domain").(string)
	userName := firstName + " " + lastName
	var dnOfUser string 
	dnOfUser += "CN=" + userName + ",CN=Users"
	domainArr := strings.Split(domain, ".") 
	for _, i := range domainArr {
		dnOfUser += ",DC=" + i
	}
	log.Printf("[DEBUG] dnOfUser: %s ", dnOfUser)
	log.Printf("[DEBUG] deleting user : %s ", userName)
	err := delUser(userName, dnOfUser, client) 
	if err != nil {
		log.Printf("[ERROR] Error in deletion :%s", err)
		return fmt.Errorf("[ERROR] Error in deletion :%s", err)
	}
	log.Printf("[DEBUG] user Deleted success: %s", userName)
	return nil
}

// Helper function for adding user:
func addUser(userName string, dnName string, adConn *ldap.Conn, upn string, lastName string, pass string) error {
	a := ldap.NewAddRequest(dnName) // returns a new AddRequest without attributes " with dn".
	a.Attribute("objectClass", []string{"organizationalPerson", "person", "top", "user"})
	a.Attribute("sAMAccountName", []string{userName})
	a.Attribute("userPrincipalName", []string{upn})
	a.Attribute("name", []string{userName})
	a.Attribute("sn", []string{lastName})
	a.Attribute("userPassword", []string{pass})

	err := adConn.Add(a)
	if err != nil {
		return err
	}
	return nil
}
//Helper function to delete user:

func delUser(userName string, dnName string, adConn *ldap.Conn) error {
	delReq := ldap.NewDelRequest(dnName, nil)
	err := adConn.Del(delReq)
	if err != nil {
		return err
	}
	return nil
}
