package handler

import (
	"encoding/json"
	"net/http"

	"github.com/stepupgo/stepupgo2-1/server/model"
	"github.com/stepupgo/stepupgo2-1/view"
)

func TopPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/available_lotteries")
		if err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		defer resp.Body.Close()

		var lotteries []*model.Lottery
		if err := json.NewDecoder(resp.Body).Decode(&lotteries); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		if err := view.ListTmpl.Execute(w, lotteries); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	}
}
