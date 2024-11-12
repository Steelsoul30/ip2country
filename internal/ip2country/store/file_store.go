package store

import "errors"

type FileStore struct {
}

func NewFileStore() *FileStore {
	return &FileStore{}
}

func (r *FileStore) GetCountryByIP(ip string) (string, error) {
	return "", errors.New("Not implemented")
}
