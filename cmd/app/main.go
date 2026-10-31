package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"request-service/internal/handler"
	"request-service/internal/repository"
	"request-service/internal/service"
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

	customers := handler.NewCustomerHandler(
		service.NewCustomerService(
			repository.NewCustomerRepository(db),
		),
	)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	customers.Register(router)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
