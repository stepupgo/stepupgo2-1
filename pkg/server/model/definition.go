package model

type DrawStatus int64

const (
	DrawStatusNotDrawn DrawStatus = 0
	DrawStatusMidDrawn DrawStatus = 1
	DrawStatusDrawn    DrawStatus = 2
)
