package handler

import (
	"encoding/json"
	"net/http"

	"github.com/stepupgo/stepupgo2-1/db"
	"github.com/stepupgo/stepupgo2-1/server/model"
	"github.com/stepupgo/stepupgo2-1/view"
)

func Result() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp1, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/result?id=" + r.FormValue("id"))
		if err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		defer resp1.Body.Close()

		var result model.Result
		if err := json.NewDecoder(resp1.Body).Decode(&result); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		resp2, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/lottery?id=" + r.FormValue("id"))
		if err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		defer resp2.Body.Close()

		var l model.Lottery
		if err := json.NewDecoder(resp2.Body).Decode(&l); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		type winner struct {
			Prize   *model.Prize
			Numbers []string
		}

		data := struct {
			model.Lottery
			Winners map[string]*winner
		}{
			Lottery: l,
			Winners: map[string]*winner{},
		}

		rows, err := db.DB.Query("SELECT number FROM purchased WHERE lottery_id = ?", l.ID)
		if err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		for rows.Next() {
			var number string
			if err := rows.Scan(&number); err != nil {
				const status = http.StatusInternalServerError
				http.Error(w, http.StatusText(status), status)
				return
			}

			for i := range result.Winners {
				for _, n := range result.Winners[i].Numbers {
					if number == n {
						prizeID := result.Winners[i].PrizeID
						if data.Winners[prizeID] == nil {
							for _, p := range l.Prizes {
								if p.ID == prizeID {
									data.Winners[prizeID] = &winner{
										Prize: p,
									}
								}
							}
						}
						data.Winners[prizeID].Numbers = append(data.Winners[prizeID].Numbers, n)
					}
				}
			}
		}

		if err := view.ResultTmpl.Execute(w, data); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	}
}
