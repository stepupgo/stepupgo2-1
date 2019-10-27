package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/stepupgo/stepupgo2-1/preview"
	"github.com/stepupgo/stepupgo2-1/types"

	_ "github.com/mattn/go-sqlite3"
)

type Handler struct {
	DB *sql.DB
}

// HomeHandler provides a handler for Home
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/available_lotteries")
	if err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	defer resp.Body.Close()

	var lotteries []*types.Lottery
	if err := json.NewDecoder(resp.Body).Decode(&lotteries); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	if err := preview.ListTmpl.Execute(w, lotteries); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
}

// PurchasePageHandler provides a handler for PurchasePage
func (h *Handler) PurchasePageHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/lottery?id=" + r.FormValue("id"))
	if err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	defer resp.Body.Close()

	var l types.Lottery
	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	data := struct {
		types.Lottery
		Remain int64
	}{
		Lottery: l,
		Remain:  l.Num, // TODO: 残りを計算する
	}
	if err := preview.PurchasePageTmpl.Execute(w, data); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
}

// PurchaseHandler provides a handler for purchasing
func (h *Handler) PurchaseHandler(w http.ResponseWriter, r *http.Request) {
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

	var l types.Lottery
	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	var count int
	if err := h.DB.QueryRow("SELECT COUNT(*) FROM Purchased WHERE lottery_id = ?", l.ID).Scan(&count); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	for i := 1; i <= num; i++ {
		const sql = "INSERT INTO Purchased(lottery_id, number) values (?,?)"
		format := fmt.Sprintf(`%%0%dd`, len(strconv.FormatInt(l.Num-1, 10)))
		n := fmt.Sprintf(format, count+i)
		if _, err := h.DB.Exec(sql, id, n); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	}

	http.Redirect(w, r, "/Purchase_page?id="+l.ID, http.StatusFound)
}

// ResultHandler provides a handler for getting result
func (h *Handler) ResultHandler(w http.ResponseWriter, r *http.Request) {
	resp1, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/result?id=" + r.FormValue("id"))
	if err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	defer resp1.Body.Close()

	var result types.Result
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

	var l types.Lottery
	if err := json.NewDecoder(resp2.Body).Decode(&l); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	type winner struct {
		Prize   *types.Prize
		Numbers []string
	}

	data := struct {
		types.Lottery
		Winners map[string]*winner
	}{
		Lottery: l,
		Winners: map[string]*winner{},
	}

	rows, err := h.DB.Query("SELECT number FROM Purchased WHERE lottery_id = ?", l.ID)
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

	if err := preview.ResultTmpl.Execute(w, data); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
}
