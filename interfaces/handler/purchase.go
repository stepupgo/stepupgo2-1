// package handler

// import (
// 	"net/http"
// 	"strconv"
// 	"log"

// 	"github.com/julienschmidt/httprouter"

// 	"stepupgo2-1/interfaces/response"
// )

// func HandlePurchase(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
// 	id := request.FormValue("id")
// 	num, err := strconv.Atoi(request.FormValue("num"))
// 	if err != nil {
// 		const status = http.StatusInternalServerError
// 		http.Error(writer, http.StatusText(status), status)
// 		return
// 	}
// 	// TODO: パラメタのバリデーション

// 	resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/lottery?id=" + id)
// 	if err != nil {
// 		log.Println(err)
// 		response.Error(writer, http.StatusInternalServerError, "Internal Server Error")
// 		return
// 	}
// 	defer resp.Body.Close()

// 	var l Lottery
// 	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
// 		log.Println(err)
// 		response.Error(writer, http.StatusInternalServerError, "Internal Server Error")
// 		return
// 	}

// 	//TODO: パラメータの分だけ購入する

// 	//TODO: 指定枚数分購入した結果をpurchasedテーブルにインサートする。


// 	http.Redirect(w, r, "/purchase_page?id="+l.ID, http.StatusFound)
// }

// func HandlePurchasePageGet(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
// 	resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/lottery?id=" + r.FormValue("id"))
// 	if err != nil {
// 		const status = http.StatusInternalServerError
// 		http.Error(w, http.StatusText(status), status)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	var l Lottery
// 	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
// 		const status = http.StatusInternalServerError
// 		http.Error(w, http.StatusText(status), status)
// 		return
// 	}

// 	data := struct {
// 		Lottery
// 		Remain int64
// 	}{
// 		Lottery: l,
// 		Remain:  l.Num, // TODO: 残りを計算する
// 	}
// 	if err := purchasePageTmpl.Execute(w, data); err != nil {
// 		const status = http.StatusInternalServerError
// 		http.Error(w, http.StatusText(status), status)
// 		return
// 	}
// }

// func HandleResult(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
// 	resp1, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/result?id=" + r.FormValue("id"))
// 	if err != nil {
// 		const status = http.StatusInternalServerError
// 		http.Error(w, http.StatusText(status), status)
// 		return
// 	}
// 	defer resp1.Body.Close()

// 	var result Result
// 	if err := json.NewDecoder(resp1.Body).Decode(&result); err != nil {
// 		const status = http.StatusInternalServerError
// 		http.Error(w, http.StatusText(status), status)
// 		return
// 	}

// 	resp2, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/lottery?id=" + r.FormValue("id"))
// 	if err != nil {
// 		const status = http.StatusInternalServerError
// 		http.Error(w, http.StatusText(status), status)
// 		return
// 	}
// 	defer resp2.Body.Close()

// 	var l Lottery
// 	if err := json.NewDecoder(resp2.Body).Decode(&l); err != nil {
// 		const status = http.StatusInternalServerError
// 		http.Error(w, http.StatusText(status), status)
// 		return
// 	}

// 	type winner struct {
// 		Prize   *Prize
// 		Numbers []string
// 	}

// 	data := struct {
// 		Lottery
// 		Winners map[string]*winner
// 	}{
// 		Lottery: l,
// 		Winners: map[string]*winner{},
// 	}

// 	rows, err := db.Query("SELECT number FROM purchased WHERE lottery_id = ?", l.ID)
// 	if err != nil {
// 		const status = http.StatusInternalServerError
// 		http.Error(w, http.StatusText(status), status)
// 		return
// 	}
// 	for rows.Next() {
// 		var number string
// 		if err := rows.Scan(&number); err != nil {
// 			const status = http.StatusInternalServerError
// 			http.Error(w, http.StatusText(status), status)
// 			return
// 		}

// 		for i := range result.Winners {
// 			for _, n := range result.Winners[i].Numbers {
// 				if number == n {
// 					prizeID := result.Winners[i].PrizeID
// 					if data.Winners[prizeID] == nil {
// 						for _, p := range l.Prizes {
// 							if p.ID == prizeID {
// 								data.Winners[prizeID] = &winner{
// 									Prize: p,
// 								}
// 							}
// 						}
// 					}
// 					data.Winners[prizeID].Numbers = append(data.Winners[prizeID].Numbers, n)
// 				}
// 			}
// 		}
// 	}

// 	if err := resultTmpl.Execute(w, data); err != nil {
// 		const status = http.StatusInternalServerError
// 		http.Error(w, http.StatusText(status), status)
// 		return
// 	}
// }

