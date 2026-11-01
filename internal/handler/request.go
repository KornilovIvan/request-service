package handler

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"request-service/internal/model"
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
	r.GET("/api/v1/requests", h.List)
	r.GET("/api/v1/requests/:id", h.Get)
	r.PATCH("/api/v1/requests/:id/status", h.UpdateStatus)
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

func (h *RequestHandler) List(c *gin.Context) {
	query := c.Request.URL.Query()
	requests, err := h.service.List(
		c.Request.Context(),
		queryParam(query, "status"),
		queryParam(query, "created_from"),
		queryParam(query, "limit"),
		time.Now(),
	)
	if errors.Is(err, service.ErrInvalidStatus) {
		respondError(c, http.StatusBadRequest, "invalid status")
		return
	}
	if errors.Is(err, service.ErrInvalidDate) {
		respondError(c, http.StatusBadRequest, "invalid date")
		return
	}
	if errors.Is(err, service.ErrInvalidLimit) {
		respondError(c, http.StatusBadRequest, "invalid limit")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	if requests == nil {
		requests = []model.Request{}
	}
	c.JSON(http.StatusOK, requests)
}

func queryParam(query url.Values, key string) *string {
	if !query.Has(key) {
		return nil
	}
	value := query.Get(key)
	return &value
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

func (h *RequestHandler) UpdateStatus(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, "invalid status")
		return
	}

	request, err := h.service.UpdateStatus(c.Request.Context(), id, body.Status)
	if errors.Is(err, service.ErrInvalidStatus) {
		respondError(c, http.StatusBadRequest, "invalid status")
		return
	}
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
