package main

import (
	"encoding/json"
	"net/http"
)

// いずれ外部APIのインターフェイスDIする
type handlers struct {
}

func (h handlers) listPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/available_lotteries")
		if err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		defer resp.Body.Close()

		var lotteries []*Lottery
		if err := json.NewDecoder(resp.Body).Decode(&lotteries); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		if err := listTmpl.Execute(w, lotteries); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	}
}
