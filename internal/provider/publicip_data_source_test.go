// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/fireflycons/terraform-provider-localos/internal/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccPublicIpDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `data "localos_public_ip" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.localos_public_ip.test", "id", amazonCheckIp),
					resource.TestCheckResourceAttrWith("data.localos_public_ip.test", "cidr", func(value string) error {
						if !assert.Regexp(t, helpers.HostCidrRegex, value) {
							return fmt.Errorf("Value %s does not match a /32 cidr", value)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("data.localos_public_ip.test", "ip", func(value string) error {
						if !assert.Regexp(t, helpers.IpRegex, value) {
							return fmt.Errorf("Value %s does not match an IP address", value)
						}
						return nil
					}),
				),
			},
		},
	})
}
