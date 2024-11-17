package store

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"ip2country/internal/config"
	store2 "ip2country/pkg/store"
)

func NewStore(cfg *config.Config, cmd *cobra.Command) (store2.Store, error) {
	switch cfg.ActiveDataStore {
	case config.API:
		slog.Info("Using API data store")
		// Initialize and return the API store implementation
		for i, db := range cfg.DB {
			if db.Name == config.API {
				return NewAPIStore(cfg.DB[i].Host), nil
			}
		}

	case config.Relational:
		slog.Info("Using some relational database data store")
		// Initialize and return the relational database store implementation
		return NewDBStore(), fmt.Errorf("relational database data store not implemented yet")

	case config.Local:
		slog.Info("Using local data store. This might take a while to load.")
		path, _ := cmd.Flags().GetString("zippath")
		return NewFileStore(path), nil

	default:
		return nil, fmt.Errorf("unknown data store type: %s", cfg.ActiveDataStore)
	}
	return nil, fmt.Errorf("unknown data store type: %s", cfg.ActiveDataStore)
}
