package interfaces

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"stepupgo2-1/interfaces/handler"

	"github.com/julienschmidt/httprouter"
)

// IsLetter function to check string is aplhanumeric only
var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

// Run start server
func Run(port int) error {
	log.Printf("Server running at http://localhost:%d/", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), Routes())
}

// Routes returns the initialized router
func Routes() *httprouter.Router {
	router := httprouter.New()

	// Index Route
	router.GET("/", handler.Index)

	//Lottery Route
	router.GET("/lottery", handler.HandleLotteryGet)//idによって取ってくる
	router.GET("/available_lotteries", handler.HandleAvailableLotteriesGet)

	//purchase Route
	// router.GET("/purchase_page", handler.HandlePurchasePageGet)
	// router.POST("/purchase", Handler.HandlePurchase)//idによって取ってくる
	// router.GET("/result")

	return router
}
