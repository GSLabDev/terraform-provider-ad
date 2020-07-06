---
layout: "ad"
page_title: "Active Directory: ad_organizational_unit"
sidebar_current: "docs-ad-resource-inventory-folder"
description: |-
  Creates Organizational Unit in active directory.
---

# ad\_organizational\_unit

Creates Organizational Unit in active directory.

## Example Usage

```hcl
# Add Organizational Unit to Active Directory
resource "ad_organizational_unit" "test" {
  ou_name                 = "sample-ou"
  ou_distinguished_name   = "OU=groups,DC=company,DC=com"
  description             = "Managed by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `ou_name` - (Required) Name of organizational unit.
* `ou_distinguished_name` - (Required) The distinguished name of the Organizational Unit of the Active Directory to add the ou to.
* `description` - (Optional) Sets the description property of the resultant ou object.
