package store_test

import (
	"errors"
	"net"
	"path/filepath"
	"testing"

	sut "ip2country/internal/ip2country/store"
	"ip2country/pkg/store"
)

// Mock implementation for API store
type mockAPIStore struct {
	info store.SubnetInfo
	err  error
}

func (m *mockAPIStore) GetInfoByIP(ip net.IP) (*store.SubnetInfo, error) {
	return &m.info, m.err
}

func TestNewFileStore(t *testing.T) {
	// Test with a valid zip path
	a, _ := filepath.Abs("geolite2-test.zip")
	_ = a
	fs := sut.NewFileStore("geolite2-test.zip")
	if fs == nil {
		t.Fatal("Expected non-nil FileStore")
	}
}

func TestFileStore_GetInfoByIP(t *testing.T) {
	fs := sut.NewFileStore("geolite2-test.zip")

	tests := []struct {
		name          string
		ip            string
		expectedError error
	}{
		{
			name:          "Valid IP",
			ip:            "5.132.126.112",
			expectedError: nil,
		},
		{
			name:          "Invalid IP",
			ip:            "invalid-ip",
			expectedError: errors.New("invalid IP address"),
		},
		{
			name:          "IP not found",
			ip:            "223.111.211.2",
			expectedError: store.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			_, err := fs.GetInfoByIP(ip)
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestFileStore_Close(t *testing.T) {
	fs := sut.NewFileStore("geolite2-test.zip")
	fs.Close()
	// Attempt to get info by IP after closing the store
	ip := net.ParseIP("8.8.8.8")
	_, err := fs.GetInfoByIP(ip)
	if err == nil {
		t.Fatal("Expected error after closing the store, got nil")
	}
}

func TestAPIStore_GetInfoByIP(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		storeInfo store.SubnetInfo
		storeErr  error
	}{
		{
			name:      "Valid IP",
			ip:        "8.8.8.8",
			storeInfo: store.SubnetInfo{Country: "USA", City: "Mountain View"},
		},
		{
			name:     "Invalid IP",
			ip:       "invalid-ip",
			storeErr: store.ErrNotFound,
		},
		{
			name:     "Nil IP",
			ip:       "",
			storeErr: errors.New("invalid IP address"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &mockAPIStore{
				info: tt.storeInfo,
				err:  tt.storeErr,
			}

			ip := net.ParseIP(tt.ip)
			_, err := mockStore.GetInfoByIP(ip)
			if err != nil && err.Error() != tt.storeErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.storeErr, err)
			}
		})
	}
}
