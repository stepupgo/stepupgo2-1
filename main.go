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
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal("Cannot open database: ", err)
	}

	if err := initDB(db); err != nil {
		log.Fatal("Cannot init database", err)
	}

	h := &handler{}
	hdb := &dbHandler{db: db}

	http.HandleFunc("/", h.getRoot)

	http.HandleFunc("/purchase_page", h.getPurchasePage)

	http.HandleFunc("/purchase", hdb.getPurchase)

	http.HandleFunc("/result", hdb.getResult)

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
