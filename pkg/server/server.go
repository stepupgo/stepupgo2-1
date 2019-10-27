package server

import (
	"net"
	"net/http"
	"os"

	"github.com/stepupgo/stepupgo2-1/pkg/db"
	"github.com/stepupgo/stepupgo2-1/pkg/di"
	"github.com/stepupgo/stepupgo2-1/pkg/server/handler"
)

func routing() {
	http.HandleFunc("/", handler.TopPage())
	http.HandleFunc("/purchase_page", handler.PurchasePage())
	http.HandleFunc("/purchase", handler.Purchase())
	http.HandleFunc("/result", handler.Result())
}

func connection() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := net.JoinHostPort("", port)
	http.ListenAndServe(addr, nil)
}

func Run() {
	db.Init()
	di.Init()
	routing()
	connection()
}
