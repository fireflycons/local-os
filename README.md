# Terraform Provider localos

The [localos provider](./docs/index.md) contains data sources that get information about the machine running `terraform apply`


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
