package postgres

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func New() (*Postgres, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &Postgres{DB: db}, nil
}
