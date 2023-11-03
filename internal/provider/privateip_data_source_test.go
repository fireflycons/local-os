// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/fireflycons/terraform-provider-localos/internal/helpers/privateip"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var nics = privateip.MustGetLocalIP4Interfaces(true)
var primary = nics.GetPrimary()
var countSecondaries = func() int {
	return len(nics.GetSecondaries())
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
	for _, nic := range nics.GetSecondaries() {
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

func NewWithMock(version string, mock privateip.LocalInterfaces) func() provider.Provider {
	return func() provider.Provider {
		return &LocalOsProvider{
			version:         version,
			localInterfaces: mock,
		}
	}
}

// Test using mocked LocalInterfaces to simulate no primary interface.
func TestAccPrivateIpDataSourceWithNoPrimary(t *testing.T) {
	mock := privateip.NewMockLocalInterfaces(t)
	secondary := &privateip.NIC{
		Ip:        "10.0.0.1",
		Network:   "10.0.0.0/8",
		Name:      "test",
		IsPrimary: false,
	}

	mock.On("ScanInterfaces").Return(nil)
	mock.On("GetPrimary").Return(nil)
	mock.On("GetPrimaryAbsentReason").Return("test set it to nil")
	mock.On("GetSecondaries").Return([]*privateip.NIC{
		secondary,
	})
	mock.On("GetFirst").Return(secondary)

	var testAccProtoV6MockProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"localos": providerserver.NewProtocol6WithError(NewWithMock("test", mock)()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6MockProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `data "localos_private_ip" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.localos_private_ip.test", "primary"),
				),
			},
		},
	})
}

// Test using mocked LocalInterfaces to simulate no interfaces at all.
func TestAccPrivateIpDataSourceWithNoInterfacesRaisesError(t *testing.T) {
	mock := privateip.NewMockLocalInterfaces(t)

	mock.On("ScanInterfaces").Return(nil)
	mock.On("GetPrimary").Return(nil)
	mock.On("GetPrimaryAbsentReason").Return("test set it to nil")
	mock.On("GetSecondaries").Return([]*privateip.NIC{})
	mock.On("GetFirst").Return(nil)

	var testAccProtoV6MockProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"localos": providerserver.NewProtocol6WithError(NewWithMock("test", mock)()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6MockProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      `data "localos_private_ip" "test" {}`,
				ExpectError: regexp.MustCompile(`No local NICs found`),
			},
		},
	})
}
