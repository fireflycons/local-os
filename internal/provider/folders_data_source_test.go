// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFoldersDataSource(t *testing.T) {
	var expectedHome, expectedSSH string

	if runtime.GOOS == "windows" {
		expectedHome = os.Getenv("USERPROFILE")
	} else {
		expectedHome = os.Getenv("HOME")
	}
	expectedSSH = filepath.Join(expectedHome, ".ssh")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `data "localos_folders" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.localos_folders.test", "home", expectedHome),
					resource.TestCheckResourceAttr("data.localos_folders.test", "ssh", expectedSSH),
				),
			},
		},
	})
}
