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
		log.Println(err)
		os.Exit(1)
	}
	if err := initDB(db); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	handler := Handler{
		db: db,
	}

	http.HandleFunc("/", handler.RouteHandle)
	http.HandleFunc("/purchase_page", handler.PurchasePageHandle)
	http.HandleFunc("/purchase", handler.PurchaseHandle)
	http.HandleFunc("/result", handler.ResultHandle)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := net.JoinHostPort("", port)
	http.ListenAndServe(addr, nil)
}
