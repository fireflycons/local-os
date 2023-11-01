package helpers

import "regexp"

var (
	IpRegex          = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
	HostCidrRegex    = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|/32$)){4})`)
	NetworkCidrRegex = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|/([0-9]|[1-2][0-9]|3[0-2])$)){4})`)
)
