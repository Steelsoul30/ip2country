package store

import "errors"

type DBStore struct {
}

func NewDBStore() *DBStore {
	return &DBStore{}
}

func (r *DBStore) GetCountryByIP(ip string) (string, error) {
	return "", errors.New("Not implemented")
}
