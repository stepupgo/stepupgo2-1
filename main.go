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
		log.Println("failed to open sqlite file")
		os.Exit(1)
	}

	if err := initDB(db); err != nil {
		log.Printf("failed to initialize db err: %#v", err)
		os.Exit(1)
	}

	h := &handler{
		db: db,
	}

	http.HandleFunc("/", h.Home)
	http.HandleFunc("/purchase_page", h.PurchasePage)
	http.HandleFunc("/purchase", h.Purchase)
	http.HandleFunc("/result", h.Result)

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
