package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"ip2country/pkg/store"
)

type APIStore struct {
	host string
}

func NewAPIStore(host string) *APIStore {
	return &APIStore{host: host}
}

func (r *APIStore) GetInfoByIP(ip net.IP) (*store.SubnetInfo, error) {

	client := &http.Client{}
	host := fmt.Sprintf("%s=%s", r.host, ip.String())
	slog.Info(fmt.Sprintf("Requesting data from the API: %s", host))
	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "keycdn-tools:https://www.github.com")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get a valid response from the server")
	}

	var result struct {
		Status      string `json:"status"`
		Description string `json:"description"`
		Data        struct {
			Geo struct {
				CountryName string `json:"country_name"`
				City        string `json:"city"`
			} `json:"geo"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, errors.New("failed to get a valid response from the server")
	}

	return &store.SubnetInfo{
		Subnet:  ip.String(),
		Country: result.Data.Geo.CountryName,
		City:    result.Data.Geo.City,
	}, nil
}
