package privateip

import (
	"testing"

	"github.com/fireflycons/terraform-provider-localos/internal/helpers"
	"github.com/stretchr/testify/require"
)

var nics = MustGetLocalIP4Interfaces(true)

func TestPrimaryInteface(t *testing.T) {

	var primary *NIC = nil

	for _, nic := range nics {
		if nic.IsPrimary {
			primary = nic
			break
		}
	}

	require.NotNil(t, primary, "No primary NIC located")
	require.Regexp(t, helpers.IpRegex, primary.Ip)
	require.Regexp(t, helpers.NetworkCidrRegex, primary.Network)
}

func TestSecondaryInterfaces(t *testing.T) {

	for _, nic := range nics {
		if !nic.IsPrimary {
			require.Regexp(t, helpers.IpRegex, nic.Ip)
			require.Regexp(t, helpers.NetworkCidrRegex, nic.Network)
		}
	}
}
