package lib

import (
	"fmt"
	"time"
)

type DNServer struct {
	nameserver []string
	timeout    time.Duration
}

func New(ns []string, tout time.Duration) (dns DNServer, err error) {
	dns = DNServer{
		nameserver: ns,
		timeout:    tout,
	}
	if len(ns) == 0 {
		return dns, fmt.Errorf("nameserver is empty")
	}
	return
}

func (s DNServer) Dig(domain string) (IP []string, err error) {
	return
}

func (s DNServer) DigWithType(domain, typ string) (IP []string, err error) {
	return
}
