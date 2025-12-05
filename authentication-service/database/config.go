package database

import (
	"authentication/utils"
	"database/sql"
	"fmt"
	"log"
	"time"
)

var counts int64 = 0

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func New() *sql.DB {
	host := utils.GetEnv("POSTGRES_HOST", "localhost")
	port := utils.GetEnv("POSTGRES_PORT", "5432")
	user := utils.GetEnv("POSTGRES_USER", "postgres")
	password := utils.GetEnv("POSTGRES_PASSWORD", "password")
	dbName := utils.GetEnv("POSTGRES_DB", "users")
	sslMode := utils.GetEnv("DB_SSLMODE", "disable")
	timezone := utils.GetEnv("DB_TIMEZONE", "Asia/Ho_Chi_Minh")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbName, port, sslMode, timezone,
	)

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgress not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		// Break connect when retrying more than 10 times
		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
