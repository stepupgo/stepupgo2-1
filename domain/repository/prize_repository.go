package repository

import (
	"database/sql"

	"stepupgo2-1/domain/model"
)

type PrizeRepository interface {
	SelectByPrimaryKey(DB *sql.DB, prizeID string) (*model.Prize, error)
}
