package persistence

import (
	"database/sql"
	"log"

	"stepupgo2-1/domain/model"
)

type PrizePersistence struct{}

func (prizePersistence PrizePersistence) SelectByPrimaryKey(DB *sql.DB, prizeID string) (*model.Prize, error) {
	row := DB.QueryRow("SELECT * FROM lottery WHERE lottery_id = ?", prizeID)
	return convertToPrize(row)
}

// convertToPrize convert type *sql.Row to Prize.
func convertToPrize(row *sql.Row) (*model.Prize, error) {
	prize := model.Prize{}
	err := row.Scan(&prize.PrizeID, &prize.Name, &prize.Number, &prize.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return &prize, nil
}
