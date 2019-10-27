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
	"log"
	"net"
	"net/http"
	"os"

	controller "github.com/stepupgo/stepupgo2-1/controller"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", controller.LotteriesTypeGet())
	http.HandleFunc("/purchase_page", controller.PurchasePage())
	http.HandleFunc("/purchase", controller.Purchase())
	http.HandleFunc("/result", controller.LotteryResult())

	addr := net.JoinHostPort("", port)
	log.Println("Server running...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}

}
