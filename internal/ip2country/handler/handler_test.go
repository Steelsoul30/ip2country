package handler_test

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"ip2country/internal/ip2country/handler"
	"ip2country/pkg/store"
)

type mockStore struct {
	info *store.SubnetInfo
	err  error
}

func (m *mockStore) GetInfoByIP(ip net.IP) (*store.SubnetInfo, error) {
	return m.info, m.err
}

func TestFindCountryHandler(t *testing.T) {
	tests := []struct {
		name           string
		ip             string
		storeInfo      *store.SubnetInfo
		storeErr       error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid IP",
			ip:             "8.8.8.8",
			storeInfo:      &store.SubnetInfo{Country: "USA", City: "Mountain View"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"USA","city":"Mountain View"}`,
		},
		{
			name:           "Missing IP parameter",
			ip:             "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"IP parameter is missing"}`,
		},
		{
			name:           "Invalid IP address",
			ip:             "invalid-ip",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid IP address"}`,
		},
		{
			name:           "IP not found",
			ip:             "8.8.8.8",
			storeErr:       store.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"not found"}`,
		},
		{
			name:           "Internal server error",
			ip:             "8.8.8.8",
			storeErr:       errors.New("internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &mockStore{
				info: tt.storeInfo,
				err:  tt.storeErr,
			}
			handler.SetStore(mockStore)

			req, _ := http.NewRequest("GET", "/v1/find-country?ip="+tt.ip, nil)
			rr := httptest.NewRecorder()

			handler.FindCountryHandler(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var responseBody map[string]string
			_ = json.Unmarshal(rr.Body.Bytes(), &responseBody)
			expectedBody := make(map[string]string)
			_ = json.Unmarshal([]byte(tt.expectedBody), &expectedBody)

			if !equal(responseBody, expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v", responseBody, expectedBody)
			}
		})
	}
}

func equal(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
