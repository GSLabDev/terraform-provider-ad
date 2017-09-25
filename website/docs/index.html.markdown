---
layout: "ad"
page_title: "Provider: Active Directory"
sidebar_current: "docs-ad-index"
description: |-
  The Active Directory provider is used to interact with the resources supported by
  Active Directory. The provider needs to be configured with the proper credentials
  before it can be used.
---

# Active Directory Provider

The Active Directory provider is used to interact with the resources supported by
Active Directory.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

~> **NOTE:** The Active Directory Provider currently represents _initial support_
and therefore may undergo significant changes as the community improves it. This
provider at this time only supports adding Computer Resource

## Example Usage

```hcl
# Configure the Active Directory Provider
provider "ad" {
  domain         = "${var.ad_server_domain}"
  user           = "${var.ad_server_user}"
  password       = "${var.ad_server_password}"
  ip             = "${var.ad_server_ip}"
}
# Add computer to Active Directory
resource "ad_computer" "foo" {
  domain        = "${var.ad_domain}"
  computer_name = "terraformSample"
}
```

## Argument Reference

The following arguments are used to configure the Active Directory Provider:

* `user` - (Required) This is the username for Active Directory Server operations. Can also
  be specified with the `AD_USER` environment variable.
* `password` - (Required) This is the password for Active Directory API operations. Can
  also be specified with the `AD_PASSWORD` environment variable.
* `ip` - (Required) This is the Active Directory server ip for Active Directory
  operations. Can also be specified with the `AD_SERVER` environment
  variable.
* `domain` - (Required) This is the domain of the Active Directory Server.

## Acceptance Tests

The Active Directory provider's acceptance tests require the above provider
configuration fields to be set using the documented environment variables.

In addition, the following environment variables are used in tests, and must be
set to valid values for your Active Directory environment:

 * AD\_COMPUTER\_DOMAIN

Once all these variables are in place, the tests can be run like this:

```
make testacc TEST=./builtin/providers/ad
```
