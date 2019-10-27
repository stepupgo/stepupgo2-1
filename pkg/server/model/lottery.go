package model

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
