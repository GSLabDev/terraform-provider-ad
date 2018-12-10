---
layout: "ad"
page_title: "Active Directory: ad_computer_to_ou"
sidebar_current: "docs-ad-resource-inventory-folder"
description: |-
  Provides a Active Directory computer resource to Organizational Unit. This can be used to create and delete computer from OU.
---

# ad\_computer\_to\_ou

Provides a Active Directory computer resource to Organizational Unit. This can be used to create and delete computer from OU of AD.

## Example Usage

```hcl
resource "ad_computer_to_ou" "bar" {
  ou_distinguished_name = "ou=SubOU,ou=MyOU,dc=terraform,dc=com"
  computer_name         = "sampleName"
}
```

## Argument Reference

The following arguments are supported:

* `ou_distinguished_name` - (Required) The distinguished name of the Organizational Unit of the Active Directory
* `computer_name` - (Required) The name of a Computer to be added to Active Directory