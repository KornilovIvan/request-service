package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"request-service/internal/model"
	"request-service/internal/service"
)

type StatusHandler struct {
	service *service.StatusService
}

func NewStatusHandler(service *service.StatusService) *StatusHandler {
	return &StatusHandler{service: service}
}

func (h *StatusHandler) Register(r gin.IRouter) {
	r.GET("/api/v1/statuses", h.List)
}

func (h *StatusHandler) List(c *gin.Context) {
	statuses, err := h.service.List(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	if statuses == nil {
		statuses = []model.Status{}
	}
	c.JSON(http.StatusOK, statuses)
}
