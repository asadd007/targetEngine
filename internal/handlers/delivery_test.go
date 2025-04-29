package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"targeting-engine/internal/models"
	"targeting-engine/internal/repository"
	"targeting-engine/internal/service"
)

func TestServeHTTP(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=true to run.")
	}
	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		postgresURI = "postgres://postgres:postgres@localhost:5432/targeting_engine_test?sslmode=disable"
	}
	ctx := context.Background()
	repo, err := repository.NewPostgresRepository(ctx, postgresURI)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer repo.Close(ctx)
	if err := repo.InitTestData(ctx); err != nil {
		t.Fatalf("Failed to initialize test data: %v", err)
	}

	targetingService := service.NewTargetingService(repo)

	tests := []struct {
		name           string
		url            string
		method         string
		expectedStatus int
		expectJSON     bool
	}{
		{
			name:           "Method not allowed",
			url:            "/v1/delivery?app=com.example.app&country=US&os=Android",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectJSON:     false,
		},
		{
			name:           "Missing app parameter",
			url:            "/v1/delivery?country=US&os=Android",
			method:         http.MethodGet,
			expectedStatus: http.StatusBadRequest,
			expectJSON:     true,
		},
		{
			name:           "Missing OS parameter",
			url:            "/v1/delivery?app=com.example.app&country=US",
			method:         http.MethodGet,
			expectedStatus: http.StatusBadRequest,
			expectJSON:     true,
		},
		{
			name:           "Missing country parameter",
			url:            "/v1/delivery?app=com.example.app&os=Android",
			method:         http.MethodGet,
			expectedStatus: http.StatusBadRequest,
			expectJSON:     true,
		},
		{
			name:           "User in US on Android",
			url:            "/v1/delivery?app=com.example.app&country=US&os=Android",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectJSON:     true,
		},
		{
			name:           "User in Canada on iOS",
			url:            "/v1/delivery?app=com.example.app&country=CA&os=iOS",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectJSON:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewDeliveryHandler(targetingService)
			req, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d but got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectJSON {
				contentType := rr.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type application/json but got %s", contentType)
				}
			}
			if tc.expectedStatus == http.StatusOK {
				var campaigns []models.CampaignResponse
				if err := json.NewDecoder(rr.Body).Decode(&campaigns); err != nil {
					t.Errorf("Failed to decode response body: %v", err)
					return
				}

				if len(campaigns) == 0 {
					t.Error("Expected at least one campaign but got none")
				}
			}
		})
	}
}
