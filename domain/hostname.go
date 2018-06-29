package domain

import "net"

type Hostname struct {
	Name    string
	Address net.IP
}

func NewHostname(name string, address net.IP) (*Hostname) {
	d := new(Hostname)
	d.Name = name
	d.Address = address

	return d
}