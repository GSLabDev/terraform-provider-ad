package ad

import (
	"fmt"
	"log"

	"gopkg.in/ldap.v3"
)

type Config struct {
	Domain   string
	IP       string
	URL      string
	Username string
	Password string
}

// Client() returns a connection for accessing AD services.
func (c *Config) Client() (*ldap.Conn, error) {
	var username string
	var url string
	username = c.Username + "@" + c.Domain

	// stay downwards compatible
	switch {
	case c.URL != "":
		url = c.URL
	case c.IP != "":
		url = fmt.Sprintf("ldap://%s:389", c.IP)
	default:
		return nil, fmt.Errorf("Need either an IP or LDAP URL to connect to AD, check provider configuration")
	}

	adConn, err := clientConnect(url, username, c.Password)

	if err != nil {
		return nil, fmt.Errorf("Error while trying to connect active directory server, Check server IP address, username or password: %s", err)
	}
	log.Printf("[DEBUG] AD connection successful for user: %s", c.Username)
	return adConn, nil
}

func clientConnect(url, username, password string) (*ldap.Conn, error) {
	adConn, err := ldap.DialURL(url)
	if err != nil {
		return nil, err
	}

	err = adConn.Bind(username, password)
	if err != nil {
		return nil, err
	}
	return adConn, nil
}
