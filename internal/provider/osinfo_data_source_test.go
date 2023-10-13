// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"runtime"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOsInfoDataSource(t *testing.T) {
	expectedName := runtime.GOOS
	expectedArch := runtime.GOARCH
	expectedId := expectedName + "/" + expectedArch
	is_windows := expectedName == "windows"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccOsInfoDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.localos_info.test", "id", expectedId),
					resource.TestCheckResourceAttr("data.localos_info.test", "name", expectedName),
					resource.TestCheckResourceAttr("data.localos_info.test", "arch", expectedArch),
					resource.TestCheckResourceAttr("data.localos_info.test", "is_windows", strconv.FormatBool(is_windows)),
				),
			},
		},
	})
}

const testAccOsInfoDataSourceConfig = `
data "localos_info" "test" {
}
`
