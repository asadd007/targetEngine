package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"targeting-engine/internal/models"

	_ "github.com/lib/pq"
)

var ErrCampaignNotFound = errors.New("campaign not found")

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(ctx context.Context, uri string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to PostgreSQL: %v", err)
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("couldn't ping PostgreSQL: %v", err)
	}

	if err := createTables(ctx, db); err != nil {
		return nil, fmt.Errorf("couldn't create tables: %v", err)
	}

	return &PostgresRepository{db: db}, nil
}

func createTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS campaigns (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			image_url VARCHAR(255) NOT NULL,
			cta VARCHAR(255) NOT NULL,
			status VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create targeting rules table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS targeting_rules (
			campaign_id VARCHAR(255) NOT NULL,
			dimension_type VARCHAR(255) NOT NULL,
			rule_type VARCHAR(255) NOT NULL,
			values TEXT[] NOT NULL,
			PRIMARY KEY (campaign_id, dimension_type),
			FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
		)
	`)
	return err
}

func (r *PostgresRepository) GetCampaigns(ctx context.Context) ([]models.Campaign, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, image_url, cta, status
		FROM campaigns
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var c models.Campaign
		if err := rows.Scan(&c.ID, &c.Name, &c.ImageURL, &c.CTA, &c.Status); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, c)
	}

	return campaigns, rows.Err()
}

func (r *PostgresRepository) SaveCampaign(ctx context.Context, campaign models.Campaign) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO campaigns (id, name, image_url, cta, status)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET name = $2, image_url = $3, cta = $4, status = $5
	`, campaign.ID, campaign.Name, campaign.ImageURL, campaign.CTA, campaign.Status)
	return err
}

func (r *PostgresRepository) GetTargetingRules(ctx context.Context) ([]models.TargetingRule, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT campaign_id, dimension_type, rule_type, values
		FROM targeting_rules
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.TargetingRule
	for rows.Next() {
		var r models.TargetingRule
		if err := rows.Scan(&r.CampaignID, &r.DimensionType, &r.RuleType, &r.Values); err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}

	return rules, rows.Err()
}

func (r *PostgresRepository) SaveTargetingRule(ctx context.Context, rule models.TargetingRule) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO targeting_rules (campaign_id, dimension_type, rule_type, values)
		VALUES ($1, $2, $3, $4)
	`, rule.CampaignID, rule.DimensionType, rule.RuleType, rule.Values)
	return err
}

func (r *PostgresRepository) Close(ctx context.Context) error {
	return r.db.Close()
}

func (r *PostgresRepository) InitTestData(ctx context.Context) error {
	spotifyAd := models.Campaign{
		ID:       "spotify",
		Name:     "Spotify - Music for everyone",
		ImageURL: "https://somelink",
		CTA:      "Download",
		Status:   models.StatusActive,
	}

	duolingoAd := models.Campaign{
		ID:       "duolingo",
		Name:     "Duolingo: Best way to learn",
		ImageURL: "https://somelink2",
		CTA:      "Install",
		Status:   models.StatusActive,
	}

	subwaySurferAd := models.Campaign{
		ID:       "subwaysurfer",
		Name:     "Subway Surfer",
		ImageURL: "https://somelink3",
		CTA:      "Play",
		Status:   models.StatusActive,
	}

	if err := r.SaveCampaign(ctx, spotifyAd); err != nil {
		return err
	}
	if err := r.SaveCampaign(ctx, duolingoAd); err != nil {
		return err
	}
	if err := r.SaveCampaign(ctx, subwaySurferAd); err != nil {
		return err
	}

	spotifyCountryRule := models.TargetingRule{
		CampaignID:    "spotify",
		DimensionType: models.DimensionCountry,
		RuleType:      models.Include,
		Values:        []string{"US", "Canada"},
	}

	duolingoOSRule := models.TargetingRule{
		CampaignID:    "duolingo",
		DimensionType: models.DimensionOS,
		RuleType:      models.Include,
		Values:        []string{"Android", "iOS"},
	}

	duolingoCountryRule := models.TargetingRule{
		CampaignID:    "duolingo",
		DimensionType: models.DimensionCountry,
		RuleType:      models.Exclude,
		Values:        []string{"US"},
	}

	subwaysurferOSRule := models.TargetingRule{
		CampaignID:    "subwaysurfer",
		DimensionType: models.DimensionOS,
		RuleType:      models.Include,
		Values:        []string{"Android"},
	}

	subwaysurferAppRule := models.TargetingRule{
		CampaignID:    "subwaysurfer",
		DimensionType: models.DimensionApp,
		RuleType:      models.Include,
		Values:        []string{"com.gametion.ludokinggame"},
	}

	if err := r.SaveTargetingRule(ctx, spotifyCountryRule); err != nil {
		return err
	}
	if err := r.SaveTargetingRule(ctx, duolingoOSRule); err != nil {
		return err
	}
	if err := r.SaveTargetingRule(ctx, duolingoCountryRule); err != nil {
		return err
	}
	if err := r.SaveTargetingRule(ctx, subwaysurferOSRule); err != nil {
		return err
	}
	if err := r.SaveTargetingRule(ctx, subwaysurferAppRule); err != nil {
		return err
	}

	return nil
}
