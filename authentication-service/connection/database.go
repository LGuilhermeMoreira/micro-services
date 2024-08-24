package connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectToDB() *sql.DB {
	dns := os.Getenv("DNS")
	var count int64
	for {
		db, err := openDB(dns)
		if err != nil {
			log.Println("Postgres not yet ready...")
		} else {
			log.Println("Connected to database")
			return db
		}
		if count > 10 {
			fmt.Println(err)
			return nil
		}
		count++
		log.Println("Wait 2 seconds and try again")
		time.Sleep(2 * time.Second)
		continue
	}
}
