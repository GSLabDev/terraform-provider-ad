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
  ou_name       = "eample-ou"
  domain        = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `ou_name` - (Required) Name of organizational unit.
* `domain` - (Required) Name of domain under which you want to place ou.
