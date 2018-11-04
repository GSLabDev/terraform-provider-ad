# Configure the Active Directory Provider
provider "ad" {
  domain         = "${var.ad_server_domain}"
  user           = "${var.ad_server_user}"
  password       = "${var.ad_server_password}"
  ip             = "${var.ad_server_ip}"
}

# Get Attributes of user "Matt"
data "ad_users" "matt" {
    base_search_dn = "OU=Users,DC=MY,DC=DOMAIN"
    username_filter = "Matthew Hodgkins"
    attributes = [ "cn", "sAMAccountName"]
}

# Get Attributes of user "Bob"
data "ad_users" "bob" {
    base_search_dn = "OU=Users,DC=MY,DC=DOMAIN"
    username_filter = "Bob *"
    attributes = [ "cn", "sAMAccountName"]
}

# Output queried attributes for "Matt"
output "matt" {
  value = "${data.ad_users.matt.user}"
}

# Output queried attributes for "Bob"
output "bob" {
  value = "${data.ad_users.bob.user}"
}
