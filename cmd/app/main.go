package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"request-service/internal/handler"
)

func main() {
	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("DATABASE_URL is not set")
	}
	addr, ok := os.LookupEnv("HTTP_ADDR")
	if !ok {
		log.Fatal("HTTP_ADDR is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := handler.NewRouter(db).Run(addr); err != nil {
		log.Fatal(err)
	}
}
