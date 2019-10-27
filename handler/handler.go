package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stepupgo/stepupgo2-1/lottery"
	"github.com/stepupgo/stepupgo2-1/templates"
	"log"
	"net/http"
	"strconv"
)

const lotteryEndpoint = "https://lottery-dot-tenntenn-samples.appspot.com"

type Handler struct{
}

func (h Handler)RootHandler(w http.ResponseWriter, r *http.Request){
	resp, err := http.Get(lotteryEndpoint+"/available_lotteries")
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

func (h Handler)PurchasePageHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(lotteryEndpoint+"/lottery?id=" + r.FormValue("id"))
	if err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	defer resp.Body.Close()

	var l lottery.Lottery
	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	data := struct {
		lottery.Lottery
		Remain int64
	}{
		Lottery: l,
		Remain:  l.Num, // TODO: 残りを計算する
	}
	if err := templates.PurchasePageTmpl.Execute(w, data); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
}

func (h Handler)PurchaseHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Println(err)
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	if err := initDB(db); err != nil {
		log.Println(err)
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
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

	var l lottery.Lottery
	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM purchased WHERE lottery_id = ?", l.ID).Scan(&count); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	for i := 1; i <= num; i++ {
		const sql = "INSERT INTO purchased(lottery_id, number) values (?,?)"
		format := fmt.Sprintf(`%%0%dd`, len(strconv.FormatInt(l.Num-1, 10)))
		n := fmt.Sprintf(format, count+i)
		if _, err := db.Exec(sql, id, n); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	}

	http.Redirect(w, r, "/purchase_page?id="+l.ID, http.StatusFound)
}

func(h Handler) ResultHandler (w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Println(err)
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	if err := initDB(db); err != nil {
		log.Println(err)
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	resp1, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/result?id=" + r.FormValue("id"))
	if err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	defer resp1.Body.Close()

	var result lottery.Result
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

	var l lottery.Lottery
	if err := json.NewDecoder(resp2.Body).Decode(&l); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	type winner struct {
		Prize   *lottery.Prize
		Numbers []string
	}

	data := struct {
		lottery.Lottery
		Winners map[string]*winner
	}{
		Lottery: l,
		Winners: map[string]*winner{},
	}

	rows, err := db.Query("SELECT number FROM purchased WHERE lottery_id = ?", l.ID)
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

	if err := templates.ResultTmpl.Execute(w, data); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
}


func initDB(db *sql.DB) error {
	const sql = `
CREATE TABLE IF NOT EXISTS purchased (
	lottery_id  TEXT NOT NULL,
	number 		TEXT NOT NULL,
	PRIMARY KEY(lottery_id, number)
);
`
	if _, err := db.Exec(sql); err != nil {
		return err
	}
	return nil
}
