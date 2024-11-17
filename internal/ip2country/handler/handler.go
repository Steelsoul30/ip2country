// Package handler Description: This package contains the handler for the ip2country service.
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"ip2country/internal/middleware"
	"ip2country/pkg/store"
)

var storeImpl store.Store

type response struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

func SetStore(s store.Store) {
	storeImpl = s
}

func FindCountryHandler(w http.ResponseWriter, r *http.Request) {
	ipStr := r.URL.Query().Get("ip")
	if ipStr == "" {
		middleware.WriteError(w, http.StatusBadRequest, "IP parameter is missing")
		return
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		slog.Error(fmt.Sprintf("Invalid IP address: %v", ipStr))
		middleware.WriteError(w, http.StatusBadRequest, "Invalid IP address")
		return
	}

	info, err := storeImpl.GetInfoByIP(ip)
	if err != nil && errors.Is(err, store.ErrNotFound) {
		middleware.WriteError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := response{
		Country: info.Country,
		City:    info.City,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
