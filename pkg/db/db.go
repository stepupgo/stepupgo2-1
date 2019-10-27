package db

import "database/sql"

var (
	DB *sql.DB
)

func initDB() error {
	const sql = `
		CREATE TABLE IF NOT EXISTS purchased (
		lottery_id  TEXT NOT NULL,
		number 		TEXT NOT NULL,
		PRIMARY KEY(lottery_id, number)
	);
	`
	if _, err := DB.Exec(sql); err != nil {
		return err
	}
	return nil
}

func Init() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}

	DB = db

	if err := initDB(); err != nil {
		panic(err)
	}
}
