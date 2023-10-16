---
page_title: "localos Provider"
description: "The localos provider gets information about the machine running terraform and makes it available as data sources."
---


# localos Provider

The localos provider provides information in the form of data sources about
the operating system and environment of the machine
on which you are running terraform.

In certain situations it can be useful to know if your configuration is running on Windows or not,
especically for storing locally created artifacts such as key pairs in the appropriate directories.


## Example Usage

```terraform
terraform {
  required_providers {
    environment = {
      source = "registry.terraform.io/fireflycons/localos"
    }
  }
}

provider "localos" {}
```