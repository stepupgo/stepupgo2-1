package usecase

import (
	"database/sql"
	"stepupgo2-1/domain/model"
	"stepupgo2-1/domain/repository"
	"stepupgo2-1/infrastructure/persistence"
)

type LotteryUsecase struct{}

func (looteryUsecase LotteryUsecase) SelectByPrimaryKey(DB *sql.DB, lotteryID string) (*model.Lottery, error) {
	lottery, err := repository.LotteryRepository(persistence.LotteryPersistence{}).SelectByPrimaryKey(DB, lotteryID)
	if err != nil {
		return nil, err
	}
	return lottery, nil
}
