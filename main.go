package main

/*
TODO
	* テストコードを足す
	* main関数を分ける
	* 可読性
	* テスタビリティ
	* エラー処理（panicを使わない）
*/

import (
	"database/sql"
	"encoding/json"
	"net"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}

	if err := initDB(db); err != nil {
		panic(err)
	}

	h := &handler{}
	hdb := &dbHandler{db: db}

	http.HandleFunc("/", h.getRoot)

	http.HandleFunc("/purchase_page", h.getPurchasePage)

	http.HandleFunc("/purchase", hdb.getPurchase)

	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		resp1, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/result?id=" + r.FormValue("id"))
		if err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
		defer resp1.Body.Close()

		var result Result
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

		var l Lottery
		if err := json.NewDecoder(resp2.Body).Decode(&l); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}

		type winner struct {
			Prize   *Prize
			Numbers []string
		}

		data := struct {
			Lottery
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

		if err := resultTmpl.Execute(w, data); err != nil {
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
			return
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := net.JoinHostPort("", port)
	http.ListenAndServe(addr, nil)
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
