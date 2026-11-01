package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"request-service/internal/repository"
	"request-service/internal/service"
)

func NewRouter(db *sql.DB) *gin.Engine {
	customers := NewCustomerHandler(
		service.NewCustomerService(
			repository.NewCustomerRepository(db),
		),
	)
	statuses := NewStatusHandler(
		service.NewStatusService(
			repository.NewStatusRepository(db),
		),
	)
	documents := NewDocumentHandler(
		service.NewDocumentService(
			repository.NewDocumentRepository(db),
		),
	)
	requests := NewRequestHandler(
		service.NewRequestService(
			repository.NewRequestRepository(db),
		),
	)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	customers.Register(router)
	statuses.Register(router)
	documents.Register(router)
	requests.Register(router)
	return router
}
