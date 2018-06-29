package lib

import "net"

type Domain struct {
	name    string
	address net.IP
}

func NewDomain(name string, address net.IP) (*Domain) {
	d := new(Domain)
	d.name = name
	d.address = address

	return d
}