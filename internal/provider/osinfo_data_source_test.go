// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOsInfoDataSource(t *testing.T) {
	expectedName := runtime.GOOS
	expectedArch := runtime.GOARCH
	expectedId := expectedName + "/" + expectedArch
	is_windows := expectedName == "windows"
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.localos_info.test", "id", expectedId),
		resource.TestCheckResourceAttr("data.localos_info.test", "name", expectedName),
		resource.TestCheckResourceAttr("data.localos_info.test", "arch", expectedArch),
		resource.TestCheckResourceAttr("data.localos_info.test", "is_windows", strconv.FormatBool(is_windows)),
	}

	for _, e := range os.Environ() {
		kv := strings.Split(e, "=")
		checks = append(checks, resource.TestCheckResourceAttr("data.localos_info.test", fmt.Sprintf("environment.%s", kv[0]), kv[1]))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `data "localos_info" "test" {}`,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}
