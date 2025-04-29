package repository

import (
	"context"
	"database/sql"
	"targeting-engine/internal/models"
)

type Repository interface {
	GetCampaigns(ctx context.Context) ([]models.Campaign, error)
	GetTargetingRules(ctx context.Context) ([]models.TargetingRule, error)
	Close(ctx context.Context) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetCampaigns(ctx context.Context) ([]models.Campaign, error) {
	query := `
		SELECT id, name, image_url, cta, status
		FROM campaigns
		WHERE status = 'ACTIVE'
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var campaign models.Campaign
		err := rows.Scan(&campaign.ID, &campaign.Name, &campaign.ImageURL, &campaign.CTA, &campaign.Status)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

func (r *repository) GetTargetingRules(ctx context.Context) ([]models.TargetingRule, error) {
	query := `
		SELECT campaign_id, dimension_type, rule_type, values
		FROM targeting_rules
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.TargetingRule
	for rows.Next() {
		var rule models.TargetingRule
		err := rows.Scan(&rule.CampaignID, &rule.DimensionType, &rule.RuleType, &rule.Values)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func (r *repository) Close(ctx context.Context) error {
	return r.db.Close()
}
