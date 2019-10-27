package main

import (
	"net/http"
	"encoding/json"
)

// Handlers HTTPハンドラを集めた型
type Handlers struct {

}


// AvailableListHandler 公開されている宝くじ情報の一覧を取得
func (hs *Handlers) AvailableListHandler(w http.ResponseWriter, r *http.Request) {
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