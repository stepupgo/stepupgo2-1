package main

import "time"

type Lottery struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Price   int64     `json:"price"`
	Num     int64     `json:"num"`
	Prizes  []*Prize  `json:"prizes"`
	StartAt time.Time `json:"start_at"`
	DrawAt  time.Time `json:"draw_at"`
}

type Prize struct {
	ID     string `json:"id"`
	Name   string `json:"name`
	Num    int64  `json:"num"`
	Amount int64  `json:"amount"`
}

type Result struct {
	Status  DrawStatus `json:"status"`
	Winners []*Winner  `json:"winners"`
}

type DrawStatus int64

const (
	DrawStatusNotDrawn DrawStatus = 0
	DrawStatusMidDrawn DrawStatus = 1
	DrawStatusDrawn    DrawStatus = 2
)

type Winner struct {
	PrizeID string   `json:"prize_id"`
	Numbers []string `json:"numbers"`
}
