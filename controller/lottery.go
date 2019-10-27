package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	model "github.com/stepupgo/stepupgo2-1/model"
	repository "github.com/stepupgo/stepupgo2-1/model/db"
	temple "github.com/stepupgo/stepupgo2-1/view"
)

//宝くじの種類を取得
func LotteriesTypeGet() http.HandlerFunc {
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

		if err := temple.ListTmpl.Execute(w, lotteries); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		fmt.Println("compleat get lotteries")
	}
}

func PurchasePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/lottery?id=" + r.FormValue("id"))
		if err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		defer resp.Body.Close()

		var l model.Lottery
		if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		data := struct {
			model.Lottery
			Remain int64
		}{
			Lottery: l,
			Remain:  l.Num, // TODO: 残りを計算する
		}

		if err := temple.PurchasePageTmpl.Execute(w, data); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	}
}

func Purchase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		num, err := strconv.Atoi(r.FormValue("num"))
		if err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		// TODO: パラメタのバリデーション

		resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/lottery?id=" + id)
		if err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		defer resp.Body.Close()

		var l model.Lottery
		if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		var count int
		if err := repository.DB.QueryRow("SELECT COUNT(*) FROM purchased WHERE lottery_id = ?", l.ID).Scan(&count); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		for i := 1; i <= num; i++ {
			const sql = "INSERT INTO purchased(lottery_id, number) values (?,?)"
			format := fmt.Sprintf(`%%0%dd`, len(strconv.FormatInt(l.Num-1, 10)))
			n := fmt.Sprintf(format, count+i)
			if _, err := repository.DB.Exec(sql, id, n); err != nil {
				const status = http.StatusInternalServerError
				http.Error(w, http.StatusText(status), status)
				return
			}
		}

		http.Redirect(w, r, "/purchase_page?id="+l.ID, http.StatusFound)
	}
}

func LotteryResult() http.HandlerFunc {
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

		rows, err := repository.DB.Query("SELECT number FROM purchased WHERE lottery_id = ?", l.ID)
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

		if err := temple.ResultTmpl.Execute(w, data); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	}
}