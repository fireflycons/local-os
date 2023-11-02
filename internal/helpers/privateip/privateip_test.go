package privateip

import (
	"testing"

	"github.com/fireflycons/terraform-provider-localos/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestPrimaryInteface(t *testing.T) {

	var primary = MustGetLocalIP4Interfaces(true).GetPrimary()

	require.NotNil(t, primary, "No primary NIC located")
	require.Regexp(t, helpers.IpRegex, primary.Ip)
	require.Regexp(t, helpers.NetworkCidrRegex, primary.Network)
}

func TestSecondaryInterfaces(t *testing.T) {

	for _, nic := range MustGetLocalIP4Interfaces(true).GetSecondaries() {
		require.Regexp(t, helpers.IpRegex, nic.Ip)
		require.Regexp(t, helpers.NetworkCidrRegex, nic.Network)
	}
}
