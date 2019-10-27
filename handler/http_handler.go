package handler

import (
    "encoding/json"
    "net/http"

    "github.com/stepupgo/stepupgo2-1/types"
    "github.com/stepupgo/stepupgo2-1/preview"
)

type Handler struct {
    
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
    resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/available_lotteries")
    if err != nil {
        const status = http.StatusInternalServerError
        http.Error(w, http.StatusText(status), status)
        return
    }
    defer resp.Body.Close()

    var lotteries []*types.Lottery
    if err := json.NewDecoder(resp.Body).Decode(&lotteries); err != nil {
        const status = http.StatusInternalServerError
        http.Error(w, http.StatusText(status), status)
        return
    }

    if err := preview.ListTmpl.Execute(w, lotteries); err != nil {
        const status = http.StatusInternalServerError
        http.Error(w, http.StatusText(status), status)
        return
    }
}


func (h *Handler) PurchasePageHandler(w http.ResponseWriter, r *http.Request) {
    resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/lottery?id=" + r.FormValue("id"))
    if err != nil {
        const status = http.StatusInternalServerError
        http.Error(w, http.StatusText(status), status)
        return
    }
    defer resp.Body.Close()

    var l types.Lottery
    if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
        const status = http.StatusInternalServerError
        http.Error(w, http.StatusText(status), status)
        return
    }

    data := struct {
        types.Lottery
        Remain int64
    }{
        Lottery: l,
        Remain:  l.Num, // TODO: 残りを計算する
    }
    if err := preview.PurchasePageTmpl.Execute(w, data); err != nil {
        const status = http.StatusInternalServerError
        http.Error(w, http.StatusText(status), status)
        return
    }
}

func (h *Handler) PurchaseHandler(w http.ResponseWriter, r *http.Request) {
    id := r.FormValue("id")
    num, err := strconv.Atoi(r.FormValue("num"))
    if err != nil {
        const status = http.StatusInternalServerError
        http.Error(w, http.StatusText(status), status)
        return
    }
    // TODO: パラメタのバリデーション

    resp, err := http.Get("https://lottery-dot-tenntenn-samples.appspot.com/lottery?id=" + id)
    if err != nil {
        const status = http.StatusInternalServerError
        http.Error(w, http.StatusText(status), status)
        return
    }
    defer resp.Body.Close()

    var l types.Lottery
    if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
        const status = http.StatusInternalServerError
        http.Error(w, http.StatusText(status), status)
        return
    }

    var count int
    if err := db.QueryRow("SELECT COUNT(*) FROM Purchased WHERE lottery_id = ?", l.ID).Scan(&count); err != nil {
        const status = http.StatusInternalServerError
        http.Error(w, http.StatusText(status), status)
        return
    }

    for i := 1; i <= num; i++ {
        const sql = "INSERT INTO Purchased(lottery_id, number) values (?,?)"
        format := fmt.Sprintf(`%%0%dd`, len(strconv.FormatInt(l.Num-1, 10)))
        n := fmt.Sprintf(format, count+i)
        if _, err := db.Exec(sql, id, n); err != nil {
            const status = http.StatusInternalServerError
            http.Error(w, http.StatusText(status), status)
            return
        }
    }

    http.Redirect(w, r, "/Purchase_page?id="+l.ID, http.StatusFound)
}
