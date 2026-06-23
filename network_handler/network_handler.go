package network_handler

import (
	"errors"
	"net"
)

type IPAddress = net.Addr

type Handler struct {
	Addr *net.IPAddr
}

func NewHandler(target string) (*Handler, error) {
	addr, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		return nil, errors.New("Unable to resolve IP Address")
	}

	return &Handler{
		Addr: addr,
	}, nil
}
