# Terraform Active Directory Provider

This is the repository for the Terraform Active Directory Provider, which one can use
with Terraform to work with Active Directory.

[1]: https://www.vmware.com/products/vcenter-server.html
[2]: https://www.vmware.com/products/esxi-and-esx.html

Coverage is currently only limited to a one resource only computer, but in the coming months we are planning release coverage for most essential Active Directory workflows.
Watch this space!

For general information about Terraform, visit the [official website][3] and the
[GitHub project page][4].

[3]: https://terraform.io/
[4]: https://github.com/hashicorp/terraform

# Using the Provider

The current version of this provider requires Terraform v0.10.2 or higher to
run.

Note that you need to run `terraform init` to fetch the provider before
deploying. Read about the provider split and other changes to TF v0.10.0 in the
official release announcement found [here][4].

[4]: https://www.hashicorp.com/blog/hashicorp-terraform-0-10/

## Full Provider Documentation

The provider is useful in adding computers to Active Directory.
### Example
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
  description   = "terraform sample server"
}
# Add computer to Organizational Unit of Active Directory
resource "ad_computer_to_ou" "bar" {
  ou_distinguished_name        = "${var.ad_ou_dn}"
  computer_name                = "terraformOuSample"
  description                  = "terraform sample server to OU"
}
# Add group to Organizational Unit of Active Directory
resource "ad_group_to_ou" "baz" {
  ou_distinguished_name        = "${var.ad_ou_dn}"
  group_name                   = "terraformGroupSample"
  description                  = "terraform sample group to OU"
}
# Add User to Active Directory
resource "ad_user" "foo1"{
    domain = "domain"
    first_name = "firstname"
    last_name = "lastname"
    logon_name    =  "logonname"
    password = "password"
}
```

# Building The Provider

**NOTE:** Unless you are [developing][7] or require a pre-release bugfix or feature,
you will want to use the officially released version of the provider (see [the
section above][8]).

[7]: #developing-the-provider
[8]: #using-the-provider


## Cloning the Project

First, you will want to clone the repository to
`$GOPATH/src/github.com/terraform-providers/terraform-provider-ad`:

```sh
mkdir -p $GOPATH/src/github.com/terraform-providers
cd $GOPATH/src/github.com/terraform-providers
git clone git@github.com:terraform-providers/terraform-provider-ad
```

## Running the Build

After the clone has been completed, you can enter the provider directory and
build the provider.

```sh
cd $GOPATH/src/github.com/terraform-providers/terraform-provider-ad
make build
```

## Installing the Local Plugin

After the build is complete, copy the `terraform-provider-ad` binary into
the same path as your `terraform` binary, and re-run `terraform init`.

After this, your project-local `.terraform/plugins/ARCH/lock.json` (where `ARCH`
matches the architecture of your machine) file should contain a SHA256 sum that
matches the local plugin. Run `shasum -a 256` on the binary to verify the values
match.

# Developing the Provider

If you wish to work on the provider, you'll first need [Go][9] installed on your
machine (version 1.9+ is **required**). You'll also need to correctly setup a
[GOPATH][10], as well as adding `$GOPATH/bin` to your `$PATH`.

[9]: https://golang.org/
[10]: http://golang.org/doc/code.html#GOPATH

See [Building the Provider][11] for details on building the provider.

[11]: #building-the-provider

# Testing the Provider

**NOTE:** Testing the Active Directory provider is currently a complex operation as it
requires having a Active Directory Server to test against.

## Configuring Environment Variables

Most of the tests in this provider require a comprehensive list of environment
variables to run. See the individual `*_test.go` files in the
[`ad/`](ad/) directory for more details. The next section also
describes how you can manage a configuration file of the test environment
variables.

### Using the `.tf-ad-devrc.mk` file

The [`tf-ad-devrc.mk.example`](tf-ad-devrc.mk.example) file contains
an up-to-date list of environment variables required to run the acceptance
tests. Copy this to `$HOME/.tf-ad-devrc.mk` and change the permissions to
something more secure (ie: `chmod 600 $HOME/.tf-ad-devrc.mk`), and
configure the variables accordingly.

## Running the Acceptance Tests

After this is done, you can run the acceptance tests by running:

```sh
$ make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the
`TESTARGS` parameter containing the run mask as per below:

```sh
make testacc TESTARGS="-run=TestAccAdComputer_Basic"
```
OR
```sh
make testacc TESTARGS="-run=TestAccAdComputerToOU_Basic"
```

This following example would run all of the acceptance tests matching
`TestAccAdComputer_Basic` OR `TestAccAdComputerToOU_Basic`. Change this for the
specific tests you want to run.
