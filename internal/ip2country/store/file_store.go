package store

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"path/filepath"

	"github.com/yl2chen/cidranger"

	"ip2country/internal/dbgenerator"
	"ip2country/pkg/store"
)

type FileStore struct {
	tree cidranger.Ranger
}

func NewFileStore(zipPath string) *FileStore {

	generator := dbgenerator.NewDbGenerator()
	defer generator.Close()
	dataPath := filepath.Dir(zipPath) + "/geodata.dat"
	tree, err := generator.TryLoadFromGob(dataPath)
	if tree != nil {
		return &FileStore{tree: tree}
	}
	tree, err = generator.DirectFromZip(zipPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating file store: %v", err))
	}
	err = generator.SaveInfo(dataPath)
	if err != nil {
		slog.Error(fmt.Sprintf("Error saving file store: %v", err))
	}
	return &FileStore{tree: tree}

}

func (r *FileStore) GetInfoByIP(ip net.IP) (*store.SubnetInfo, error) {
	if r.tree == nil {
		return nil, errors.New("tree is nil")
	}

	if ip == nil {
		slog.Error(fmt.Sprintf("Invalid IP address: %v", ip))
		return nil, errors.New("invalid IP address")
	}

	entries, err := r.tree.ContainingNetworks(ip)
	if err != nil {
		slog.Error(fmt.Sprintf("Error finding IP: %v", err))
		return nil, err
	}
	if len(entries) == 0 {
		return nil, store.ErrNotFound
	}

	entry := entries[0].(*store.CustomTreeEntry)
	return &entry.Info, nil
}

func (r *FileStore) Close() {
	r.tree = nil

}
