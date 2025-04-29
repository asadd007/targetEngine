package service

import (
	"context"
	"os"
	"testing"

	"targeting-engine/internal/models"
	"targeting-engine/internal/repository"
)

type MockRepository struct {
	campaigns []models.Campaign
	rules     []models.TargetingRule
}

func (m *MockRepository) GetCampaigns(ctx context.Context) ([]models.Campaign, error) {
	return m.campaigns, nil
}

func (m *MockRepository) GetCampaignByID(ctx context.Context, id string) (*models.Campaign, error) {
	for _, c := range m.campaigns {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, repository.ErrCampaignNotFound
}

func (m *MockRepository) SaveCampaign(ctx context.Context, campaign models.Campaign) error {
	m.campaigns = append(m.campaigns, campaign)
	return nil
}

func (m *MockRepository) GetTargetingRules(ctx context.Context) ([]models.TargetingRule, error) {
	return m.rules, nil
}

func (m *MockRepository) GetTargetingRulesByCampaignID(ctx context.Context, campaignID string) ([]models.TargetingRule, error) {
	var campaignRules []models.TargetingRule
	for _, r := range m.rules {
		if r.CampaignID == campaignID {
			campaignRules = append(campaignRules, r)
		}
	}
	return campaignRules, nil
}

func (m *MockRepository) SaveTargetingRule(ctx context.Context, rule models.TargetingRule) error {
	m.rules = append(m.rules, rule)
	return nil
}

func (m *MockRepository) Close(ctx context.Context) error {
	return nil
}

func TestGetMatchingCampaigns(t *testing.T) {
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

	service := NewTargetingService(repo)

	tests := []struct {
		name           string
		request        models.DeliveryRequest
		expectedCount  int
		expectedIDs    []string
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Missing app parameter",
			request: models.DeliveryRequest{
				Country: "US",
				OS:      "Android",
			},
			expectError: true,
		},
		{
			name: "User in US on Android",
			request: models.DeliveryRequest{
				App:     "com.example.app",
				Country: "US",
				OS:      "Android",
			},
			expectedCount: 1,
			expectedIDs:   []string{"spotify"},
		},
		{
			name: "User in Canada on iOS",
			request: models.DeliveryRequest{
				App:     "com.example.app",
				Country: "CA",
				OS:      "iOS",
			},
			expectedCount: 2,
			expectedIDs:   []string{"spotify", "duolingo"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			campaigns, err := service.GetMatchingCampaigns(ctx, tc.request)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(campaigns) != tc.expectedCount {
				t.Errorf("Expected %d campaigns but got %d", tc.expectedCount, len(campaigns))
				return
			}

			for i, id := range tc.expectedIDs {
				if i < len(campaigns) && campaigns[i].CID != id {
					t.Errorf("Expected campaign ID %s but got %s", id, campaigns[i].CID)
				}
			}
		})
	}
}
