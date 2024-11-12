package store

import "errors"

type APIStore struct {
}

func NewAPIStore() *APIStore {
	return &APIStore{}
}

func (r *APIStore) GetCountryByIP(ip string) (string, error) {
	return "", errors.New("not implemented")
}
