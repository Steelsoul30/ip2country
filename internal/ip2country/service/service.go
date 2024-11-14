// Package service contains the ip2country core business logic
package service

import "ip2country/pkg/store"

type Ip2countryService struct {
	store store.Store
}

func NewIp2countryService(store store.Store) *Ip2countryService {
	return &Ip2countryService{store: store}
}

func (s *Ip2countryService) GetCountryByIP(ip string) (string, error) {
	return s.store.GetCountryByIP(ip)
}
