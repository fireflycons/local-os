package privateip

import (
	"encoding/binary"
	"net"

	"github.com/jackpal/gateway"
)

type NIC struct {
	Name      string
	Ip        string
	Network   string
	IsPrimary bool
}

func GetPrimary(nics []*NIC) *NIC {
	for _, nic := range nics {
		if nic.IsPrimary {
			return nic
		}
	}

	return nil
}

// Get all non-loopback interfaces for this host
// Primary is defined as the interface that routes to default gateway
func GetLocalIP4Interfaces(includeLinkLocal bool) []*NIC {
	var (
		nic      net.Interface
		nics     []net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
		err      error
	)

	results := make([]*NIC, 0, 4)
	// Will be nil if no interface has a default gateway
	// Terraform wouldn't be possible without it
	primaryIP, _ := gateway.DiscoverInterface()

	if nics, err = net.Interfaces(); err != nil {
		return nil
	}

	for _, nic = range nics {

		if addrs, err = nic.Addrs(); err != nil { // get addresses
			continue
		}

		for _, addr := range addrs { // get ipv4 address
			n := addr.(*net.IPNet)
			if ipv4Addr = n.IP.To4(); ipv4Addr != nil {
				if n.IP.IsLoopback() {
					// always ignore
					continue
				}

				if primaryIP.Equal(n.IP) {
					results = append(results, &NIC{
						Name:      nic.Name,
						IsPrimary: true,
						Ip:        ipv4Addr.String(),
						Network:   getNetworkForHost(n).String(),
					})
				} else if !n.IP.IsLinkLocalUnicast() || (includeLinkLocal && n.IP.IsLinkLocalUnicast()) {
					results = append(results, &NIC{
						Name:      nic.Name,
						IsPrimary: false,
						Ip:        ipv4Addr.String(),
						Network:   getNetworkForHost(n).String(),
					})
				}
			}
		}
	}

	return results
}

// Get subnet CIDR for given host
func getNetworkForHost(host *net.IPNet) (network *net.IPNet) {

	network = &net.IPNet{
		IP:   make(net.IP, 4),
		Mask: host.Mask,
	}

	// Subnet address is IP of host BITWISE-AND netmask
	binary.BigEndian.PutUint32(
		network.IP,
		binary.BigEndian.Uint32(host.IP.To4())&binary.BigEndian.Uint32(host.Mask),
	)

	return
}
