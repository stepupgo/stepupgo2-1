package handler

import (
	"encoding/json"
	// "github.com/stepupgo/stepupgo2-1/lottery"
	// "github.com/stepupgo/stepupgo2-1/templates"
	"../lottery"
	"../templates"
	"net/http"
)



type Handler struct{
}

func (h Handler)RootHandler(w http.ResponseWriter, r *http.Request){
	resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/available_lotteries")
	if err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	defer resp.Body.Close()

	var lotteries []*lottery.Lottery
	if err := json.NewDecoder(resp.Body).Decode(&lotteries); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	if err := templates.ListTmpl.Execute(w, lotteries); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
}

