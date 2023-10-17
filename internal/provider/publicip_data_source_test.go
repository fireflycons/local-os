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

	ip_regex := regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
	cidr_regex := regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|/32$)){4})`)
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
					resource.TestCheckResourceAttrWith("data.localos_public_ip.test", "ip", func(value string) error {
						if !ip_regex.MatchString(value) {
							return fmt.Errorf("Value %s does not match an IP address", value)
						}
						return nil
					}),
				),
			},
		},
	})
}

const testAccPublicIpDataSourceConfig = `data "localos_public_ip" "test" {}`
