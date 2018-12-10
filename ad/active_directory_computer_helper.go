package ad

import ldap "gopkg.in/ldap.v2"

func addComputerToAD(computerName string, dnName string, adConn *ldap.Conn, desc string) error {
	addRequest := ldap.NewAddRequest(dnName)
	addRequest.Attribute("objectClass", []string{"computer"})
	addRequest.Attribute("sAMAccountName", []string{computerName})
	addRequest.Attribute("userAccountControl", []string{"4096"})
	if desc != "" {
		addRequest.Attribute("description", []string{desc})
	}
	err := adConn.Add(addRequest)
	if err != nil {
		return err
	}
	return nil
}

func deleteComputerFromAD(dnName string, adConn *ldap.Conn) error {
	delRequest := ldap.NewDelRequest(dnName, nil)
	err := adConn.Del(delRequest)
	if err != nil {
		return err
	}
	return nil
}
