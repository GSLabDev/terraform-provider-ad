package ad

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	ldap "gopkg.in/ldap.v3"
)

func addGroupToAD(groupName string, dnName string, adConn *ldap.Conn, desc string, gidNumber string) error {
	addRequest := ldap.NewAddRequest(dnName, nil)
	addRequest.Attribute("objectClass", []string{"group"})
	addRequest.Attribute("sAMAccountName", []string{groupName})
	if desc != "" {
		addRequest.Attribute("description", []string{desc})
	}
	if gidNumber != "" {
		addRequest.Attribute("gidNumber", []string{gidNumber})
	}
	err := adConn.Add(addRequest)
	if err != nil {
		return err
	}
	return nil
}

func deleteGroupFromAD(dnName string, adConn *ldap.Conn) error {
	delRequest := ldap.NewDelRequest(dnName, nil)
	err := adConn.Del(delRequest)
	if err != nil {
		return err
	}
	return nil
}

func find_next_gidNumber(dnName string, adConn *ldap.Conn, auto_gid_min int, auto_gid_max int) (error, int) {
	log.Printf("[DEBUG] Searching next available gidNumber between %d and %d for %s.\n", auto_gid_min, auto_gid_max, dnName)
	baseDN := strings.Split(dnName, ",")
	searchRequest := ldap.NewSearchRequest(
		baseDN[len(baseDN)-2]+","+baseDN[len(baseDN)-1], // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=group)(gidNumber>="+strconv.Itoa(auto_gid_min)+")(gidNumber<="+strconv.Itoa(auto_gid_max)+"))", // The filter to apply
		[]string{"dn", "gidNumber"}, // A list attributes to retrieve
		nil,
	)
	searchResult, err := adConn.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
		return err, 0
	}

	var used_gid_numbers []int
	for _, entry := range searchResult.Entries {
		gid_number, _ := strconv.Atoi(entry.GetAttributeValues("gidNumber")[0])
		used_gid_numbers = append(used_gid_numbers, gid_number)
	}

	sort.Ints(used_gid_numbers)

	var next_available_gid int
	for potential_gid := auto_gid_min; potential_gid <= auto_gid_max; potential_gid++ {
		var potential_gid_used bool = false
		for _, used_gid := range used_gid_numbers {
			if used_gid == potential_gid {
				potential_gid_used = true
				break
			}
		}
		if potential_gid_used != true {
			next_available_gid = potential_gid
			break
		}
	}

	if next_available_gid == 0 {
		return fmt.Errorf("No available gidNumbers remaining."), 0
	}

	return nil, next_available_gid
}

func find_duplicate_gidNumber(dnName string, adConn *ldap.Conn, gidNumber int, auto_gid_min int, auto_gid_max int) (error, bool) {
	log.Printf("[DEBUG] Searching for duplicate groups with gidNumber %d.", gidNumber)
	baseDN := strings.Split(dnName, ",")
	searchRequest := ldap.NewSearchRequest(
		baseDN[len(baseDN)-2]+","+baseDN[len(baseDN)-1], // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=group)(gidNumber>="+strconv.Itoa(auto_gid_min)+")(gidNumber<="+strconv.Itoa(auto_gid_max)+"))", // The filter to apply
		[]string{"dn", "gidNumber"}, // A list attributes to retrieve
		nil,
	)
	searchResult, err := adConn.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
		return err, false
	}

	var used_gid_numbers []int
	for _, entry := range searchResult.Entries {
		gid_number, _ := strconv.Atoi(entry.GetAttributeValues("gidNumber")[0])
		used_gid_numbers = append(used_gid_numbers, gid_number)
	}

	var gidNumber_duplicate bool = false
	var gidNumber_times_used int = 0
	for _, used_gid := range used_gid_numbers {
		if used_gid == gidNumber {
			gidNumber_times_used++
			if gidNumber_times_used > 1 {
				gidNumber_duplicate = true
				break
			}
		}
	}

	return err, gidNumber_duplicate
}

func update_gidNumber(dnName string, adConn *ldap.Conn, gidNumber int) error {
	log.Printf("[DEBUG] Setting gidNumber %d on %s", gidNumber, dnName)
	modifyRequest := ldap.NewModifyRequest(dnName, nil)
	modifyRequest.Replace("gidNumber", []string{strconv.Itoa(gidNumber)})

	err := adConn.Modify(modifyRequest)
	if err != nil {
		return err
	}
	return nil
}
