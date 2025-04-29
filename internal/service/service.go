package service

import (
	"context"

	"targeting-engine/internal/models"
)

type Service interface {
	GetMatchingCampaigns(ctx context.Context, req models.DeliveryRequest) ([]models.CampaignResponse, error)
}
