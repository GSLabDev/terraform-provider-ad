package ad

import (
	"fmt"
	"log"

	"gopkg.in/ldap.v2"
)

type Config struct {
	Domain   string
	IP       string
	Username string
	Password string
}

// Client() returns a connection for accessing AD services.
func (c *Config) Client() (*ldap.Conn, error) {
	var username string
	username = c.Username + "@" + c.Domain
	adConn, err := clientConnect(c.IP, username, c.Password)

	if err != nil {
		return nil, fmt.Errorf("Error while trying to connect active directory server, Check server IP address, username or password: %s", err)
	}
	log.Printf("[DEBUG] AD connection successful for user: %s", c.Username)
	return adConn, nil
}

func clientConnect(ip, username, password string) (*ldap.Conn, error) {
	adConn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ip, 389))
	if err != nil {
		return nil, err
	}

	err = adConn.Bind(username, password)
	if err != nil {
		return nil, err
	}
	return adConn, nil
}
