package domain

import "net"

type Hostname struct {
	name    string
	address net.IP
}

func NewHostname(name string, address net.IP) (*Hostname) {
	d := new(Hostname)
	d.name = name
	d.address = address

	return d
}