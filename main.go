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
	_ "github.com/mattn/go-sqlite3"
	"github.com/stepupgo/stepupgo2-1/handler"
	"net"
	"net/http"
	"os"
)


func main() {


	var h handler.Handler

	http.HandleFunc("/", h.RootHandler)

	http.HandleFunc("/purchase_page", h.PurchasePageHandler)

	http.HandleFunc("/purchase", h.PurchaseHandler)

	http.HandleFunc("/result", h.ResultHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := net.JoinHostPort("", port)
	http.ListenAndServe(addr, nil)
}

