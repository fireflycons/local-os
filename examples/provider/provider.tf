terraform {
  required_providers {
    environment = {
      source = "registry.terraform.io/fireflycons/localos"
    }
  }
}

provider "localos" {}
