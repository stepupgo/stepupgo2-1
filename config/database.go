package config

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func init() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := initDB(db); err != nil {
		log.Fatal(err)
	}
}

func initDB(db *sql.DB) error {
	const sql = `
	CREATE TABLE IF NOT EXISTS purchased (
		lottery_id  TEXT NOT NULL,
		number 		TEXT NOT NULL,
		PRIMARY KEY(lottery_id, number)
	);
	`
	if _, err := db.Exec(sql); err != nil {
		return err
	}
	return nil
}
