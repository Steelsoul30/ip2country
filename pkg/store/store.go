// Package store Description: This package contains the repositories for the ip2country service.
package store

import (
	"errors"
	"net"
)

var ErrNotFound = errors.New("not found")

type Store interface {
	// GetCountryByIP Description: This method returns the country details for the given IP address.
	GetInfoByIP(ip net.IP) (*SubnetInfo, error)
}

// SubnetInfo holds information about each subnet
type SubnetInfo struct {
	Subnet  string // CIDR notation of the subnet
	Country string // Country name associated with the subnet
	City    string // City name associated with the subnet
}

type CustomTreeEntry struct {
	IPNet net.IPNet
	Info  SubnetInfo
}

func (e *CustomTreeEntry) Network() net.IPNet {
	return e.IPNet
}
