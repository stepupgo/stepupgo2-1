package db

import (
	"database/sql"
	"log"
)

var (
	DB *sql.DB
)

func createTableSQL() string {
	const sql = `
		CREATE TABLE IF NOT EXISTS purchased (
		lottery_id  TEXT NOT NULL,
		number 		TEXT NOT NULL,
		PRIMARY KEY(lottery_id, number)
	);
	`

	return sql
}

func initDB() error {
	sql := createTableSQL()
	if _, err := DB.Exec(sql); err != nil {
		return err
	}
	return nil
}

func Init() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Println(err)
	}

	DB = db

	if err := initDB(); err != nil {
		log.Println(err)
	}
}
