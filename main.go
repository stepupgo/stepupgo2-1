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

	"github.com/stepupgo/stepupgo2-1/handler"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Openに失敗したら
	// それ以降の処理はできないのでos.Exit(1)する
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Printf("failed to open : %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// dbのinitializeに失敗したら
	// それ以降の処理はできないのでos.Exit(1)する
	if err := initDB(db); err != nil {
		log.Printf("failed to initialize : %v\n", err)
		os.Exit(1)
	}

	v := &handler.Handler{DB: db}
	http.HandleFunc("/", v.HomeHandler)

	http.HandleFunc("/purchase_page", v.PurchasePageHandler)

	http.HandleFunc("/purchase", v.PurchaseHandler)

	http.HandleFunc("/result", v.ResultHandler)

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
