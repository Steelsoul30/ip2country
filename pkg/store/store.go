// Package store Description: This package contains the repositories for the ip2country service.
package store

type Store interface {
	// GetCountryByIP Description: This method returns the country details for the given IP address.
	GetCountryByIP(ip string) (string, error)
}
