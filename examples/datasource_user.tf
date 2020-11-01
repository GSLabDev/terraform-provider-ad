provider "ad" {
    domain   = var.domain
    ip       = var.ip
    user     = var.user
    password = var.password
}

# search for a user in a domain by sAMAccountName
data "ad_user" "find" {
    logon_name = "test"
    domain     = var.domain
}

resource "ad_add_to_group" "add" {
    dns_to_add = [
         data.ad_user.find.dn,
    ]
    target_group = "CN=test-group,OU=test-ou,DC=domain,DC=com"
}
