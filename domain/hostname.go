package domain

import "net"

type hostname struct {
	name    string
	address net.IP
}

func NewHostname(name string, address net.IP) (*hostname) {
	d := new(hostname)
	d.name = name
	d.address = address

	return d
}