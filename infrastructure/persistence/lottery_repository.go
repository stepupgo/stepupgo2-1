package persistence

type LotteryPersistence struct{}

func (lotteryPersistence LotteryPersistence) SelectByPrimaryKey(DB *sql.DB, LotteryID string) (*model.Lottery, error) {
	row := DB.QueryRow("SELECT * FROM lottery WHERE lottery_id = ?", lotteryID)
	return convertToLottery(row)
}

// convertToLottery convert type *sql.Row to Lottery.
func convertToLottery(row *sql.Row) (*model.Lottery, error) {
	lottery := model.Lottery{}
	err := row.Scan(&lottery.LotteryID, &lottery.Name, &lottery.Price, &lottery.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return &lottery, nil
}