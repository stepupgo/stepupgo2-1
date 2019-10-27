package api

import "net/http"

type top struct{}

type ITop interface {
	GetAvailableLotteries() (*http.Response, error)
}

func NewITop() ITop {
	return &top{}
}

func (t *top) GetAvailableLotteries() (*http.Response, error) {
	resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/available_lotteries")
	return resp, err
}
