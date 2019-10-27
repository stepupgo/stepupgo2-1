package model

type Result struct {
	Status  DrawStatus `json:"status"`
	Winners []*Winner  `json:"winners"`
}
