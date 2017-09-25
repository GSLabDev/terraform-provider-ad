---
layout: "ad"
page_title: "Active Directory: ad_computer"
sidebar_current: "docs-ad-resource-inventory-folder"
description: |-
  Provides a Active Directory computer resource. This can be used to create and delete computer.
---

# ad\_computer

Provides a Active Directory computer resource. This can be used to create and delete computers from AD.

## Example Usage

```hcl
resource "ad_computer" "web" {
  domain        = "terraform.com"
  computer_name = "sampleName"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain of the Active Directory
* `computer_name` - (Required) The name of a Computer to be added to Active Directory