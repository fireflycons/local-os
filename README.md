# Terraform Provider localos

The [localos provider](./docs/index.md) contains data sources that get information about the machine running `terraform apply`

The documentation for the provider can be found on the [Terraform Registry](https://registry.terraform.io/providers/fireflycons/localos/latest/docs)

## Example

```hcl
terraform {
  required_providers {
    localos = {
      source = "fireflycons/localos"
      version = "0.1.1"
    }
  }
}

provider "localos" {}

data "localos_folders" "folders" {}
data "localos_public_ip" "my_ip" {}

# Create a key pair
resource "tls_private_key" "pk" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "generated_key" {
  key_name   = "test-kp"
  public_key = tls_private_key.pk.public_key_openssh
}

# Save the PK to OS specific SSH keys folder
resource "local_sensitive_file" "private_key" {
  content  = tls_private_key.pk.private_key_pem
  filename = "${data.localos_folders.folders.ssh}/test-kp-pvt.pem"
}

# Create a security group that restricts access to only my public IP
resource "aws_security_group" "only_me" {
  name        = "only_me"
  description = "Allow all access only from my workstation"
  vpc_id      = var.vpc_id

  ingress {
    description = "Any from my IP"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [
      data.localos_public_ip.my_ip.cidr
    ]
  }
}
```


## Requirements to build/develop

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Using the provider

The data sources are

* [localos_info](./docs/data-sources/info.md) - Retrieves operating system (windows, linux etc), architecture (amd64, arm64 etc), and all environment variables.
* [localos_folders](./docs/data-sources/folders.md) - Gets paths to local folders of interest, currently user's home and ssh key directories.
* [localos_public_ip](./docs/data-sources/public_ip.md) - Gets the public IP of your workstation as an IP address and a /32 CIDR. Useful for configuring routes, firewalls etc for private infrastructure.


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
