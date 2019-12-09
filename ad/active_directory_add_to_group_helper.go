package ad

import ldap "gopkg.in/ldap.v3"

func addToGroup(groupToAddName string, targetGroupName string, adConn *ldap.Conn) error {
	modifyRequest := ldap.NewModifyRequest(targetGroupName, nil)
	modifyRequest.Add("member", []string{groupToAddName})

	err := adConn.Modify(modifyRequest)
	if err != nil {
		return err
	}
	return nil
}

func removeFromGroup(groupToRemoveName string, targetGroupName string, adConn *ldap.Conn) error {
	modifyRequest := ldap.NewModifyRequest(targetGroupName, nil)
	modifyRequest.Delete("member", []string{groupToRemoveName})

	err := adConn.Modify(modifyRequest)
	if err != nil {
		return err
	}
	return nil
}
