package helper

import "net"

func AddressToIp(address string) (ip net.IP, err error) {
	ip = net.ParseIP(address)
	if ip != nil {
		return
	}

	resolvedIp, err := net.LookupIP(address)
	if err != nil {
		return
	}

	return resolvedIp[0], err
}
