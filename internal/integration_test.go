package internal

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ip2country/internal/config"
	"ip2country/internal/ip2country/handler"
	storeImpl "ip2country/internal/ip2country/store"
	"ip2country/internal/router"
	"ip2country/pkg/store"
)

type mockStore struct{}

func (m *mockStore) GetInfoByIP(ip net.IP) (*store.SubnetInfo, error) {
	if ip.String() == "2.22.233.255" {
		return &store.SubnetInfo{Country: "United Kingdom", City: "London"}, nil
	}
	return nil, store.ErrNotFound
}

func TestIntegrationMockStore(t *testing.T) {
	// Initialize the mock store
	mock := &mockStore{}
	handler.SetStore(mock)
	cfg, _ := config.LoadConfig()
	// Create a new router
	r := router.NewRouter(cfg)

	// Create a test server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Send a request to the test server
	resp, err := http.Get(ts.URL + "/v1/find-country?ip=2.22.233.255")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the status code
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"country":"United Kingdom","city":"London"}`
	body := make([]byte, len(expected))
	_, err = resp.Body.Read(body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", string(body), expected)
	}
}

func TestIntegrationFileStore(t *testing.T) {
	fileStore := storeImpl.NewFileStore("../db/geolite2.zip")
	handler.SetStore(fileStore)

	cfg, _ := config.LoadConfig()

	// Create a new router
	r := router.NewRouter(cfg)

	// Create a test server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Send a request to the test server
	resp, err := http.Get(ts.URL + "/v1/find-country?ip=2.22.233.255")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the status code
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"country":"Israel","city":"Rosh Ha‘Ayin"}`
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	actual := strings.TrimSpace(string(body))
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}
