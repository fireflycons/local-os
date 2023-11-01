// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/fireflycons/terraform-provider-localos/internal/helpers/privateip"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var nics = privateip.GetLocalIP4Interfaces(true)
var primary = privateip.GetPrimary(nics)
var countSecondaries = func() int {
	cnt := 0
	for _, n := range nics {
		if !n.IsPrimary {
			cnt++
		}
	}

	return cnt
}

func TestAccPrivateIpDataSource(t *testing.T) {

	var expectedResourceId = primary.Ip + "_" + strings.ReplaceAll(primary.Network, "/", "_")

	var checks = []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "id", expectedResourceId),
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "primary.ip", primary.Ip),
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "primary.cidr", primary.Ip+"/32"),
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "primary.network", primary.Network),
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "primary.name", primary.Name),
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "secondaries.#", strconv.Itoa(countSecondaries())),
	}

	ind := 0
	for _, nic := range nics {
		if !nic.IsPrimary {
			checks = append(checks, resource.TestCheckResourceAttr("data.localos_private_ip.test", fmt.Sprintf("secondaries.%d.ip", ind), nic.Ip))
			checks = append(checks, resource.TestCheckResourceAttr("data.localos_private_ip.test", fmt.Sprintf("secondaries.%d.cidr", ind), nic.Ip+"/32"))
			checks = append(checks, resource.TestCheckResourceAttr("data.localos_private_ip.test", fmt.Sprintf("secondaries.%d.network", ind), nic.Network))
			checks = append(checks, resource.TestCheckResourceAttr("data.localos_private_ip.test", fmt.Sprintf("secondaries.%d.name", ind), nic.Name))
			ind++
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `data "localos_private_ip" "test" {}`,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}
