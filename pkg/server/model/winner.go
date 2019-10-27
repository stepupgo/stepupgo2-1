package model

type Winner struct {
	PrizeID string   `json:"prize_id"`
	Numbers []string `json:"numbers"`
}
