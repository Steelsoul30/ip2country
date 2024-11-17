package store

import (
	"errors"
	"net"

	"ip2country/pkg/store"
)

type DBStore struct {
}

func NewDBStore() *DBStore {
	return &DBStore{}
}

func (r *DBStore) GetInfoByIP(ip net.IP) (*store.SubnetInfo, error) {
	_ = ip
	return nil, errors.New("not implemented")
}
