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

// Test with mock interfaces returning primary and secondaries.
func TestAccPrivateIpDataSource(t *testing.T) {
	mock := privateip.NewMockLocalInterfaces(t)
	primary := &privateip.NIC{
		Ip:        "172.31.0.1",
		Network:   "172.31.0.0/16",
		Name:      "test-primary",
		IsPrimary: true,
	}

	mock.On("ScanInterfaces").Return(nil)
	mock.On("GetPrimary").Return(primary)
	mock.On("GetSecondaries").Return([]*privateip.NIC{
		{
			Ip:        "192.168.0.1",
			Network:   "192.168.0.0/24",
			Name:      "test-secondary1",
			IsPrimary: false,
		},
		{
			Ip:        "169.254.0.1",
			Network:   "169.254.0.0/16",
			Name:      "test-secondary2",
			IsPrimary: false,
		},
	})
	mock.On("GetFirst").Return(primary)

	var testAccProtoV6MockProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"localos": providerserver.NewProtocol6WithError(newProviderWithMock("test", mock)()),
	}

	var expectedResourceId = mock.GetFirst().Ip + "_" + strings.ReplaceAll(mock.GetFirst().Network, "/", "_")

	var checks = []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "id", expectedResourceId),
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "primary.ip", primary.Ip),
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "primary.cidr", primary.Ip+"/32"),
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "primary.network", primary.Network),
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "primary.name", primary.Name),
		resource.TestCheckResourceAttr("data.localos_private_ip.test", "secondaries.#", strconv.Itoa(len(mock.GetSecondaries()))),
	}

	for ind, nic := range mock.GetSecondaries() {
		checks = append(checks, resource.TestCheckResourceAttr("data.localos_private_ip.test", fmt.Sprintf("secondaries.%d.ip", ind), nic.Ip))
		checks = append(checks, resource.TestCheckResourceAttr("data.localos_private_ip.test", fmt.Sprintf("secondaries.%d.cidr", ind), nic.Ip+"/32"))
		checks = append(checks, resource.TestCheckResourceAttr("data.localos_private_ip.test", fmt.Sprintf("secondaries.%d.network", ind), nic.Network))
		checks = append(checks, resource.TestCheckResourceAttr("data.localos_private_ip.test", fmt.Sprintf("secondaries.%d.name", ind), nic.Name))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6MockProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `data "localos_private_ip" "test" {}`,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

// Test using mocked LocalInterfaces to simulate no primary interface.
// BSD versions don't return a primary due to github.com/jackpal/gateway not implementing it.
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
		"localos": providerserver.NewProtocol6WithError(newProviderWithMock("test", mock)()),
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
		"localos": providerserver.NewProtocol6WithError(newProviderWithMock("test", mock)()),
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

func newProviderWithMock(version string, mock privateip.LocalInterfaces) func() provider.Provider {
	return func() provider.Provider {
		return &LocalOsProvider{
			version:         version,
			localInterfaces: mock,
		}
	}
}
