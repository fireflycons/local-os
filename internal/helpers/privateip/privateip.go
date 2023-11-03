package privateip

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/jackpal/gateway"
)

type NIC struct {
	Name      string
	Ip        string
	Network   string
	IsPrimary bool
}

//go:generate mockery --name LocalInterfaces
type LocalInterfaces interface {
	// ScanInteterfaces reads the network interfaces connected to
	// this machine and populates internal data strucutures.
	ScanInterfaces() error

	// GetPrimary returns the interface connected to the
	// default gateway, or nil if such could not be determined.
	GetPrimary() *NIC

	// GetSecondaries returns all other interfaces.
	GetSecondaries() []*NIC

	// GetFirst returns the first NIC found.
	// This is normally the primary, but if no primary then the
	// first seconday in the list; else nil if no interfaces
	GetFirst() *NIC

	// Returns a reason as to why there is no primary interface
	GetPrimaryAbsentReason() string
}

var _ LocalInterfaces = &LocalInterfacesImpl{}

type LocalInterfacesImpl struct {
	nics                []*NIC
	primaryAbsentReason error
}

func New() LocalInterfaces {
	return &LocalInterfacesImpl{
		nics:                nil,
		primaryAbsentReason: errors.New("method ScanInterfaces() has not been called"),
	}
}

func (i *LocalInterfacesImpl) GetPrimary() *NIC {

	if i.nics == nil {
		return nil
	}

	for _, nic := range i.nics {
		if nic.IsPrimary {
			return nic
		}
	}

	return nil
}

func (i *LocalInterfacesImpl) GetSecondaries() []*NIC {

	if i.nics == nil {
		return make([]*NIC, 0, 1)
	}

	result := make([]*NIC, 0, 4)

	for _, nic := range i.nics {
		if !nic.IsPrimary {
			result = append(result, nic)
		}
	}

	return result
}

func (i *LocalInterfacesImpl) GetFirst() *NIC {

	nic := i.GetPrimary()

	if nic != nil {
		return nic
	}

	for _, nic = range i.GetSecondaries() {
		return nic
	}

	return nil
}

func (i *LocalInterfacesImpl) GetPrimaryAbsentReason() string {
	if i.primaryAbsentReason != nil {
		return i.primaryAbsentReason.Error()
	}

	return ""
}

func (i *LocalInterfacesImpl) ScanInterfaces() error {
	res, err := GetLocalIP4Interfaces(true)

	if err != nil {
		return err
	}

	interfaces, ok := res.(*LocalInterfacesImpl)

	if !ok {
		return fmt.Errorf("unexpected type %T in ScanInterfaces", res)
	}

	i.nics = interfaces.nics
	i.primaryAbsentReason = interfaces.primaryAbsentReason

	return nil
}

func MustGetLocalIP4Interfaces(includeLinkLocal bool) LocalInterfaces {
	n, e := GetLocalIP4Interfaces(includeLinkLocal)
	if e != nil {
		panic(e)
	}
	return n
}

// Get all non-loopback interfaces for this host.
// Primary is defined as the interface that routes to default gateway.
func GetLocalIP4Interfaces(includeLinkLocal bool) (LocalInterfaces, error) {
	var (
		nic      net.Interface
		nics     []net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
		err      error
	)

	result := &LocalInterfacesImpl{
		nics: make([]*NIC, 0, 4),
	}

	// Will be nil if no interface has a default gateway,
	// or the gateway package doesn't support it
	primaryIP, err := gateway.DiscoverInterface()

	if err != nil {
		result.primaryAbsentReason = err
	}

	if nics, err = net.Interfaces(); err != nil {
		return nil, err
	}

	for _, nic = range nics {

		if addrs, err = nic.Addrs(); err != nil { // get addresses
			continue
		}

		for _, addr := range addrs { // get ipv4 address
			n, ok := addr.(*net.IPNet)
			if !ok {
				return nil, errors.New("unable to cast 'net.Addr' to '*net.IPNet'")
			}
			if ipv4Addr = n.IP.To4(); ipv4Addr != nil {
				if n.IP.IsLoopback() {
					// always ignore
					continue
				}

				if primaryIP.Equal(n.IP) {
					result.nics = append(result.nics, &NIC{
						Name:      nic.Name,
						IsPrimary: true,
						Ip:        ipv4Addr.String(),
						Network:   getNetworkForHost(n).String(),
					})
				} else if !n.IP.IsLinkLocalUnicast() || (includeLinkLocal && n.IP.IsLinkLocalUnicast()) {
					result.nics = append(result.nics, &NIC{
						Name:      nic.Name,
						IsPrimary: false,
						Ip:        ipv4Addr.String(),
						Network:   getNetworkForHost(n).String(),
					})
				}
			}
		}
	}

	return result, nil
}

// Get subnet CIDR for given host.
func getNetworkForHost(host *net.IPNet) (network *net.IPNet) {

	network = &net.IPNet{
		IP:   make(net.IP, 4),
		Mask: host.Mask,
	}

	// Subnet address is IP of host BITWISE-AND netmask.
	binary.BigEndian.PutUint32(
		network.IP,
		binary.BigEndian.Uint32(host.IP.To4())&binary.BigEndian.Uint32(host.Mask),
	)

	return
}
