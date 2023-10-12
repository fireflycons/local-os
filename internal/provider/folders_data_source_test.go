// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFoldersDataSource(t *testing.T) {
	expectedHome := os.Getenv("HOME")
	expectedSSH := path.Join(expectedHome, ".ssh")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccFoldersDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.localos_folders.test", "home", expectedHome),
					resource.TestCheckResourceAttr("data.localos_folders.test", "ssh", expectedSSH),
				),
			},
		},
	})
}

const testAccFoldersDataSourceConfig = `
data "localos_folders" "test" {
}
`
