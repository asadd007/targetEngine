package service

import (
	"context"
	"errors"
	"strings"

	"targeting-engine/internal/models"
	"targeting-engine/internal/repository"
)

var (
	ErrInvalidRequest = errors.New("invalid request: missing required parameters")
)

type TargetingService struct {
	repo repository.Repository
}

func NewTargetingService(repo repository.Repository) *TargetingService {
	return &TargetingService{
		repo: repo,
	}
}

func (s *TargetingService) GetMatchingCampaigns(ctx context.Context, req models.DeliveryRequest) ([]models.CampaignResponse, error) {
	if req.App == "" || req.OS == "" || req.Country == "" {
		return nil, ErrInvalidRequest
	}

	campaigns, err := s.repo.GetCampaigns(ctx)
	if err != nil {
		return nil, err
	}

	rules, err := s.repo.GetTargetingRules(ctx)
	if err != nil {
		return nil, err
	}

	rulesByCampaign := make(map[string]map[models.DimensionType]models.TargetingRule)

	for _, rule := range rules {
		if _, ok := rulesByCampaign[rule.CampaignID]; !ok {
			rulesByCampaign[rule.CampaignID] = make(map[models.DimensionType]models.TargetingRule)
		}
		rulesByCampaign[rule.CampaignID][rule.DimensionType] = rule
	}

	var matchingAds []models.CampaignResponse
	for _, campaign := range campaigns {
		// Skip inactive ads // but we have picked only actives
		if campaign.Status != models.StatusActive {
			continue
		}

		if s.campaignMatchesRules(campaign.ID, req, rulesByCampaign) {
			matchingAds = append(matchingAds, campaign.ToCampaignResponse())
		}
	}

	return matchingAds, nil
}

func (s *TargetingService) campaignMatchesRules(campaignID string, req models.DeliveryRequest, rulesByCampaign map[string]map[models.DimensionType]models.TargetingRule) bool {
	rules, exists := rulesByCampaign[campaignID]
	if !exists {
		return true
	}

	if appRule, exists := rules[models.DimensionApp]; exists {
		if !s.matchesDimensionRule(req.App, appRule) {
			return false
		}
	}

	if countryRule, exists := rules[models.DimensionCountry]; exists {
		if !s.matchesDimensionRule(req.Country, countryRule) {
			return false
		}
	}

	if osRule, exists := rules[models.DimensionOS]; exists {
		if !s.matchesDimensionRule(req.OS, osRule) {
			return false
		}
	}

	return true
}

func (s *TargetingService) matchesDimensionRule(value string, rule models.TargetingRule) bool {
	normalizedValue := strings.ToLower(value)

	valueInRules := false
	for _, ruleValue := range rule.Values {
		if strings.ToLower(ruleValue) == normalizedValue {
			valueInRules = true
			break
		}
	}

	if rule.RuleType == models.Include {
		return valueInRules
	} else {
		return !valueInRules
	}
}
