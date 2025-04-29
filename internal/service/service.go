package service

import (
	"context"

	"targeting-engine/internal/models"
)

// Service defines the interface for the targeting service
type Service interface {
	// GetMatchingCampaigns returns campaigns that match the given delivery request
	GetMatchingCampaigns(ctx context.Context, req models.DeliveryRequest) ([]models.CampaignResponse, error)
}
