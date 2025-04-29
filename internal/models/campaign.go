package models

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)

type Campaign struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	CTA      string `json:"cta"`
	Status   Status `json:"status"`
}

type CampaignResponse struct {
	CID string `json:"cid"`
	Img string `json:"img"`
	CTA string `json:"cta"`
}

func (c *Campaign) ToCampaignResponse() CampaignResponse {
	return CampaignResponse{
		CID: c.ID,
		Img: c.ImageURL,
		CTA: c.CTA,
	}
}
