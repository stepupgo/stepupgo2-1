package repository

import (
	"database/sql"

	"stepupgo2-1/domain/model"
)


type LotteryRepository interface {
	SelectByPrimaryKey(DB *sql.DB, lotteryID string) (*model.Lottery, error)
}
