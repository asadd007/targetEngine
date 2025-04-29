package models

import (
	"github.com/lib/pq"
)

type RuleType string

const (
	Include RuleType = "INCLUDE"
	Exclude RuleType = "EXCLUDE"
)

type DimensionType string

const (
	DimensionApp     DimensionType = "APP"
	DimensionCountry DimensionType = "COUNTRY"
	DimensionOS      DimensionType = "OS"
)

type TargetingRule struct {
	CampaignID    string         `json:"campaign_id"`
	DimensionType DimensionType  `json:"dimension_type"`
	RuleType      RuleType       `json:"rule_type"`
	Values        pq.StringArray `json:"values"`
}

type DeliveryRequest struct {
	App     string `json:"app"`
	OS      string `json:"os"`
	Country string `json:"country"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
