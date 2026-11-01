package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"request-service/internal/service"
)

type RequestHandler struct {
	service *service.RequestService
}

func NewRequestHandler(service *service.RequestService) *RequestHandler {
	return &RequestHandler{service: service}
}

func (h *RequestHandler) Register(r gin.IRouter) {
	r.POST("/api/v1/requests", h.Create)
	r.GET("/api/v1/requests/:id", h.Get)
}

func (h *RequestHandler) Create(c *gin.Context) {
	var body struct {
		CustomerID int    `json:"customer_id"`
		Title      string `json:"title"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, "title is required")
		return
	}

	request, err := h.service.Create(c.Request.Context(), body.CustomerID, body.Title, body.Status)
	if errors.Is(err, service.ErrInvalidTitle) {
		respondError(c, http.StatusBadRequest, "title is required")
		return
	}
	if errors.Is(err, service.ErrInvalidStatus) {
		respondError(c, http.StatusBadRequest, "invalid status")
		return
	}
	if errors.Is(err, service.ErrCustomerNotFound) {
		respondError(c, http.StatusBadRequest, "customer not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusCreated, request)
}

func (h *RequestHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}

	request, err := h.service.Get(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		respondError(c, http.StatusNotFound, "request not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, request)
}
