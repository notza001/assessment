package expenses

import (
	"database/sql"
	"log"
)

var Db *sql.DB

func open(url string) *sql.DB {
	// url := os.Getenv("DATABASE_URL")
	var err error
	Db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	return Db
}

func InitTable(dbUrl string) error {
	createTable := `
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			title TEXT,
			amount FLOAT,
			note TEXT,
			tags TEXT[]
		);
	`
	_, err := open(dbUrl).Exec(createTable)
	if err != nil {
		return err
	}
	return nil
}
