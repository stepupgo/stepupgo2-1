package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/stepupgo/stepupgo2-1/pkg/db"
	"github.com/stepupgo/stepupgo2-1/pkg/server/model"
	"github.com/stepupgo/stepupgo2-1/pkg/view"
)

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
		if err := view.PurchasePageTmpl.Execute(w, data); err != nil {
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
		if err := db.DB.QueryRow("SELECT COUNT(*) FROM purchased WHERE lottery_id = ?", l.ID).Scan(&count); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		for i := 1; i <= num; i++ {
			const sql = "INSERT INTO purchased(lottery_id, number) values (?,?)"
			format := fmt.Sprintf(`%%0%dd`, len(strconv.FormatInt(l.Num-1, 10)))
			n := fmt.Sprintf(format, count+i)
			if _, err := db.DB.Exec(sql, id, n); err != nil {
				const status = http.StatusInternalServerError
				http.Error(w, http.StatusText(status), status)
				return
			}
		}

		http.Redirect(w, r, "/purchase_page?id="+l.ID, http.StatusFound)
	}
}
