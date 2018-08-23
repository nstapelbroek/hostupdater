package traefik

import "net"

type Address struct {
	IP net.IP
	PortNumber int16
}