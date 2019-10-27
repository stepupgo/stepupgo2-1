package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"stepupgo2-1/application/usecase"
	"stepupgo2-1/config"
	"stepupgo2-1/interfaces/response"

	"github.com/julienschmidt/httprouter"
)

func HandleLotteryGet(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	id := request.FormValue("id")
	if id == "" {
		log.Println(errors.New("lottery_id is nil"))
		response.Error(writer, http.StatusBadRequest, "Bad Request")
		return
	}

	lotteryUsecase := usecase.LotteryUsecase{}
	var err error
	lotteryUsecase, err = usecase.LotteryUsecase{}.SelectByPrimaryKey(config.DB, id)
	if err != nil {
		log.Println(err)
		response.Error(writer, http.StatusInternalServerError, "Internal Server Error")
	}
	response.JSON(writer, http.StatusOK, lotteryUsecase)

}

func HandleAvailableLotteriesGet(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/available_lotteries")
	if err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	defer resp.Body.Close()

	var lotteries []*Lottery
	if err := json.NewDecoder(resp.Body).Decode(&lotteries); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	if err := listTmpl.Execute(w, lotteries); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
}
