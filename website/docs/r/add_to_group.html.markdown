---
layout: "ad"
page_title: "Active Directory: ad_add_to_group"
sidebar_current: "docs-ad-resource-inventory-folder"
description: |-
  Adds distinguished names to members of an Active Directory group object.
---

# ad\_add\_to\_group

Adds distinguished names to members of an Active Directory group object. Can be used to add users or groups to a group.

## Example Usage

```hcl
resource "ad_group_to_ou" "admins" {
  ou_distinguished_name        = "OU=groups,DC=company,DC=com"
  group_name                   = "admins"
  description                  = "Managed by terraform."
  gid_number                   = 9001
}

resource "ad_add_to_group" "main" {
  target_group = "CN=admins,OU=groups,DC=company,DC=com"
  dns_to_add = [
    "CN=alice,OU=users,DC=company,DC=com",
    "CN=bob,OU=users,DC=company,DC=com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `target_group` - (Required) The distinguished name of the target group you're adding members to.
* `dns_to_add` - (Required) A list of distinguished names to add to target_group. Can be users or groups.

## Attributes Reference

The following attributes are exported:

* `id` - A concatenation of each distinguished name in the `dns_to_add` list, seperated by a pipe `|`. e.g. 'CN=alice,OU=users,DC=company,DC=com|CN=bob,OU=users,DC=company,DC=com'
