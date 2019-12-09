---
layout: "ad"
page_title: "Active Directory: ad_group_to_ou"
sidebar_current: "docs-ad-resource-inventory-folder"
description: |-
  Creates a group object in an Active Directory Organizational Unit.
---

# ad\_group\_to\_ou

Creates a group object in an Active Directory Organizational Unit.

## Example Usage

```hcl
resource "ad_group_to_ou" "admins" {
  ou_distinguished_name        = "OU=groups,DC=company,DC=com"
  group_name                   = "admins"
  description                  = "Managed by terraform."
  gid_number                   = 9001
}
```

## Advanced auto_gid usage

```hcl
resource "ad_group_to_ou" "main" {
  count = 10
  ou_distinguished_name        = "OU=groups,DC=company,DC=com"
  group_name                   = "sample_group${count.index + 1}"
  description                  = "Managed by terraform."
  auto_gid                     = true
  auto_gid_min                 = 9001
  auto_gid_max                 = 9010
}
```

## Argument Reference

The following arguments are supported:

* `ou_distinguished_name` - (Required) The distinguished name of the Organizational Unit of the Active Directory to add the group to.
* `group_name` - (Required) The name of the group to be added.
* `description` - (Optional) Sets the description property of the resultant group object.
* `gid_number` - (Optional) Statically sets the 'gidNumber' attribute on the resultant group for use by Linux systems.
* `auto_gid` - (Optional) Boolean to automatically set the 'gidNumber' attribute on the resultant group, and ensure that it's unique. Does nothing when `gid_number` is set.
* `auto_gid_min` - (Optional) The lower bounds of automatically assignable gid numbers. Does nothing when `auto_gid` is not set, or set to 'false'.
* `auto_gid_max` - (Optional) The upper bounds of automatically assignable gid numbers. Does nothing when `auto_gid` is not set, or set to 'false'.

## Attributes Reference

The following attributes are exported:

* `id` - A concatenation of `group_name`/`ou_distinguished_name`.
* `auto_gid_number` - The 'gidNumber' that was set with `auto_gid`.
