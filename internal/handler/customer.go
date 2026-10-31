package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"request-service/internal/service"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) Register(r gin.IRouter) {
	r.GET("/api/v1/customers/:id", h.Get)
}

func (h *CustomerHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	customer, err := h.service.Get(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		respondError(c, http.StatusNotFound, "customer not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, customer)
}
