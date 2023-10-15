// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPublicIpDataSource(t *testing.T) {

	cidr_regex := regexp.MustCompile(`\d+\.\d+.\d+.\d+/32`)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccPublicIpDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.localos_public_ip.test", "id", amazonCheckIp),
					resource.TestCheckResourceAttrWith("data.localos_public_ip.test", "cidr", func(value string) error {
						if !cidr_regex.MatchString(value) {
							return fmt.Errorf("Value %s does not match a /32 cidr", value)
						}
						return nil
					}),
				),
			},
		},
	})
}

const testAccPublicIpDataSourceConfig = `data "localos_public_ip" "test" {}`
