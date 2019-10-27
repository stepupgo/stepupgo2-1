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
	_ "github.com/mattn/go-sqlite3"
	"github.com/stepupgo/stepupgo2-1/handler"
	"net"
	"net/http"
	"os"
)

func main() {

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}

	if err := initDB(db); err != nil {
		panic(err)
	}

	h := &handler.Handler{}
	hd := &handler.HandlerUseDB{DB: db}

	http.HandleFunc("/", h.LotteriesHandler)

	http.HandleFunc("/purchase_page", h.PurchasePageHandler)

	http.HandleFunc("/purchase", hd.PurchaseHandler)

	http.HandleFunc("/result", hd.ResultHandler)

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
